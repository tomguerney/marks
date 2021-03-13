package config

import (
	"strings"

	"github.com/abitofoldtom/marks/marks"
)

type loader struct {
	provider
	validator
}

type provider interface {
	GetString(string) string
	SetDefault(string, interface{})
}

type validator interface {
	validate(*marks.Config, []validationRule) error
}

type concreteValidator struct{}

type validationRule func(*marks.Config) error

func NewLoader(provider provider) *loader {
	return &loader{provider, &concreteValidator{}}
}

func (l *loader) Load() (*marks.Config, error) {
	l.setDefaults()
	config := l.loadConfig()
	rules := l.loadValidationRules()
	if err := l.validate(config, rules); err != nil {
		return nil, err
	}
	return config, nil
}

func (l *loader) loadConfig() *marks.Config {
	userConfig := l.loadUserConfig()
	appConfig := l.loadAppConfg()
	return &marks.Config{AppConfig: appConfig, UserConfig: userConfig}
}

func (l *loader) setDefaults() {
	l.SetDefault("yaml", "bookmarks.yaml")
	l.SetDefault("idColor", "green")
	l.SetDefault("urlColor", "blue")
	l.SetDefault("tagsColor", "yellow")
	l.SetDefault("browserColor", "red")
	l.SetDefault("chromeOpenArgs", "-a \"Google Chrome\" {{.Url}}")
	l.SetDefault("firefoxOpenargs", "-a firefox {{.Url}}")
	l.SetDefault("browser", "chrome")
}

func (l *loader) loadUserConfig() *marks.UserConfig {
	return &marks.UserConfig{
		ContentPath:     l.GetString("contentpath"),
		MarksYamlFile:   l.GetString("yaml"),
		ChromeOpenArgs:  l.GetString("chromeOpenArgs"),
		FirefoxOpenArgs: l.GetString("firefoxOpenargs"),
		IdColor:         strings.ToLower(l.GetString("idColor")),
		UrlColor:        strings.ToLower(l.GetString("urlColor")),
		TagsColor:       strings.ToLower(l.GetString("tagsColor")),
		BrowserColor:    strings.ToLower(l.GetString("browserColor")),
		Browser:         strings.ToLower(l.GetString("browser")),
	}
}

func (l *loader) loadAppConfg() *marks.AppConfig {
	return &marks.AppConfig{
		FullFormat:        []string{"id", "url", "tags"},
		MarksYamlFileMode: 0644,
		SupportedBrowsers: []string{"chrome", "firefox"},
		SupportedColors:   []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"},
	}
}

func (l *loader) loadValidationRules() []validationRule {
	return []validationRule{
		browserMustBeSupported,
		colorsMustBeSupported,
	}
}

func (v *concreteValidator) validate(config *marks.Config, rules []validationRule) error {
	for _, rule := range rules {
		if err := rule(config); err != nil {
			return err
		}
	}
	return nil
}
