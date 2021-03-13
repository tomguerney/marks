package runner

import (
	"github.com/abitofoldtom/marks/marks"
)

type add struct {
	args         *AddArgs
	config       *marks.Config
	marksService marks.MarkService
	printer      marks.Printer
}

type AddArgs struct {
	id   string
	url  string
	tags []string
}

func NewAddRunner(
	args *AddArgs,
	config *marks.Config,
	marks marks.MarkService,
	printer marks.Printer,
) *add {
	return &add{
		args,
		config,
		marks,
		printer,
	}
}

func NewAddArgs(id, url string, tags []string) *AddArgs {
	return &AddArgs{id, url, tags}
}

func (a *add) Run() error {

	exists, err := a.marksService.Contains(a.args.id)
	if err != nil {
		return err
	}

	if exists {
		a.printer.Error("Bookmark with id \"%v\" already exists", a.args.id)
		return nil
	}

	mark := &marks.Mark{
		Id:   a.args.id,
		Url:  a.args.url,
		Tags: a.args.tags,
	}

	err = a.marksService.Create(mark)
	if err != nil {
		return err
	}

	fullMark, err := a.printer.FullMark(mark)
	if err != nil {
		return err
	}

	a.printer.Msg("Bookmark created: %v", fullMark)

	return nil
}
