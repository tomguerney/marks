package runner

import (
	"errors"
	"fmt"

	"github.com/abitofoldtom/marks/marks"
)

type runner struct {
	config      *marks.Config
	markService marks.MarkService
	printer     marks.Printer
	prompter    marks.Prompter
}

func newRunner(config *marks.Config, markService marks.MarkService, printer marks.Printer, prompter marks.Prompter) *runner {
	return &runner{
		config:      config,
		markService: markService,
		printer:     printer,
		prompter:    prompter,
	}
}

func (r *runner) filter(prompt, id, url string, tags []string) (*marks.Mark, error) {

	filtered, err := r.markService.Filter(id, url, tags)
	if err != nil {
		return nil, err
	}

	if len(filtered) == 0 {

		filterMark := &marks.Mark{Id: id, Url: url, Tags: tags}

		printFilterMark, err := r.printer.FullMarkWithFields(filterMark)
		if err != nil {
			return nil, err
		}

		return nil, &runnerError{fmt.Sprintf("No bookmarks found matching: %v", printFilterMark)}
	}

	var i int

	if len(filtered) == 1 {
		i = 0
	} else {
		table, err := r.printer.Tabulate(filtered)
		if err != nil {
			return nil, err
		}
		i, err = r.prompter.Select(prompt, table)
		if err != nil {
			return nil, err
		}
	}

	if i >= len(filtered) {
		return nil, errors.New(fmt.Sprintf("no mark at index %v", i))
	}

	return filtered[i], nil
}
