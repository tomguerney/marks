package printer

type printMark struct {
	id   string
	url  string
	tags string
}

func (m *printMark) split() []string {
	return []string{m.id, m.url, m.tags}
}
