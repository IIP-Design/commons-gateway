package logs

import (
	"encoding/json"
	"log"
)

// LogError formats error messages with a bespoke message and log it to CloudWatch.
func LogError(err error, message string) {
	errMessage := map[string]interface{}{
		"Message": message,
		"Error":   err.Error(),
	}

	fullMessage, _ := json.MarshalIndent(errMessage, "", "  ")

	log.Println(string(fullMessage))
}
