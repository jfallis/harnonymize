package harnonymize

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/chromedp/cdproto/har"
)

const (
	cookieName   = "cookie"
	authName     = "auth"
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

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", file.Path, file.Name))
	if err != nil {
		return err
	}

	return file.HAR.UnmarshalJSON(content)
}

func (c *Config) Write(file HAR) error {
	data, err := file.HAR.MarshalJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/anonymized_%s", file.Path, file.Name), data, fs.ModePerm)
}

func (c *Config) Anonymize(file HAR) {
	for i := range file.HAR.Log.Entries {
		entry := file.HAR.Log.Entries[i]

		entry.Request.Cookies = nil
		entry.Response.Cookies = nil

		redactHeaders(entry.Request.Headers)
		redactHeaders(entry.Response.Headers)

		c.redactContent(entry.Response.Content)
	}
}

func (c *Config) redactContent(content *har.Content) {
	for _, keyword := range c.BlockContentKeywords {
		if strings.Contains(strings.ToLower(content.Text), keyword) {
			content.Text = redactedText
			break
		}
	}
}

func redactHeaders(headers []*har.NameValuePair) {
	for i := range headers {
		header := headers[i]
		if strings.Contains(strings.ToLower(header.Name), authName) || strings.Contains(strings.ToLower(header.Name), cookieName) {
			header.Value = redactedText
		}
	}
}
