package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	DownloadArtifact(projectId int, jobId int, artifactFileName string, glToken string) (artifactDir string, err error)
	SendNote(note string, projectId int, mergeRequestIid int, glToken string) (err error)
}

const (
	artifactsBaseDir           = "artifacts"
	contentTypeHeaderKey       = "Content-Type"
	contentTypeJsonHeaderValue = "application/json"
)

func newBaseGetRequest(path string, glToken string, host string) (*http.Request, error) {
	fullLink := "https://" + host + path
	req, err := http.NewRequest(http.MethodGet, fullLink, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating GET request for job artifact '%s': %s\n", path, err))
	}
	req.Header.Set(privateTokenHeaderKey, glToken)

	return req, nil
}

func newBasePostRequest(path string, body io.Reader, glToken string, host string) (*http.Request, error) {
	fullLink := "https://" + host + path
	req, err := http.NewRequest(http.MethodPost, fullLink, body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating POST request for job artifact '%s': %s\n", path, err))
	}
	req.Header.Set(privateTokenHeaderKey, glToken)
	req.Header.Set(contentTypeHeaderKey, contentTypeJsonHeaderValue)

	return req, nil
}
