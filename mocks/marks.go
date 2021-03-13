package mocks

import "github.com/abitofoldtom/marks/marks"

var DefaultMarks = []*marks.Mark{
	&marks.Mark{
		Id:   "Abc News",
		Url:  "https://www.abc.net.au/news/",
		Tags: []string{"news", "current affairs"},
	},
	&marks.Mark{
		Id:   "Google",
		Url:  "https://www.google.com",
		Tags: []string{"search"},
	},
	&marks.Mark{
		Id:   "BBC News",
		Url:  "https://www.bbc.com/news",
		Tags: []string{"news", "uk"},
	},
}
