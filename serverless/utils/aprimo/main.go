package aprimo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

type TokenResponse struct {
	Expiration int    `json:"expires_in"`
	Scope      string `json:"scope"`
	Token      string `json:"access_token"`
	Type       string `json:"token_type"`
}

type FileRecordInitEvent struct {
	AprimoId string `json:"aprimoId"`
	Filename string `json:"filename"`
	FileType string `json:"filetype"`
}

type FileRecordUpdateEvent struct {
	AprimoId  string `json:"aprimoId"`
	Filename  string `json:"filename"`
	FileToken string `json:"fileToken"`
}

type UploadSegmentResponse struct {
	Uri string `json:"uri"`
}

type UploadCommitResponse struct {
	Token string `json:"token"`
}

type FileSegment struct {
	Segment  int
	FileType string
	Data     *bytes.Buffer
}

// GetEndpointURL constructs an URL for a given Aprimo API endpoint. Authorization
// is handled by the Aprimo Marketing Operations API, while other operations should
// be directed to the DAM API by setting the `auth` parameter to `false`.
func GetEndpointURL(endpoint string, auth bool) string {
	var url string

	if auth {
		url = fmt.Sprintf("https://%s.aprimo.com/%s", os.Getenv("APRIMO_TENANT"), endpoint)
	} else {
		url = fmt.Sprintf("https://%s.dam.aprimo.com/api/core/%s", os.Getenv("APRIMO_TENANT"), endpoint)
	}

	return url
}

// GetAuthToken requests a bearer token from the Aprimo Marketing Operations API.
// The resultant token can be used to authenticate to the Aprimo DAM API for record
// creation and file uploading.
func GetAuthToken() (string, error) {
	var err error
	var token string

	endpoint := GetEndpointURL("login/connect/token", true)

	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	body.Set("scope", "api")
	body.Set("client_id", os.Getenv("APRIMO_CLIENT_ID"))
	body.Set("client_secret", os.Getenv("APRIMO_CLIENT_SECRET"))
	encodedBody := body.Encode()

	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(encodedBody))

	if err != nil {
		logs.LogError(err, "Retrieve Aprimo Auth Token Error")

		return token, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		logs.LogError(err, "Error Reading Response Body")

		return token, err
	}

	var res TokenResponse

	err = json.Unmarshal(respBody, &res)

	if err != nil {
		logs.LogError(err, "Data Un-marshalling Error")

		return token, err
	}

	return res.Token, err
}

func PostJsonData(endpoint string, token string, reqBody string) ([]byte, int, error) {
	var statusCode int
	var res []byte
	var err error

	url := GetEndpointURL(endpoint, false)
	jsonData := []byte(reqBody)
	bodyReader := bytes.NewReader(jsonData)

	client := &http.Client{}
	request, err := http.NewRequest(
		http.MethodPost,
		url,
		bodyReader,
	)

	if err != nil {
		logs.LogError(err, "Error Preparing Aprimo Request")
		return res, statusCode, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("API-VERSION", "1")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)

	if err != nil {
		logs.LogError(err, "Aprimo Response Error")
		return res, statusCode, err
	}
	statusCode = resp.StatusCode

	defer resp.Body.Close()
	res, err = io.ReadAll(resp.Body)

	if err != nil {
		logs.LogError(err, "Error Reading Response Body")
		return res, statusCode, err
	}

	return res, statusCode, nil
}

func InitFileUpload(filename string, token string) string {
	var uri string

	reqBody := fmt.Sprintf(`{
		"filename":"%s"
	}`, filename)

	resp, statusCode, err := PostJsonData("uploads/segments", token, reqBody)
	if err == nil && statusCode == 200 {
		var uriResp UploadSegmentResponse
		json.Unmarshal(resp, &uriResp)
		uri = uriResp.Uri
	}

	return uri
}

func CommitFileUpload(filename string, segments int, uri string, token string) (string, error) {
	var respToken string
	var err error

	reqBody := fmt.Sprintf(`{
		"filename":"%s",
		"segmentcount": "%d"
	}`, filename, segments)

	resp, statusCode, err := PostJsonData(fmt.Sprintf("%s/commit", uri[1:]), token, reqBody)
	if err == nil && statusCode == 200 {
		var commitResp UploadCommitResponse
		json.Unmarshal(resp, &commitResp)
		respToken = commitResp.Token
	}

	return respToken, err
}

func UploadSegment(filename string, uri string, seg *FileSegment, token string) (bool, error) {
	success := false
	var err error

	url := GetEndpointURL(fmt.Sprintf("%s?index=%d", uri[1:], seg.Segment), false)

	// Add file data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", seg.FileType)
	partHeaders.Set("Content-Disposition", fmt.Sprintf(`form-data; name="segment%d"; filename="%s.segment%d"`, seg.Segment, filename, seg.Segment))
	part, _ := writer.CreatePart(partHeaders)

	io.Copy(part, seg.Data)
	writer.Close()

	client := &http.Client{}
	request, err := http.NewRequest(
		http.MethodPost,
		url,
		body,
	)

	if err != nil {
		logs.LogError(err, "Error Preparing Aprimo Request")
		return success, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	request.Header.Set("Accept", "*/*")
	request.Header.Set("API-VERSION", "1")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make the request
	resp, err := client.Do(request)

	if err != nil {
		logs.LogError(err, "Aprimo File Segment Upload Error")
		return success, err
	}
	success = resp.StatusCode == 202

	defer resp.Body.Close()
	return success, nil
}
