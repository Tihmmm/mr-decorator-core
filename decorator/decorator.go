package decorator

import (
	"github.com/Tihmmm/mr-decorator-core/client"
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
	c    client.Client
}

func NewDecorator(m string, c client.Client) Decorator {
	return &MRDecorator{
		mode: m,
		c:    c,
	}
}

const waitTime = 4 * time.Second // waiting for artifacts to be loaded

func (d *MRDecorator) Decorate(mr *models.MRRequest, prsr parser.Parser) error {
	if d.mode == ModeServer {
		time.Sleep(waitTime)
	}

	log.Printf("%s Started processing request for project: %d, merge request id: %d, job id: %d\n", time.Now().Format(time.DateTime), mr.ProjectId, mr.MergeRequestIid, mr.JobId)

	artifactsDir := ""
	if d.mode == ModeCli && mr.FilePath == "" {
		artifactsDir, err := d.c.DownloadArtifact(mr.ProjectId, mr.JobId, mr.ArtifactFileName, mr.AuthToken)
		if err != nil {
			log.Printf("Error getting artifact: %v\n", err)
			return err
		}
		defer file.DeleteDirectory(artifactsDir)
	} else {
		artifactsDir, mr.ArtifactFileName = filepath.Split(mr.FilePath)
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
