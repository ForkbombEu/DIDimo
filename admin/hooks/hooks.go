package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"golang.org/x/exp/slices"
)

func Register(app *pocketbase.PocketBase) error {
	modelHandler := func(event string) func(e *core.ModelEvent) error {
		return func(e *core.ModelEvent) error {
			table := e.Model.TableName()
			// we don't want to executeEventActions if the event is a system event (e.g. "_collections" changes)
			if record, ok := e.Model.(*models.Record); ok {
				executeEventActions(app, event, table, record)
			} else {
				log.Println("Skipping executeEventActions for table:", table)
			}
			return nil
		}
	}
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		app.OnModelAfterCreate().Add(modelHandler("insert"))
		app.OnModelAfterUpdate().Add(modelHandler("update"))
		app.OnModelAfterDelete().Add(modelHandler("delete"))
		app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
			return doAudit(app, "insert", e.Record, e.HttpContext)
		})
		app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
			return doAudit(app, "update", e.Record, e.HttpContext)
		})
		app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
			return doAudit(app, "delete", e.Record, e.HttpContext)
		})
		return nil
	})
	return nil
}

// collection names to be audit logged
var collections = strings.Split(os.Getenv("AUDITLOG"), ",")

func doAudit(app *pocketbase.PocketBase, event string, record *models.Record, ctx echo.Context) error {
	collection := record.Collection().Name
	// exclude logging "auditlog" and include only what's in AUDITLOG env var
	if collection != "auditlog" && slices.Contains(collections, collection) {
		var user, admin string
		if u, ok := ctx.Get(apis.ContextAdminKey).(*models.Admin); ok {
			admin = u.Id
		}
		if u, ok := ctx.Get(apis.ContextAuthRecordKey).(*models.Record); ok {
			user = u.Id
		}
		log.Printf("AuditLog:%s:%s:%s:%s:%s\n", collection, record.Id, event, user, admin)
		target, err := app.Dao().FindCollectionByNameOrId("auditlog")
		if err != nil {
			return err
		}
		auditlog := models.NewRecord(target)
		auditlog.Set("collection", collection)
		auditlog.Set("record", record.Id)
		auditlog.Set("event", event)
		auditlog.Set("user", user)
		auditlog.Set("admin", admin)
		auditlog.Set("data", record)

		return app.Dao().SaveRecord(auditlog)
	}
	return nil
}

func executeEventActions(app *pocketbase.PocketBase, event string, table string, record *models.Record) {
	// TODO: Load and cache this. Reload only on changes to "hooks" table
	rows := []dbx.NullStringMap{}
	app.DB().Select("action_type", "action", "action_params", "expands").
		From("hooks").
		Where(dbx.HashExp{"collection": table, "event": event, "disabled": false}).
		All(&rows)
	for _, row := range rows {
		action_type := row["action_type"].String
		action := row["action"].String
		action_params := row["action_params"].String
		expands := strings.Split(row["expands"].String, ",")
		app.Dao().ExpandRecord(record, expands, func(c *models.Collection, ids []string) ([]*models.Record, error) {
			return app.Dao().FindRecordsByIds(c.Name, ids, nil)
		})
		switch action_type {
		case "sendmail":
			if err := doSendMail(app, action, action_params, record); err != nil {
				log.Println("ERROR", err)
			}
		default:
			if err := executeEventAction(event, table, action_type, action, action_params, record); err != nil {
				log.Println("ERROR", err)
			}
		}
	}
}

func executeEventAction(event, table, action_type, action, action_params string, record *models.Record) error {
	log.Printf("event:%s, table: %s, action: %s\n", event, table, action)
	switch action_type {
	case "command":
		return doCommand(action, action_params, record)
	case "post":
		return doPost(action, action_params, record)
	case "restroom-mw":
		return doRestroomMW(action, action_params, record)
	default:
		return errors.New(fmt.Sprintf("Unknown action_type: %s", action_type))
	}
}

func doSendMail(app *pocketbase.PocketBase, action, action_params string, record *models.Record) error {
	params := struct {
		Subject    string `json:"subject"`
		OwnerField string `json:"ownerField"`
	}{
		Subject:    "New message",
		OwnerField: "owner",
	}
	if action_params != "" {
		err := json.Unmarshal([]byte(action_params), &params)
		if err != nil {
			return err
		}
	}
	var emails []string
	owner := record.Get(params.OwnerField)
	if o, ok := owner.(string); ok {
		userTo, err := app.Dao().FindRecordById("users", o)
		if err != nil {
			return err
		}
		emails = []string{userTo.Email()}
	} else if os, ok := owner.([]string); ok {
		records, err := app.Dao().FindRecordsByIds("users", os)
		if err != nil {
			return err
		}
		for _, x := range records {
			emails = append(emails, x.Email())
		}
	} else {
		return errors.New(fmt.Sprintf("Unknown record"))
	}
	// TODO: send mail in multiple go routines
	var err error = nil
	for _, email := range emails {
		message := &mailer.Message{
			From: mail.Address{
				Address: app.Settings().Meta.SenderAddress,
				Name:    app.Settings().Meta.SenderName,
			},
			To: []mail.Address{
				{Address: email},
			},
			Subject: params.Subject,
			HTML:    action,
		}
		e := app.NewMailClient().Send(message)
		if e != nil {
			if err == nil {
				err = e
			} else {
				err = fmt.Errorf("%w; %w", err, e)
			}
		}
	}

	return err
}

func doCommand(action, action_params string, record *models.Record) error {
	cmd := exec.Command(action, action_params)
	if w, err := cmd.StdinPipe(); err != nil {
		return err
	} else {
		if r, err := cmd.StdoutPipe(); err != nil {
			return err
		} else {
			go func() {
				defer w.Close()
				defer r.Close()
				log.Println("-------------------------------")
				defer log.Println("-------------------------------")
				if err := cmd.Start(); err != nil {
					log.Printf("command start failed: %s %+v\n", action, err)
				} else {
					// write JSON into the pipe and close
					json.NewEncoder(w).Encode(record)
					w.Close()
					if err := cmd.Wait(); err != nil {
						log.Printf("command wait failed: %s %+v\n", action, err)
					}
				}
			}()
			// read pipe's stdout and copy to ours (in parallel to the above goroutine)
			io.Copy(os.Stdout, r)
		}
	}
	return nil
}

func doPost(action, action_params string, record *models.Record) error {
	r, w := io.Pipe()
	defer w.Close()
	go func() {
		defer r.Close()
		if resp, err := http.Post(action, "application/json", r); err != nil {
			log.Println("POST failed", action, err)
		} else {
			io.Copy(os.Stdout, resp.Body)
		}
	}()
	if err := json.NewEncoder(w).Encode(record); err != nil {
		log.Println("ERROR writing to pipe", err)
	}
	return nil
}

func doRestroomMW(action, action_params string, record *models.Record) error {
	// Parse action params
	params := struct {
		Wrapper string `json:"wrapper"`
		Method  string `json:"method"`
	}{
		Wrapper: "",
		Method:  "post",
	}
	if action_params != "" {
		err := json.Unmarshal([]byte(action_params), &params)
		if err != nil {
			return err
		}
	}
	r, w := io.Pipe()
	defer w.Close()

	var reqf func() (*http.Response, error)
	if params.Method == "post" {
		reqf = func() (*http.Response, error) {
			return http.Post(action, "application/json", r)
		}
	} else if params.Method == "get" {
		reqf = func() (*http.Response, error) {
			buffer := new(bytes.Buffer)
			buffer.ReadFrom(r)
			return http.Get(fmt.Sprintf("%s?%s", action, buffer.String()))
		}
	} else {
		return fmt.Errorf("Unknown method %s", params.Method)
	}

	go func() {
		defer r.Close()
		if resp, err := reqf(); err != nil {
			log.Println("HTTP request failed", action, err)
		} else {
			io.Copy(os.Stdout, resp.Body)
		}
	}()

	var reqObj interface{}

	// Build request object
	if params.Wrapper != "" {
		reqObj = map[string]models.Record{
			params.Wrapper: *record,
		}
	} else {
		reqObj = record
	}

	if params.Method == "post" {
		reqObj = map[string]interface{}{
			"data": reqObj,
			"keys": map[string]interface{}{},
		}

		if err := json.NewEncoder(w).Encode(reqObj); err != nil {
			log.Println("ERROR writing to pipe", err)
		}
	} else if params.Method == "get" {
		dataParam, _ := json.Marshal(reqObj)
		fmt.Fprintf(w, "data=%s", url.QueryEscape(string(dataParam)))
	}
	return nil
}
