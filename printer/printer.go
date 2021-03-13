package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/abitofoldtom/marks/marks"
)

type printer struct {
	out       io.Writer
	config    *marks.Config
	colorizer colorizer
}

type colorizer interface {
	Colorize(string, string) (string, error)
	Black(string) string
	Red(string) string
	Green(string) string
	Yellow(string) string
	Blue(string) string
	Magenta(string) string
	Cyan(string) string
	White(string) string
}

func NewPrinter(config *marks.Config, colorizer colorizer) *printer {
	return &printer{os.Stdout, config, colorizer}
}

func (p *printer) Msg(text string, a ...interface{}) {
	textln := fmt.Sprintf("%s\n", text)
	fmt.Fprintf(p.out, textln, a...)
}

func (p *printer) Error(text string, a ...interface{}) {
	textln := fmt.Sprintf("Error: %s\n", text)
	fmt.Fprintf(p.out, textln, a...)
}

func (p *printer) Tabulate(mks []*marks.Mark) ([]string, error) {
	printMks, err := p.colorizeAll(mks)

	if err != nil {
		return nil, err
	}

	builder := strings.Builder{}
	writer := tabwriter.NewWriter(&builder, 0, 8, 4, ' ', 0)

	for _, printM := range printMks {
		fmt.Fprintln(writer, fmt.Sprintf("%s", strings.Join(printM.split(), "\t")))
	}

	writer.Flush()
	table := strings.Split(builder.String(), "\n")

	return table[:len(table)-1], nil
}

func (p *printer) FullMark(m *marks.Mark) (string, error) {
	pm, err := p.colorize(m)
	if err != nil {
		return "", err
	}

	output := []string{}

	output = append(output, fmt.Sprintf("%v", pm.id))

	if m.Url != "" {
		output = append(output, fmt.Sprintf("%v", pm.url))
	}

	if len(m.Tags) > 0 {
		output = append(output, fmt.Sprintf("%v", pm.tags))
	}

	return strings.Join(output, " "), nil
}

func (p *printer) FullMarkWithFields(m *marks.Mark) (string, error) {
	pm, err := p.colorize(m)
	if err != nil {
		return "", err
	}

	output := []string{}

	output = append(output, fmt.Sprintf("Id: %v", pm.id))

	if m.Url != "" {
		output = append(output, fmt.Sprintf("Url: %v", pm.url))
	}

	if len(m.Tags) > 0 {
		output = append(output, fmt.Sprintf("Tags: %v", pm.tags))
	}

	return strings.Join(output, ", "), nil
}

func (p *printer) colorizeAll(mks []*marks.Mark) (printMks []*printMark, err error) {
	for _, m := range mks {
		printM, err := p.colorize(m)
		if err != nil {
			return nil, err
		}
		printMks = append(printMks, printM)
	}
	return printMks, nil
}

func (p *printer) Id(id string) (string, error) {
	return p.colorizer.Colorize(p.config.IdColor, id)
}

func (p *printer) Url(url string) (string, error) {
	return p.colorizer.Colorize(p.config.UrlColor, url)
}

func (p *printer) Tags(tags []string) (string, error) {
	if len(tags) > 0 {
		bracketedTags := fmt.Sprintf("[%v]", strings.Join(tags, ", "))
		return p.colorizer.Colorize(p.config.TagsColor, bracketedTags)
	}
	return "", nil
}

func (p *printer) Browser(browser string) (string, error) {
	return p.colorizer.Colorize(p.config.BrowserColor, browser)
}

func (p *printer) colorize(m *marks.Mark) (*printMark, error) {
	id, err := p.Id(m.Id)
	if err != nil {
		return nil, err
	}
	url, err := p.Url(m.Url)
	if err != nil {
		return nil, err
	}
	tags, err := p.Tags(m.Tags)
	if err != nil {
		return nil, err
	}
	return &printMark{
		id:   id,
		url:  url,
		tags: tags,
	}, nil
}
