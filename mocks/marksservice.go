package mocks

import "github.com/abitofoldtom/marks/marks"

type MarkService struct {
	MarkFn           func(id string) (*marks.Mark, error)
	MarksFn          func() ([]*marks.Mark, error)
	CreateFn         func(m *marks.Mark) error
	UpdateFn         func(id string, m *marks.Mark) error
	DeleteFn         func(id string) error
	ContainsFn       func(id string) (bool, error)
	FilterFn         func(id, url string, tags []string) ([]*marks.Mark, error)
	MarkFnCalled     bool
	MarksFnCalled    bool
	CreateFnCalled   bool
	UpdateFnCalled   bool
	DeleteFnCalled   bool
	ContainsFnCalled bool
	FilterFnCalled   bool
}

func NewMarkService() *MarkService {
	return &MarkService{
		MarkFn:     defaultMarkFn,
		MarksFn:    defaultMarksFn,
		CreateFn:   defaultCreateFn,
		UpdateFn:   defaultUpdateFn,
		DeleteFn:   defaultDeleteFn,
		ContainsFn: defaultContainsFn,
		FilterFn:   defaultFilterFn,
	}
}

var defaultMarkFn = func(id string) (*marks.Mark, error) {
	return DefaultMarks[0], nil
}

var defaultMarksFn = func() ([]*marks.Mark, error) {
	return DefaultMarks, nil
}

var defaultCreateFn = func(m *marks.Mark) error {
	return nil
}

var defaultUpdateFn = func(id string, m *marks.Mark) error {
	return nil
}

var defaultDeleteFn = func(id string) error {
	return nil
}

var defaultContainsFn = func(id string) (bool, error) {
	return false, nil
}

var defaultFilterFn = func(id, url string, tags []string) ([]*marks.Mark, error) {
	return []*marks.Mark{DefaultMarks[0], DefaultMarks[1]}, nil
}

func (s *MarkService) Mark(id string) (*marks.Mark, error) {
	s.MarkFnCalled = true
	return s.MarkFn(id)
}

func (s *MarkService) Marks() ([]*marks.Mark, error) {
	s.MarksFnCalled = true
	return s.MarksFn()
}

func (s *MarkService) Create(m *marks.Mark) error {
	s.CreateFnCalled = true
	return s.CreateFn(m)
}

func (s *MarkService) Update(id string, m *marks.Mark) error {
	s.UpdateFnCalled = true
	return s.UpdateFn(id, m)
}

func (s *MarkService) Delete(id string) error {
	s.DeleteFnCalled = true
	return s.DeleteFn(id)
}

func (s *MarkService) Contains(id string) (bool, error) {
	s.ContainsFnCalled = true
	return s.ContainsFn(id)
}

func (s *MarkService) Filter(id, url string, tags []string) ([]*marks.Mark, error) {
	s.FilterFnCalled = true
	return s.FilterFn(id, url, tags)
}
