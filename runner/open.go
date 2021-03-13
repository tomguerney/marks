package runner

import (
	"github.com/abitofoldtom/marks/marks"
)

type open struct {
	*runner
	args   *OpenArgs
	opener opener
}

type OpenArgs struct {
	id   string
	url  string
	tags []string
}

type opener interface {
	Open(url, browser string) error
}

func NewOpenRunner(
	args *OpenArgs,
	config *marks.Config,
	markService marks.MarkService,
	printer marks.Printer,
	prompter marks.Prompter,
	opener opener,
) *open {
	return &open{
		newRunner(config, markService, printer, prompter),
		args,
		opener,
	}
}

func NewOpenArgs(id, url string, tags []string) *OpenArgs {
	return &OpenArgs{id, url, tags}
}

func (o *open) Run() error {

	selected, err := o.filter("Select bookmark to open", o.args.id, o.args.url, o.args.tags)

	if err != nil {
		if err, ok := err.(*runnerError); ok {
			o.printer.Error(err.Error())
			return nil
		} else {
			return err
		}
	}

	err = o.opener.Open(selected.Url, o.config.Browser)
	if err != nil {
		return err
	}

	printUrl, err := o.printer.Url(selected.Url)
	if err != nil {
		return err
	}

	printBrowser, err := o.printer.Browser(o.config.Browser)
	if err != nil {
		return err
	}

	o.printer.Msg("Url opened in %v: %v", printBrowser, printUrl)

	return nil
}
