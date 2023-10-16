package harnonymize

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/chromedp/cdproto/har"
)

const (
	cookieName = "cookie"
	authName   = "auth"
)

const (
	redactedText = "** HARnonymize removed content **"
)

type HAR struct {
	Path  string
	Name  string
	HAR   *har.HAR
	Error error
}

func NewHAR(path, name string) HAR {
	return HAR{
		Path: path,
		Name: name,
		HAR:  &har.HAR{},
	}
}

func (c *Config) Read(file HAR) error {
	if !strings.HasSuffix(file.Name, ".har") {
		return ErrNotHARFile
	}

	name := fmt.Sprintf("%s/%s", file.Path, file.Name)
	context, readErr := os.ReadFile(name)
	if readErr != nil {
		return readErr
	}

	return file.HAR.UnmarshalJSON(context)
}

func (c *Config) Write(file HAR) error {
	name := fmt.Sprintf("%s/anonymized_%s", file.Path, file.Name)

	data, err := file.HAR.MarshalJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(name, data, fs.ModePerm)
}

func (c *Config) Anonymize(file HAR) {
	for x := range file.HAR.Log.Entries {
		file.HAR.Log.Entries[x].Request.Cookies = nil
		file.HAR.Log.Entries[x].Response.Cookies = nil
		redactHeaders(file.HAR.Log.Entries[x].Request.Headers)
		redactHeaders(file.HAR.Log.Entries[x].Response.Headers)
		c.redactByContent(file.HAR.Log.Entries[x].Response.Content)
	}
}

func (c *Config) redactByContent(content *har.Content) {
	for _, keywords := range c.BlockContentKeywords {
		if stringContains(content.Text, keywords) {
			if strings.EqualFold(content.MimeType, "application/json") {
				content.Text = fmt.Sprintf(`{text: %q}`, redactedText)
				break
			}
			content.Text = redactedText
			break
		}
	}
}

func redactHeaders(headers []*har.NameValuePair) {
	for y := range headers {
		if stringContains(headers[y].Name, authName) {
			headers[y].Value = redactedText
		}
		if stringContains(headers[y].Name, cookieName) {
			headers[y].Value = redactedText
		}
	}
}

func stringContains(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), substr)
}
