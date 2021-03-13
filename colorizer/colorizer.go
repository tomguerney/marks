package colorizer

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

type colorizer struct{}

func NewColorizer() *colorizer {
	return &colorizer{}
}

func (c *colorizer) Colorize(colorName string, text string) (string, error) {
	switch colorName {
	case "black":
		return c.Black(text), nil
	case "red":
		return c.Red(text), nil
	case "green":
		return c.Green(text), nil
	case "yellow":
		return c.Yellow(text), nil
	case "blue":
		return c.Blue(text), nil
	case "magenta":
		return c.Magenta(text), nil
	case "cyan":
		return c.Cyan(text), nil
	case "white":
		return c.White(text), nil
	}
	return "", errors.New(fmt.Sprintf("color \"%v\" not available", text))
}

func (c *colorizer) Black(text string) string {
	return color.BlackString(text)
}

func (c *colorizer) Red(text string) string {
	return color.RedString(text)
}

func (c *colorizer) Green(text string) string {
	return color.GreenString(text)
}

func (c *colorizer) Yellow(text string) string {
	return color.YellowString(text)
}

func (c *colorizer) Blue(text string) string {
	return color.BlueString(text)
}

func (c *colorizer) Magenta(text string) string {
	return color.MagentaString(text)
}

func (c *colorizer) Cyan(text string) string {
	return color.CyanString(text)
}

func (c *colorizer) White(text string) string {
	return color.WhiteString(text)
}
