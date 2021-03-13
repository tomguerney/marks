package clipper

import (
	"github.com/atotto/clipboard"
)

type clipper struct{}

func NewClipper() *clipper {
	return &clipper{}
}

func (c *clipper) Copy(text string) error {
	return clipboard.WriteAll(text)
}
	