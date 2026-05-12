package superpanel

import (
	"regexp"
	"strings"
)

type ContentType string

const (
	TypeText  ContentType = "text"
	TypeURL   ContentType = "url"
	TypeFile  ContentType = "file"
	TypeImage ContentType = "image"
	TypeCode  ContentType = "code"
)

type Detection struct {
	Type    ContentType `json:"type"`
	Content string      `json:"content"`
	Actions []Action    `json:"actions"`
}

type Action struct {
	Label string `json:"label"`
	Code  string `json:"code"`
}

var (
	urlPattern  = regexp.MustCompile(`^https?://[^\s]+$`)
	codePattern = regexp.MustCompile(`(?s)^(?:function|class|import|def|const|var|package)\b`)
)

func Detect(content string) Detection {
	trimmed := strings.TrimSpace(content)
	detection := Detection{Content: trimmed}

	if urlPattern.MatchString(trimmed) {
		detection.Type = TypeURL
		detection.Actions = []Action{
			{Label: "Open in Browser", Code: "open_url"},
			{Label: "Copy URL", Code: "copy"},
		}
		return detection
	}

	if len(trimmed) > 100 && codePattern.MatchString(trimmed) {
		detection.Type = TypeCode
		detection.Actions = []Action{
			{Label: "Copy Code", Code: "copy"},
		}
		return detection
	}

	detection.Type = TypeText
	detection.Actions = []Action{
		{Label: "Search", Code: "search"},
		{Label: "Copy", Code: "copy"},
	}
	return detection
}
