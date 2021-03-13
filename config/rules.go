package config

import (
	"errors"
	"fmt"

	"github.com/abitofoldtom/marks/marks"
)

var browserMustBeSupported = func(c *marks.Config) error {
	browser := c.UserConfig.Browser
	for _, supportedBrowser := range c.AppConfig.SupportedBrowsers {
		if browser == supportedBrowser {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%v is not a supported browser", browser))
}

var colorsMustBeSupported = func(c *marks.Config) error {

	colorSupported := func(userColor string) bool {
		for _, supportedColor := range c.AppConfig.SupportedColors {
			if userColor == supportedColor {
				return true
			}
		}
		return false
	}

	allColorsSupported := func(userColors []string) (bool, string) {
		for _, userColor := range userColors {
			if !colorSupported(userColor) {
				return false, userColor
			}
		}
		return true, ""
	}

	userColors := []string{c.UserConfig.IdColor, c.UserConfig.UrlColor, c.UserConfig.TagsColor, c.UserConfig.BrowserColor}

	if ok, color := allColorsSupported(userColors); ok {
		return nil
	} else {
		return errors.New(fmt.Sprintf("%v is not a supported color", color))
	}
}
