package tests

import (
	"harnonymise/pkg/harnonymize"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	anon := harnonymize.New()
	anon.BlockContentKeywords = []string{"john fallis"}
	har := harnonymize.NewHAR("./data", "test.har")

	assert.NoError(t, anon.Read(har))

	anon.Anonymize(har)

	data, err := har.HAR.MarshalJSON()
	assert.NoError(t, err)

	harContent := string(data)

	for e := range har.HAR.Log.Entries {
		assert.Nil(t, har.HAR.Log.Entries[e].Request.Cookies)
		assert.Nil(t, har.HAR.Log.Entries[e].Response.Cookies)

		for h := range har.HAR.Log.Entries[e].Request.Headers {
			if har.HAR.Log.Entries[e].Request.Headers[h].Name != "cookie" {
				continue
			}
			assert.Equal(t, "** HARnonymize removed content **", har.HAR.Log.Entries[e].Request.Headers[h].Value)
			assert.Equal(t, "** HARnonymize removed content **", har.HAR.Log.Entries[e].Request.Headers[h].Value)
		}
	}

	assert.Contains(t, harContent, "** HARnonymize removed content **")
	assert.NotContains(t, harContent, "john fallis")
}
