package test

import (
	"os"
	"strings"
)

func AddToEnv() {
	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range Env {
		if !currentEnv[key] {
			_ = os.Setenv(key, value)
		}
	}
}
