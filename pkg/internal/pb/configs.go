package pb

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func getTemplatesByFolder(folder string) ([]*os.File, error) {
	var templates []*os.File
	err := filepath.Walk("./config_templates/"+folder, func(path string, info os.FileInfo, err error) error {
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

func RouteGetConfigsTemplates(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/conformance-checks/configs/get-configs-templates", func(e *core.RequestEvent) error {
			testId := e.Request.URL.Query().Get("test_id")
			if testId == "" {
				testId = "OpenID4VP_Wallet/OpenID_Foundation"
			}
			files, err := getTemplatesByFolder(testId)
			if err != nil {
				return apis.NewBadRequestError("Error reading test suite folder", err)
			}
			var variants []string
			for _, file := range files {
				variants = append(variants, strings.Replace(file.Name(), testId+"/", "", 1))
			}
			return e.JSON(http.StatusOK, map[string]interface{}{
				"variants": variants,
			})
		})
		return se.Next()
	})
}
