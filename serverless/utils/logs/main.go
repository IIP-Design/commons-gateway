package logs

import (
	"encoding/json"
	"errors"
	"log"
)

// LogError formats error messages with a bespoke message and log it to CloudWatch.
func LogError(err error, message string) {
	if err == nil {
		err = errors.New("unknown error")
	}

	errMessage := map[string]interface{}{
		"Message": message,
		"Error":   err.Error(),
	}

	fullMessage, _ := json.MarshalIndent(errMessage, "", "  ")

	log.Println(string(fullMessage))
}
