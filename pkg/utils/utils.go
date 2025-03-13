package utils

import (
	"os"
	"strconv"
	"math"
	"fmt"
)

// GetEnvVariable retrieves the value of an environment variable.
//
// Parameters:
//   - name: The name of the environment variable to retrieve.
//   - defaultValue: The value to return if the environment variable is not set.
//   - required: A boolean indicating whether the environment variable is required.
//
// Returns:
//   - string: The value of the environment variable, or the default value if not set.
//     If the variable is required and not set, the function panics.
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

// GetEnvVariableAsInt retrieves an environment variable and converts it to an integer.
//
// Parameters:
//   - name: The name of the environment variable to retrieve.
//   - defaultValue: The default integer value to return if the environment variable is not set.
//   - required: A boolean indicating whether the environment variable is required.
//
// Returns:
//   - int: The integer value of the environment variable, or the default value if not set.
//   - error: An error if the variable cannot be parsed as an integer or if it's out of range for int type.
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


