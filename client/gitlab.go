package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cfg "github.com/Tihmmm/mr-decorator-core/config"
	custErrors "github.com/Tihmmm/mr-decorator-core/errors"
	"github.com/Tihmmm/mr-decorator-core/models"
	"github.com/doyensec/safeurl"
	"github.com/google/uuid"
)

type GitlabClient struct {
	client *safeurl.WrappedClient
	cfg    cfg.BaseClientConfig
}

func NewGitlabClient(cfg cfg.BaseClientConfig) Client {
	config := safeurl.GetConfigBuilder().
		SetAllowedIPs(cfg.Ip).
		SetTimeout(time.Duration(cfg.TimeoutSeconds) * time.Second).
		Build()
	httpClient := &GitlabClient{
		cfg:    cfg,
		client: safeurl.Client(config),
	}
	return httpClient
}

const (
	jobArtifactsEndpointBasePath      = "/api/v4/projects/%d/jobs/%d/artifacts/%s"
	mergeRequestNotesEndpointBasePath = "/api/v4/projects/%d/merge_requests/%d/notes"
	privateTokenHeaderKey             = "PRIVATE-TOKEN"
)

func (c *GitlabClient) DownloadArtifact(projectId int, jobId int, artifactFileName string, glToken string) (artifactDir string, err error) {
	jobArtifactPath := fmt.Sprintf(jobArtifactsEndpointBasePath, projectId, jobId, artifactFileName)
	req, err := newBaseGetRequest(jobArtifactPath, glToken, c.cfg.Host)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error creating GET request for job artifact '%s': %s\n", jobArtifactPath, err))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error downloading artifact for project: %d, job: %d, err: %s\n", projectId, jobId, err))
	}
	if resp.StatusCode != http.StatusOK {
		return "", &custErrors.DownloadError{Err: fmt.Sprintf("Error downloading artifact for project: %d, job: %d. Gitlab response status: %d\n", projectId, jobId, resp.StatusCode)}
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Printf("Error closing response body for project: %d, job: %d, err: %s\n", projectId, jobId, err)
			return
		}
	}(resp.Body)

	dirPath := filepath.Join(artifactsBaseDir, uuid.New().String())
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return "", errors.New(fmt.Sprintf("Error creating artifact directory: %s\n", err))
	}
	dirRoot, err := os.OpenRoot(dirPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error opening artifact root dir '%s': %s\n", dirPath, err))
	}
	out, err := dirRoot.Create(artifactFileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error creating artifact file: %s\n", err))
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", errors.New(fmt.Sprintf("Error copying artifact for project: %d, job: %d, err: %v\n", projectId, jobId, err))
	}

	return dirPath, nil
}

func (c *GitlabClient) SendNote(note string, projectId int, mergeRequestIid int, glToken string) (err error) {
	body := models.NoteRequest{Body: note}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing body: %s\n", err))
	}

	notePath := fmt.Sprintf(mergeRequestNotesEndpointBasePath, projectId, mergeRequestIid)

	req, err := newBasePostRequest(notePath, bytes.NewBuffer(bodyBytes), glToken, c.cfg.Host)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating POST request to send node for job artifact '%s': %v\n", notePath, err))
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Error making note request: %v", err))
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Printf("Error closing response body for job artifact '%s': %v\n", notePath, err)
			return
		}
	}(resp.Body)

	var respBuf []byte
	_, err = resp.Body.Read(respBuf)

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Error sending note. Gitlab response status code: %d\nbody: %s\n", resp.StatusCode, string(respBuf)))
	}

	return nil
}
