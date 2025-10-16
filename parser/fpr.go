package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/Tihmmm/mr-decorator-core/pkg/file"
)

const (
	fpruCritScriptPath = "./FPRU_crit.sh"
	fpruHighScriptPath = "./FPRU_high.sh"
	criticalCountFile  = "critical_count.txt"
	highCountFile      = "high_count.txt"
	criticalCsv        = "critical.csv"
	highCsv            = "high.csv"
)

type fpr struct {
	highRecords     []fprRecord
	criticalRecords []fprRecord
	highCount       int
	criticalCount   int
}

type fprRecord struct {
	category        string
	path            string
	sscVulnInstance string
}

func (f *fpr) vulnCount() int {
	return f.criticalCount + f.highCount
}

func ParseFprFile(dir string, dest *fpr) (err error) {
	if err := extractVulns(dir); err != nil {
		return errors.New(fmt.Sprintf("Error parsing fpr: %s\n", err))
	}

	dest.criticalCount, err = extractVulnCount(dir, criticalCountFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing critical count: %s\n", err))
	}
	dest.highCount, err = extractVulnCount(dir, highCountFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing high count: %s\n", err))
	}

	criticalRecords, err := extractRecords(dir, criticalCsv)
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing critical records: %s\n", err))
	}
	dest.criticalRecords = criticalRecords

	highRecords, err := extractRecords(dir, highCsv)
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing high records: %s\n", err))
	}
	dest.highRecords = highRecords

	return nil
}

func extractVulns(fileDir string) error {
	if err := exec.Command(fpruCritScriptPath, fileDir).Run(); err != nil {
		return errors.New(fmt.Sprintf("Error extracting critical vulns: %s\n", err))
	}
	if err := exec.Command(fpruHighScriptPath, fileDir).Run(); err != nil {
		return errors.New(fmt.Sprintf("Error extracting high vulns: %s\n", err))
	}
	return nil
}

func extractVulnCount(dir, subpath string) (int, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Error opening directory root: %s\n", err))
	}
	defer func(root *os.Root) {
		if err := root.Close(); err != nil {
			log.Printf("Error closing root: %s\n", err)
			return
		}
	}(root)

	vulnCountFile, err := root.Open(subpath)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Error opening vulns count file: %s\n", err))
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %s\n", err)
			return
		}
	}(vulnCountFile)

	var count int
	scanner := bufio.NewScanner(vulnCountFile)
	for scanner.Scan() {
		lineStr := scanner.Text()
		count, _ = strconv.Atoi(lineStr)
	}
	return count, nil
}

func extractRecords(dir, subpath string) ([]fprRecord, error) {
	records, err := file.ReadCsv(dir, subpath)
	if err != nil {
		return []fprRecord{}, errors.New(fmt.Sprintf("Error extracting fpr records"))
	}
	var fprRecords []fprRecord
	for i := 1; i < len(records); i++ {
		fprRec := fprRecord{
			category:        records[i][1],
			path:            records[i][2],
			sscVulnInstance: records[i][0],
		}
		fprRecords = append(fprRecords, fprRec)
	}

	return fprRecords, nil
}
