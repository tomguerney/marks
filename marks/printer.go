package marks

type Printer interface {
	Msg(string, ...interface{})
	Error(string, ...interface{})
	Tabulate([]*Mark) ([]string, error)
	FullMark(*Mark) (string, error)
	FullMarkWithFields(*Mark) (string, error)
	Id(string) (string, error)
	Url(string) (string, error)
	Tags([]string) (string, error)
	Browser(string) (string, error)
}
