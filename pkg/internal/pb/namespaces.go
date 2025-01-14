package pb

import (
	"context"
	"log"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"google.golang.org/protobuf/types/known/durationpb"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func HookNamespaceOrgs(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess("organizations").BindFunc(func(e *core.RecordEvent) error {
		c, err := client.NewNamespaceClient(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer c.Close()

		err = c.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
			Namespace:                        e.Record.Get("name").(string),
			WorkflowExecutionRetentionPeriod: durationpb.New(7 * 24 * time.Hour),
		})

		if err != nil {
			log.Fatalf("Unable to create namespace %s: %v", e.Record.Get("name").(string), err)
		}

		log.Default().Printf("Namespace %s created", e.Record.Get("name").(string))
		return e.Next()
	})
}
