package roverapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testConfig() Config {
	conf := Defaults
	conf.Scheme = "http"
	conf.Host = "test.server"
	conf.APIKey = "test-api-key"
	return conf
}

func TestConfigValidate(t *testing.T) {
	conf := testConfig()
	assert.NoError(t, conf.Validate())
}

func TestConfigValidateMissingScheme(t *testing.T) {
	conf := testConfig()
	conf.Scheme = ""
	assert.Error(t, conf.Validate())
}

func TestConfigValidateMissingHost(t *testing.T) {
	conf := testConfig()
	conf.Host = ""
	assert.Error(t, conf.Validate())
}

func TestConfigValidateMissingAPIKey(t *testing.T) {
	conf := testConfig()
	conf.APIKey = ""
	assert.Error(t, conf.Validate())
}

func TestPhotosURL(t *testing.T) {
	conf := testConfig()
	url := conf.PhotosURL(time.Date(2020, time.May, 11, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, "http://test.server/mars-photos/api/v1/rovers/curiosity/photos?api_key=test-api-key&earth_date=2020-05-11", url)
}
