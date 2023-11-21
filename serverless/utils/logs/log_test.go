package logs

import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"testing"
)

func logCommon(err error) string {
	body := &bytes.Buffer{}
	log.SetOutput(body)
	LogError(err, "test")
	return body.String()
}

func matchLog(l string, matches []string, t *testing.T) {
	re := regexp.MustCompile(`(?mi)\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} {\s+\"Error\": \"([\w\s]+)\",\s+\"Message\": \"([\w\s]+)\"\s+}`)

	for i, match := range re.FindAllString(l, -1) {
		if i > 0 && match != matches[i-1] {
			t.Fatalf(`LogError failure, have %s, want %s`, match, matches[i-1])
		}
	}
}

func TestNoError(t *testing.T) {
	ret := logCommon(nil)
	matchLog(ret, []string{"unknown error", "test"}, t)
}

func TestRealError(t *testing.T) {
	ret := logCommon(errors.New("test"))
	matchLog(ret, []string{"test", "test"}, t)
}
