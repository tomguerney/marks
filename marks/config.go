package marks

type Config struct {
	*AppConfig
	*UserConfig
}

type AppConfig struct {
	FullFormat        []string
	MarksYamlFileMode uint32
	SupportedBrowsers []string
	SupportedColors   []string
}

type UserConfig struct {
	ContentPath     string
	MarksYamlFile   string
	IdColor         string
	UrlColor        string
	TagsColor       string
	BrowserColor    string
	ChromeOpenArgs  string
	FirefoxOpenArgs string
	Browser         string
}
