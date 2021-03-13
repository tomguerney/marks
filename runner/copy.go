package runner

import (
	"github.com/abitofoldtom/marks/marks"
)

type copyRunner struct {
	*runner
	args    *CopyArgs
	clipper clipper
}

type CopyArgs struct {
	id   string
	url  string
	tags []string
}

type clipper interface {
	Copy(string) error
}

func NewCopyRunner(
	args *CopyArgs,
	config *marks.Config,
	marks marks.MarkService,
	printer marks.Printer,
	prompter marks.Prompter,
	clipper clipper,
) *copyRunner {
	return &copyRunner{
		&runner{
			config,
			marks,
			printer,
			prompter,
		},
		args,
		clipper,
	}
}

func NewCopyArgs(id, url string, tags []string) *CopyArgs {
	return &CopyArgs{id, url, tags}
}

func (c *copyRunner) Run() error {

	selected, err := c.filter("Select bookmark to copy", c.args.id, c.args.url, c.args.tags)

	if err != nil {
		if err, ok := err.(*runnerError); ok {
			c.printer.Error(err.Error())
			return nil
		} else {
			return err
		}
	}

	err = c.clipper.Copy(selected.Url)
	if err != nil {
		return err
	}

	printUrl, err := c.printer.Url(selected.Url)
	if err != nil {
		return err
	}

	c.printer.Msg("Url copied to clipboard: %v", printUrl)

	return nil
}
