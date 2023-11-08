package testConfig

import (
	"os"
	"strings"
)

func AddToEnv(env map[string]string) {
	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range env {
		if !currentEnv[key] {
			_ = os.Setenv(key, value)
		}
	}
}

func ConfigureAprimo() {
	AddToEnv(AprimoEnv)
}

func ConfigureAws() {
	AddToEnv(AwsEnv)
}

func ConfigureDb() {
	AddToEnv(DbEnv)
}

func ConfigureEmail() {
	AddToEnv(EmailEnv)
}
