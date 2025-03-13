package utils

import (
	"os"
	"strconv"
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

func GetEnvVariableAsInt(name string, defaultValue int, required bool) (int, error) {
	output := GetEnvVariable(name, "", required)
	if output == "" {
		return defaultValue, nil
	}
	outputAsInt, err := strconv.ParseInt(output, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(outputAsInt), nil
}


