package decorator

import (
	"errors"
	"github.com/Tihmmm/mr-decorator-core/client"
	"github.com/Tihmmm/mr-decorator-core/config"
	custErrors "github.com/Tihmmm/mr-decorator-core/errors"
	"github.com/Tihmmm/mr-decorator-core/models"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/pkg/file"
	"log"
	"path/filepath"
	"time"
)

const (
	ModeServer = "server"
	ModeCli    = "cli"
)

type Decorator interface {
	Decorate(mr *models.MRRequest, prsr parser.Parser) error
}

type MRDecorator struct {
	mode string // either `cli` or `server`
	cfg  config.DecoratorConfig
	c    client.Client
}

func NewDecorator(m string, cfg config.DecoratorConfig, c client.Client) Decorator {
	return &MRDecorator{
		mode: m,
		cfg:  cfg,
		c:    c,
	}
}

func (d *MRDecorator) Decorate(mr *models.MRRequest, prsr parser.Parser) error {
	log.Printf("%s Started processing request for project: %d, merge request id: %d, job id: %d\n", time.Now().Format(time.DateTime), mr.ProjectId, mr.MergeRequestIid, mr.JobId)

	artifactsDir := ""
	var err error
	if d.mode == ModeCli && mr.FilePath != "" {
		artifactsDir, mr.ArtifactFileName = filepath.Split(mr.FilePath)
	} else {
		retryCount := 0
		for retryCount < d.cfg.ArtifactDownloadMaxRetries {
			artifactsDir, err = d.c.DownloadArtifact(mr.ProjectId, mr.JobId, mr.ArtifactFileName, mr.AuthToken)
			if err != nil && errors.Is(err, &custErrors.DownloadError{}) {
				retryCount++
				time.Sleep(time.Duration(d.cfg.ArtifactDownloadRetryDelay) * time.Second)
				continue
			}
		}
		if err != nil {
			log.Printf("Error getting artifact: %v\n", err)
			return err
		}
	}
	if d.mode == ModeServer {
		defer file.DeleteDirectory(artifactsDir)
	}

	note, err := prsr.GetNoteFromReportFile(artifactsDir, mr.ArtifactFileName, mr.VulnerabilityMgmtId)

	err = d.c.SendNote(note, mr.ProjectId, mr.MergeRequestIid, mr.AuthToken)
	if err != nil {
		log.Printf("Error sending note for mr iid: '%d' in project %d: %v\n", mr.MergeRequestIid, mr.ProjectId, err)
		return err
	}

	log.Printf("%s Finished processing request for project: %d, merge request id: %d, job id: %d\n", time.Now().Format(time.DateTime), mr.ProjectId, mr.MergeRequestIid, mr.JobId)

	return nil
}
