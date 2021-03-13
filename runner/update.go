package runner

import (
	"github.com/abitofoldtom/marks/marks"
)

type update struct {
	*runner
	args *UpdateArgs
}

type UpdateArgs struct {
	id         string
	url        string
	tags       []string
	newId      string
	newUrl     string
	newTags    []string
	removeTags []string
	removeUrl  bool
}

func NewUpdateRunner(
	args *UpdateArgs,
	config *marks.Config,
	markService marks.MarkService,
	printer marks.Printer,
	prompter marks.Prompter,
) *update {
	return &update{
		newRunner(config, markService, printer, prompter),
		args,
	}
}

func NewUpdateArgs(id, url, newId, newUrl string, tags, newTags, removeTags []string, removeUrl bool) *UpdateArgs {
	return &UpdateArgs{
		id:         id,
		url:        url,
		tags:       tags,
		newId:      newId,
		newUrl:     newUrl,
		newTags:    newTags,
		removeTags: removeTags,
		removeUrl:  removeUrl,
	}
}

func (u *update) Run() error {

	selected, err := u.filter("Select bookmark to update", u.args.id, u.args.url, u.args.tags)

	if err != nil {
		if err, ok := err.(*runnerError); ok {
			u.printer.Error(err.Error())
			return nil
		} else {
			return err
		}
	}

	updated := &marks.Mark{
		Id:   u.updatedId(selected),
		Url:  u.updatedUrl(selected),
		Tags: u.updatedTags(selected),
	}

	if _, ok := u.containsRemoveTags(updated); !ok {
		u.printer.Error("Mark %v does not contain tag %v")
		return nil
	}

	updated.Tags = u.removeTags(updated.Tags)

	if u.args.removeUrl {
		updated.Url = ""
	}

	if err := u.markService.Update(selected.Id, updated); err != nil {
		return err
	}

	printSelected, err := u.printer.FullMark(selected)
	if err != nil {
		return err
	}

	printUpdated, err := u.printer.FullMark(updated)
	if err != nil {
		return err
	}

	u.printer.Msg("Mark updated from:\n%v\nto:\n%v", printSelected, printUpdated)

	return nil
}

func (u *update) updatedId(selected *marks.Mark) string {
	if u.args.newId != "" {
		return u.args.newId
	} else {
		return selected.Id
	}
}

func (u *update) updatedUrl(selected *marks.Mark) string {
	if u.args.newUrl != "" {
		return u.args.newUrl
	} else {
		return selected.Url
	}
}

func (u *update) updatedTags(selected *marks.Mark) []string {
	l := len(selected.Tags)
	c := len(selected.Tags) + len(u.args.newTags)
	selectedTags := make([]string, l, c)
	copy(selectedTags, selected.Tags)
	return append(selectedTags, u.args.newTags...)
}

func (u *update) containsRemoveTags(mark *marks.Mark) (tag string, ok bool) {
	for _, forRemoval := range u.args.removeTags {
		if !mark.ContainsTag(forRemoval) {
			return forRemoval, false
		}
	}
	return "", true
}

func (u *update) removeTag(tag string) bool {
	for _, forRemoval := range u.args.removeTags {
		if tag == forRemoval {
			return true
		}
	}
	return false
}

func (u *update) removeTags(tags []string) []string {
	remainingTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		if !u.removeTag(tag) {
			remainingTags = append(remainingTags, tag)
		}
	}
	return remainingTags
}
