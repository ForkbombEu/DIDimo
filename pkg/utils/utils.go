package utils

import (
	"os"
	"strconv"
	"math"
	"fmt"
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
	if outputAsInt > math.MaxInt || outputAsInt < math.MinInt {
		return 0, fmt.Errorf("value out of range for int: %d", outputAsInt)
	}
	return int(outputAsInt), nil
}


