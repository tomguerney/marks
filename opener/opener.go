package opener

import (
	"errors"
	"fmt"
	"html/template"
	"os/exec"
	"strings"

	"github.com/abitofoldtom/marks/marks"
	"github.com/apex/log"
	"github.com/mattn/go-shellwords"
)

type opener struct {
	config    *marks.Config
	commander commander
}

type commander interface {
	Command(name string, arg ...string) combinedOutputter
}

type combinedOutputter interface {
	CombinedOutput() ([]byte, error)
}

type concreteCommander struct{}

func (c *concreteCommander) Command(name string, arg ...string) combinedOutputter {
	return exec.Command(name, arg...)
}

func NewOpener(config *marks.Config) *opener {
	return &opener{config: config, commander: &concreteCommander{}}
}

func (o *opener) Open(url, browser string) error {
	argTemplate, err := o.template(browser)
	if err != nil {
		return err
	}

	argString, err := o.interpolateTemplate(argTemplate, url)
	if err != nil {
		return err
	}

	argSlice, err := o.sliceArgs(argString)
	if err != nil {
		return err
	}

	cmd := o.commander.Command("open", argSlice...)

	out, err := cmd.CombinedOutput()

	log.Infof("Open output: %v", string(out))

	if err != nil {
		return err
	}
	return nil
}

func (o *opener) template(browser string) (string, error) {
	switch browser {
	case "chrome":
		return o.config.ChromeOpenArgs, nil
	case "firefox":
		return o.config.FirefoxOpenArgs, nil
	default:
		return "", errors.New(fmt.Sprintf("Browser \"%v\" not supported\n", browser))
	}

}

func (o *opener) interpolateTemplate(argTemplate, url string) (string, error) {
	builder := strings.Builder{}
	tmpl, err := template.New("open").Parse(argTemplate)
	if err != nil {
		return "", err
	}
	data := struct{ Url string }{Url: url}
	err = tmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (o *opener) sliceArgs(args string) ([]string, error) {
	slice, err := shellwords.Parse(args)
	if err != nil {
		return nil, err
	}
	return slice, nil
}
