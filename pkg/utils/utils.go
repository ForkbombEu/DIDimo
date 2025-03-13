package utils

import (
	"os"
)

func GetEnvVariable(name string, defaultValue string, required bool) string {
	output := os.Getenv(name)
	if output == "" {
		output = defaultValue
	}
	if output == "" && required {
		panic("The environment variable " + name + " is required")
	}
	return output
}


