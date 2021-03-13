package marks

import (
	"fmt"
	"strings"
)

type Mark struct {
	Id   string
	Url  string
	Tags []string
}

func (m *Mark) ContainsAllTags(subtags []string) bool {
	for _, subtag := range subtags {
		if !m.ContainsTag(subtag) {
			return false
		}
	}
	return true
}

func (m *Mark) ContainsTag(subtag string) bool {
	for _, tag := range m.Tags {
		if strings.ToLower(subtag) == strings.ToLower(tag) {
			return true
		}
	}
	return false
}

func (m *Mark) String() string {
	return fmt.Sprintf("Id: %v, Url: %v, Tags: [%v]", m.Id, m.Url, strings.Join(m.Tags, ", "))
}
