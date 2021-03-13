package marks

type Prompter interface {
	Select(string, []string) (int, error)
	Confirm(string) bool
}
