package config

import (
	"testing"

	"github.com/abitofoldtom/marks/mocks"
)

func TestBrowserMustBeSupportedPass(t *testing.T) {
	config := mocks.NewConfig()
	config.AppConfig.SupportedBrowsers = []string{"one", "two"}
	config.UserConfig.Browser = "two"
	err := browserMustBeSupported(config)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBrowserMustBeSupportedFail(t *testing.T) {
	config := mocks.NewConfig()
	config.AppConfig.SupportedBrowsers = []string{"one", "two"}
	config.UserConfig.Browser = "three"
	err := browserMustBeSupported(config)
	if err == nil {
		t.Fatal("Should cause error")
	}
}

func TestColorsMustBeSupportedPass(t *testing.T) {
	config := mocks.NewConfig()
	config.AppConfig.SupportedColors = []string{"one", "two"}
	config.UserConfig.IdColor = "one"
	config.UserConfig.UrlColor = "two"
	config.UserConfig.TagsColor = "one"
	config.UserConfig.BrowserColor = "two"
	err := colorsMustBeSupported(config)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestColorsMustBeSupportedFail(t *testing.T) {
	config := mocks.NewConfig()
	config.AppConfig.SupportedColors = []string{"one", "two"}
	config.UserConfig.IdColor = "one"
	config.UserConfig.UrlColor = "two"
	config.UserConfig.TagsColor = "three"
	config.UserConfig.BrowserColor = "two"
	err := colorsMustBeSupported(config)
	if err == nil {
		t.Fatal("Should cause error")
	}
}
