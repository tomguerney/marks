package mocks

import "github.com/tomguerney/marks/marks"

func NewConfig() *marks.Config {
	return &marks.Config{
		AppConfig:  &marks.AppConfig{},
		UserConfig: &marks.UserConfig{},
	}
}
