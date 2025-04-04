package temporalclient

import (
	"fmt"
	"os"

	"go.temporal.io/sdk/client"
)

func GetTemporalClient(namespace string) (client.Client, error) {
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = "localhost:7233"
	}
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
		Namespace: namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %v", err)
	}
	return c, nil
}
