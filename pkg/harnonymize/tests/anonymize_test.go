package tests

import (
	"harnonymise/pkg/harnonymize"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHAR(t *testing.T) {
	t.Parallel()

	har := harnonymize.NewHAR("./data", "test.har")
	assert.Equal(t, "./data", har.Path)
	assert.Equal(t, "test.har", har.Name)
}

func TestReadFileDoesNotExist(t *testing.T) {
	t.Parallel()

	anon := harnonymize.New()
	har := harnonymize.NewHAR("./data", "missing.har")

	assert.Error(t, anon.Read(har))
}

func TestReadInvalidFileType(t *testing.T) {
	t.Parallel()

	anon := harnonymize.New()
	har := harnonymize.NewHAR("./data", "test")

	assert.ErrorIs(t, anon.Read(har), harnonymize.ErrNotHARFile)
}

func TestReadEmptyFile(t *testing.T) {
	t.Parallel()

	anon := harnonymize.New()
	har := harnonymize.NewHAR("./data", "empty.har")

	assert.ErrorIs(t, anon.Read(har), harnonymize.ErrNotHARFile)
}

func TestAnonymizeAndWrite(t *testing.T) {
	t.Parallel()

	anon := harnonymize.New()
	anon.BlockContentKeywords = []string{"john fallis"}
	har := harnonymize.NewHAR("./data", "test.har")

	assert.NoError(t, anon.Read(har))

	anon.Anonymize(har)
	assert.NoError(t, anon.Write(har))
}

func TestAnonymizeRemovesCookies(t *testing.T) {
	t.Parallel()

	anon := harnonymize.New()
	har := harnonymize.NewHAR("./data", "test.har")

	assert.NoError(t, anon.Read(har))

	anon.Anonymize(har)

	for e := range har.HAR.Log.Entries {
		assert.Nil(t, har.HAR.Log.Entries[e].Request.Cookies)
		assert.Nil(t, har.HAR.Log.Entries[e].Response.Cookies)
	}
}

func TestAnonymizeRedactsHeaders(t *testing.T) {
	t.Parallel()

	anon := harnonymize.New()
	har := harnonymize.NewHAR("./data", "test.har")

	assert.NoError(t, anon.Read(har))

	anon.Anonymize(har)

	for e := range har.HAR.Log.Entries {
		for h := range har.HAR.Log.Entries[e].Request.Headers {
			if har.HAR.Log.Entries[e].Request.Headers[h].Name != "cookie" {
				continue
			}
			assert.Equal(t, "** HARnonymize removed content **", har.HAR.Log.Entries[e].Request.Headers[h].Value)
		}
	}
}

func TestAnonymizeRedactsContent(t *testing.T) {
	t.Parallel()

	anon := harnonymize.New()
	anon.BlockContentKeywords = []string{"john fallis"}
	har := harnonymize.NewHAR("./data", "test.har")

	assert.NoError(t, anon.Read(har))

	anon.Anonymize(har)

	data, err := har.HAR.MarshalJSON()
	assert.NoError(t, err)

	harContent := string(data)

	assert.Contains(t, harContent, "** HARnonymize removed content **")
	assert.NotContains(t, harContent, "john fallis")
}
