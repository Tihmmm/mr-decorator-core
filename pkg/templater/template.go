package templater

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

func ExecToString(baseTemplate string, data any) (string, error) {
	compiled, err := template.New("").Parse(baseTemplate)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error compiling template: %s", err))
	}

	var buff bytes.Buffer
	if err := compiled.Execute(&buff, data); err != nil {
		return "", errors.New(fmt.Sprintf("Error executing template: %s", err))
	}

	note := regexp.MustCompile(`[\n{}]`).ReplaceAllString(buff.String(), "")
	note = strings.ReplaceAll(note, "\"", "`")

	return note, nil
}
