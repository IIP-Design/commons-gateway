// The MIT License (MIT)

// Copyright (c) 2023 CyberMonkey SP. Z O.O.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package turnstile

import (
	"io"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Turnstile struct {
	Secret       string
	TurnstileURL string
}

type Response struct {
	// Success indicates if the challenge was passed
	Success bool `json:"success"`
	// ChallengeTs is the timestamp of the captcha
	ChallengeTs string `json:"challenge_ts"`
	// Hostname is the hostname of the passed captcha
	Hostname string `json:"hostname"`
	// ErrorCodes contains error codes returned by hCaptcha (optional)
	ErrorCodes []string `json:"error-codes"`
	// Action  is the customer widget identifier passed to the widget on the client side
	Action string `json:"action"`
	// CData is the customer data passed to the widget on the client side
	CData string `json:"cdata"`
}

func NewTS(secret string) *Turnstile {
	return &Turnstile{
		Secret:       secret,
		TurnstileURL: "https://challenges.cloudflare.com/turnstile/v0/siteverify",
	}
}

// Verify verifies a "h-captcha-response" data field, with an optional remote IP set.
func (t *Turnstile) Verify(response, remoteip string) (*Response, error) {
	values := url.Values{"secret": {t.Secret}, "response": {response}}
	if remoteip != "" {
		values.Set("remoteip", remoteip)
	}

	resp, err := http.PostForm(t.TurnstileURL, values)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP read error: %w", err)
	}

	r := Response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("JSON error: %w", err)
	}

	return &r, nil
}

func VerifyToken(token string, remoteAddr string, privateKey string) error {
	ts := NewTS(privateKey)
	resp, err := ts.Verify(token, remoteAddr)

	if err != nil || resp == nil || !resp.Success {
		return err
	} else {
		return nil
	}
}
