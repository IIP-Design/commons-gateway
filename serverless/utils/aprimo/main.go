package aprimo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	Key       string `json:"key"`
	FileType  string `json:"filetype"`
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

type RecordCreationResponse struct {
	Id string `json:"id"`
}

// GetEndpointURL constructs an URL for a given Aprimo API endpoint. Authorization
// is handled by the Aprimo Marketing Operations API, while most other operations should
// be directed to the DAM API by setting the `dam` parameter to `true`.
func GetEndpointURL(endpoint string, dam bool) string {
	var url string

	if dam {
		url = fmt.Sprintf("https://%s.dam.aprimo.com/api/core/%s", os.Getenv("APRIMO_TENANT"), endpoint)
	} else {
		url = fmt.Sprintf("https://%s.aprimo.com/%s", os.Getenv("APRIMO_TENANT"), endpoint)
	}

	return url
}

// GetAuthToken requests a bearer token from the Aprimo Marketing Operations API.
// The resultant token can be used to authenticate to the Aprimo DAM API for record
// creation and file uploading.
func GetAuthToken() (string, error) {
	var err error
	var token string

	endpoint := GetEndpointURL("login/connect/token", false)

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

func PostJsonData(endpoint string, token string, reqBody string, useDam bool) ([]byte, int, error) {
	url := GetEndpointURL(endpoint, useDam)
	return PostJsonDataWithFqUrl(url, token, reqBody)
}

func PostJsonDataWithFqUrl(url string, token string, reqBody string) ([]byte, int, error) {
	var statusCode int
	var res []byte
	var err error

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

	resp, statusCode, err := PostJsonData("uploads/segments", token, reqBody, false)
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

	resp, statusCode, err := PostJsonDataWithFqUrl(fmt.Sprintf("%s/commit", uri), token, reqBody)
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

	// Uses returned URI directly
	url := fmt.Sprintf("%s?index=%d", uri, seg.Segment)

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

func UploadFile(filename string, fileType string, data *bytes.Buffer, token string) (string, error) {
	var uploadToken string
	var err error

	// Using DAM will only work for very small files, e.g., 2.8 MB will fail
	url := GetEndpointURL("uploads", false)

	// Add file data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", fileType)
	partHeaders.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, filename, filename))
	part, _ := writer.CreatePart(partHeaders)

	io.Copy(part, data)
	writer.Close()

	client := &http.Client{}
	request, err := http.NewRequest(
		http.MethodPost,
		url,
		body,
	)

	if err != nil {
		logs.LogError(err, "Error Preparing Aprimo Request")
		return uploadToken, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Set("Accept", "*/*")
	request.Header.Set("API-VERSION", "1")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make the request
	resp, err := client.Do(request)

	if err != nil || resp.StatusCode >= 400 {
		logs.LogError(err, "Aprimo File Upload Error")
		return uploadToken, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.LogError(err, "Error Reading Response Body")
		return uploadToken, err
	}

	var commitResp UploadCommitResponse
	json.Unmarshal(respBody, &commitResp)
	uploadToken = commitResp.Token

	return uploadToken, nil
}

func SubmitRecord(description string, event FileRecordInitEvent, team string, accessToken string) (string, error) {
	var id string
	var err error

	reqBody := fmt.Sprintf(`{
		"status":"draft",
		"fields": {
			"addOrUpdate": [
				{
					"Name": "Description",
					"localizedValues": [
						{ "value": "%s" }
					]
				},
				{
					"Name": "DisplayTitle",
					"localizedValues": [
						{ "value": "%s" }
					]
				},
				{
					"Name": "Team",
					"localizedValues": [
						{ "values": ["%s"] }
					]
			}
			]
		},
		"files": {
			"master": "%s",
			"addOrUpdate": [
				{
					"versions": {
						"addOrUpdate": [
							{
								"id": "%s",
								"filename": "%s"
							}
						]
					}
				}
			]
		}
	}`, description, event.Key, team, event.FileToken, event.FileToken, event.Key)

	respBody, statusCode, err := PostJsonData("records", accessToken, reqBody, true)
	if err != nil {
		return id, err
	} else if statusCode >= 400 {
		log.Printf("Return status: %d\n", statusCode)
	}

	var res RecordCreationResponse
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return id, err
	}

	return res.Id, nil
}

func DeleteRecord(recordId string, token string) error {
	var err error

	url := GetEndpointURL(fmt.Sprintf("record/%s", recordId), true)

	client := &http.Client{}
	request, err := http.NewRequest(
		http.MethodDelete,
		url,
		nil,
	)

	if err != nil {
		logs.LogError(err, "Error Preparing Aprimo Record Delete")
		return err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("API-VERSION", "1")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(request)

	if err != nil {
		logs.LogError(err, "Error Deleting Aprimo Record")
		return err
	} else if resp.StatusCode != 204 {
		logs.LogError(err, fmt.Sprintf("Error response from Aprimo: %d", resp.StatusCode))
		return err
	}

	return nil
}
