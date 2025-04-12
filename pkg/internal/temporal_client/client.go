// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package temporalclient

import (
	"fmt"

	"github.com/forkbombeu/credimi/pkg/utils"
	"go.temporal.io/sdk/client"
)

func getTemporalClient(args ...string) (client.Client, error) {
	namespace := "default"
	if len(args) > 0 {
		namespace = args[0]
	}

	hostPort := utils.GetEnvironmentVariable("TEMPORAL_HOST_PORT", "localhost:7233")

	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %v", err)
	}
	return c, nil
}

func GetTemporalClient() (client.Client, error) {
	c, err := getTemporalClient()
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %v", err)
	}

	return c, nil
}

func GetTemporalClientWithNamespace(namespace string) (client.Client, error) {
	c, err := getTemporalClient(namespace)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %v", err)
	}

	return c, nil
}
