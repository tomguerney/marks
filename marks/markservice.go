package marks

type MarkService interface {
	Mark(id string) (*Mark, error)
	Marks() ([]*Mark, error)
	Create(m *Mark) error
	Update(id string, m *Mark) error
	Delete(id string) error
	Contains(id string) (bool, error)
	Filter(id, url string, tags []string) ([]*Mark, error)
}

type MarkAlreadyExistsError struct{}

func (e MarkAlreadyExistsError) Error() string {
	return "mark already exists"
}

type MarkDoesNotExistError struct{}

func (e MarkDoesNotExistError) Error() string {
	return "mark does not exist"
}
