package temporalclient

import (
	"fmt"

	"go.temporal.io/sdk/client"

	"github.com/forkbombeu/didimo/pkg/utils"
)

func GetTemporalClient() (client.Client, error) {
	hostPort := utils.GetEnvironmentVariable("TEMPORAL_ADDRESS", "localhost:7233")
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %v", err)
	}
	return c, nil
}
