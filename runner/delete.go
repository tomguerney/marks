package runner

import (
	"github.com/abitofoldtom/marks/marks"
)

type deleteRunner struct {
	*runner
	args *DeleteArgs
}

type DeleteArgs struct {
	id   string
	url  string
	tags []string
}

func NewDeleteRunner(
	args *DeleteArgs,
	config *marks.Config,
	marks marks.MarkService,
	printer marks.Printer,
	prompter marks.Prompter,
) *deleteRunner {
	return &deleteRunner{
		&runner{
			config,
			marks,
			printer,
			prompter,
		},
		args,
	}
}

func NewDeleteArgs(id, url string, tags []string) *DeleteArgs {
	return &DeleteArgs{id, url, tags}
}

func (d *deleteRunner) Run() error {

	selected, err := d.filter("Select bookmark to delete", d.args.id, d.args.url, d.args.tags)

	if err != nil {
		if err, ok := err.(*runnerError); ok {
			d.printer.Error(err.Error())
			return nil
		} else {
			return err
		}
	}

	printSelected, err := d.printer.FullMark(selected)
	if err != nil {
		return err
	}

	d.printer.Msg("Selected: %v", printSelected)

	if !d.prompter.Confirm("Are sure you want to delete?") {
		d.printer.Msg("Exiting")
		return nil
	}

	err = d.markService.Delete(selected.Id)
	if err != nil {
		return err
	}

	d.printer.Msg("Deleted")

	return nil
}
