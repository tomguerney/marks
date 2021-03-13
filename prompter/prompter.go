package prompter

import (
	"github.com/manifoldco/promptui"
)

type prompter struct{}

func NewPrompter() *prompter {
	return &prompter{}
}

func (p *prompter) Select(label string, table []string) (i int, err error) {

	prompt := promptui.Select{
		Label: label,
		Items: table,
	}

	i, _, err = prompt.Run()

	if err != nil {
		return 0, err
	}

	return i, nil
}

func (p *prompter) Confirm(label string) bool {

	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()

	if err != nil {
		return false
	}

	return true
}
