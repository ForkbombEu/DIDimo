package pb

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	engine "github.com/forkbombeu/didimo/pkg/template_engine"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func RouteParseTestsConfig(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/parse-config", func(e *core.RequestEvent) error {
			testId := e.Request.URL.Query().Get("test_id")
			if testId == "" {
				return apis.NewBadRequestError("test_id is required", nil)
			}
			provider := e.Request.URL.Query().Get("provider")
			if provider == "" {
				return apis.NewBadRequestError("provider is required", nil)
			}
			templates, err := getTemplatesByFolderId(testId)
			if err != nil {
				return apis.NewBadRequestError("Error reading test suite folder", err)
			}

			variables, err := e.App.FindRecordsByFilter("config_variables", "provider = {:provider} && test_suite = {:testId}", "", 0, 0, dbx.Params{"provider": provider, "testId": testId})
			if err != nil {
				return apis.NewBadRequestError("Error fetching variables", err)
			}

			if len(variables) == 0 {
				return apis.NewNotFoundError("variables not found", nil)
			}

			readers := parseFilesAsReaders(templates)

			neededVariables, err := engine.GetPlaceholders(readers)
			if err != nil {
				return apis.NewBadRequestError("Error getting placeholders", err)
			}

			for _, variable := range neededVariables {
				found := false
				for _, v := range variables {
					if v.Get("field_name") == variable {
						found = true
						break
					}
				}
				if !found {
					return apis.NewBadRequestError("Variable "+variable+" not found", nil)
				}
			}

			var renderedTemplates []string
			for _, template := range templates {
				template.Seek(0, 0)
				template_variables, err := engine.GetPlaceholders([]io.Reader{template})
				if err != nil {
					return apis.NewBadRequestError("Error getting placeholders", err)
				}
				values := make(map[string]interface{})
				for _, variable := range variables {
					name, ok := variable.Get("field_name").(string)
					if !ok {
						return apis.NewBadRequestError("Invalid variable name type", nil)
					}
					if slices.Contains(template_variables, name) {
						values[name] = variable.Get("value")
					}
				}
				template.Seek(0, 0)
				rendered, err := engine.RenderTemplate(template, values)
				if err != nil {
					return apis.NewInternalServerError("Error rendering template", err)
				}
				renderedTemplates = append(renderedTemplates, rendered)
			}

			return e.JSON(http.StatusOK, map[string]interface{}{
				"templates": renderedTemplates,
			})
		})
		// .Bind(apis.RequireAuth())
		return se.Next()
	})
}

func RouteNormalizedPlaceholders(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/normalized-placeholders", func(e *core.RequestEvent) error {
			testId := e.Request.URL.Query().Get("test_id")
			if testId == "" {
				return apis.NewBadRequestError("test_id is required", nil)
			}
			templates, err := getTemplatesByFolderId(testId)
			if err != nil {
				return apis.NewBadRequestError("Error reading test suite folder", err)
			}
			if templates == nil {
				return apis.NewBadRequestError("Error reading test suite folder", nil)
			}

			readers := parseFilesAsReaders(templates)

			placeholders, err := engine.GetPlaceholders(readers)
			if err != nil {
				return apis.NewBadRequestError("Error getting placeholders", err)
			}

			return e.JSON(http.StatusOK, map[string]interface{}{
				"placeholders": placeholders,
			})
		})
		return se.Next()
	})
}

func getTemplatesByFolderId(folderId string) ([]*os.File, error) {
	var templates []*os.File
	err := filepath.Walk("./"+folderId, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		templates = append(templates, file)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return templates, nil
}

func parseFilesAsReaders(files []*os.File) []io.Reader {
	readers := make([]io.Reader, len(files))
	for i, file := range files {
		readers[i] = file
	}
	return readers
}

