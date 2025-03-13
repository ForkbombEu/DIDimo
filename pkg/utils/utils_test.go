package utils

import (
	"testing"
)

func Test_GetEnvVariable(t *testing.T) {
	t.Run("test with an existing variable", func(t *testing.T) {
		envVar := GetEnvVariable("GOPATH", "", false)
		if envVar == "" {
			t.Errorf("Expected a value for the environment variable GOPATH, got an empty string")
		}
	})

	t.Run("test with a non-existing variable", func(t *testing.T) {
		envVar := GetEnvVariable("NON_EXISTING_ENV_VAR", "", false)
		if envVar != "" {
			t.Errorf("Expected an empty string for the environment variable NON_EXISTING_ENV_VAR, got a value")
		}
	})

	t.Run("pass default value to a non existing variable", func(t *testing.T) {
		got := GetEnvVariable("NON_EXISTING_ENV_VAR", "default", false)
		want := "default"
		if got != want {
			t.Errorf("Expected 'default' for the environment variable NON_EXISTING_ENV_VAR, got a value")
		}
	})
	t.Run("test a required variable", func(t *testing.T) {
		envVar := GetEnvVariable("GOPATH", "", true)
		if envVar == "" {
			t.Errorf("Expected a value for the environment variable GOPATH, got an empty string")
		}
	})

	t.Run("test a required variable non existing should return a error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		GetEnvVariable("NON_EXISTING_ENV_VAR", "", true)
	})

	t.Run("test a required  non existing variable with default value should not return a error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The code panicked")
			}
		}()
		GetEnvVariable("NON_EXISTING_ENV_VAR", "default", true)
	})
}