package yaml

import (
	"path"
	"strings"

	"github.com/abitofoldtom/marks/marks"
	"gopkg.in/yaml.v2"
)

type markService struct {
	config       *marks.Config
	readerWriter ReaderWriter
}

type ReaderWriter interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte, uint32) error
}

func NewMarkService(config *marks.Config, readerWriter ReaderWriter) *markService {
	return &markService{config, readerWriter}
}

func (s *markService) Mark(id string) (*marks.Mark, error) {
	marks, err := s.loadMarks()
	if err != nil {
		return nil, err
	}
	for _, mark := range marks {
		if strings.ToLower(mark.Id) == strings.ToLower(id) {
			return mark, nil
		}
	}
	return nil, nil
}

func (s *markService) Marks() ([]*marks.Mark, error) {
	return s.loadMarks()
}

func (s *markService) Create(m *marks.Mark) error {
	exists, err := s.Contains(m.Id)
	if err != nil {
		return err
	}
	if exists {
		return marks.MarkAlreadyExistsError{}
	}
	marks, err := s.loadMarks()
	if err != nil {
		return err
	}
	marks = append(marks, m)
	err = s.saveMarks(marks)
	if err != nil {
		return err
	}
	return nil
}

func (s *markService) Update(id string, new *marks.Mark) error {
	updateFn := func(i int, marks []*marks.Mark) []*marks.Mark {
		marks[i] = new
		return marks
	}
	return s.modify(id, updateFn)
}

func (s *markService) Delete(id string) error {
	deleteFn := func(i int, marks []*marks.Mark) []*marks.Mark {
		return append(marks[:i], marks[i+1:]...)
	}
	return s.modify(id, deleteFn)
}

func (s *markService) Contains(id string) (bool, error) {
	marks, err := s.loadMarks()
	if err != nil {
		return false, err
	}
	for _, mark := range marks {
		if strings.ToLower(mark.Id) == strings.ToLower(id) {
			return true, nil
		}
	}
	return false, nil
}

func (s *markService) modify(id string, modifyFn func(int, []*marks.Mark) []*marks.Mark) error {
	exists, err := s.Contains(id)
	if err != nil {
		return err
	}
	if !exists {
		return marks.MarkDoesNotExistError{}
	}
	marks, err := s.loadMarks()
	if err != nil {
		return err
	}
	for i, mark := range marks {
		if strings.ToLower(mark.Id) == strings.ToLower(id) {
			marks = modifyFn(i, marks)
			break
		}
	}
	return s.saveMarks(marks)
}

func (s *markService) Filter(id, url string, tags []string) ([]*marks.Mark, error) {
	filtered, err := s.loadMarks()
	if err != nil {
		return nil, err
	}
	filtered = filterId(filtered, id)
	filtered = filterUrl(filtered, url)
	filtered = filterTags(filtered, tags)
	return filtered, nil
}

func filterId(unfiltered []*marks.Mark, id string) (filtered []*marks.Mark) {
	if id == "" {
		return unfiltered
	}
	for _, mark := range unfiltered {
		if strings.Contains(strings.ToLower(mark.Id), strings.ToLower(id)) {
			filtered = append(filtered, mark)
		}
	}
	return filtered
}

func filterUrl(unfiltered []*marks.Mark, url string) (filtered []*marks.Mark) {
	if url == "" {
		return unfiltered
	}
	for _, mark := range unfiltered {
		if strings.Contains(strings.ToLower(mark.Url), strings.ToLower(url)) {
			filtered = append(filtered, mark)
		}
	}
	return filtered
}

func filterTags(unfiltered []*marks.Mark, tags []string) (filtered []*marks.Mark) {
	if len(tags) == 0 {
		return unfiltered
	}
	for _, mark := range unfiltered {
		if mark.ContainsAllTags(tags) {
			filtered = append(filtered, mark)
		}
	}
	return filtered
}

func (s *markService) loadMarks() ([]*marks.Mark, error) {
	marksYaml, err := s.readerWriter.ReadFile(s.yamlPath())
	if err != nil {
		return nil, err
	}
	marks := []*marks.Mark{}
	if err := yaml.Unmarshal(marksYaml, &marks); err != nil {
		return nil, err
	}
	return marks, nil
}

func (s *markService) saveMarks(marks []*marks.Mark) error {
	marksYaml, err := yaml.Marshal(marks)
	if err != nil {
		return err
	}
	if err := s.readerWriter.WriteFile(s.yamlPath(), marksYaml, s.config.MarksYamlFileMode); err != nil {
		return err
	}
	return nil
}

func (s *markService) yamlPath() string {
	return path.Join(s.config.ContentPath, s.config.MarksYamlFile)
}
