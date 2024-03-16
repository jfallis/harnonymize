package harnonymize

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

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

	content, rErr := os.ReadFile(fmt.Sprintf("%s/%s", file.Path, file.Name))
	if rErr != nil {
		return rErr
	}

	if err := file.HAR.UnmarshalJSON(content); err != nil {
		return err
	}

	if file.HAR == nil || file.HAR.Log == nil {
		return ErrNotHARFile
	}

	return nil
}

func (c *Config) Write(file HAR) error {
	data, err := file.HAR.MarshalJSON()
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/anonymized_%s", file.Path, file.Name))
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(string(data))
	if err != nil {
		return err
	}

	return writer.Flush()
}

func (c *Config) Anonymize(file HAR) {
	var wg sync.WaitGroup
	for i := range file.HAR.Log.Entries {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			entry := file.HAR.Log.Entries[i]

			entry.Request.Cookies = nil
			entry.Response.Cookies = nil

			redactHeaders(entry.Request.Headers)
			redactHeaders(entry.Response.Headers)
			c.redactByContent(entry.Response.Content)
		}(i)
	}
	wg.Wait()
}

func (c *Config) redactByContent(content *har.Content) {
	lowerContent := strings.ToLower(content.Text)
	for _, keyword := range c.BlockContentKeywords {
		if strings.Contains(lowerContent, keyword) {
			content.Text = redactedText
			break
		}
	}
}

func redactHeaders(headers []*har.NameValuePair) {
	for i := range headers {
		header := headers[i]
		lowerHeaderName := strings.ToLower(header.Name)
		if strings.Contains(lowerHeaderName, authName) || strings.Contains(lowerHeaderName, cookieName) {
			header.Value = redactedText
		}
	}
}
