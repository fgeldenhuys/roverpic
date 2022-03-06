package roverapi

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	Scheme    string
	Host      string
	APIKey    string
	Endpoints Endpoints
}

type Endpoint struct {
	Path   string
	Method string
}

type Endpoints struct {
	Photos Endpoint
}

var Defaults = Config{
	Endpoints: Endpoints{
		Photos: Endpoint{
			Path:   "mars-photos/api/v1/rovers/curiosity/photos",
			Method: "GET",
		},
	},
}

func (conf Config) Validate() error {
	if conf.Scheme == "" {
		return fmt.Errorf("Scheme is required")
	}
	if conf.Host == "" {
		return fmt.Errorf("Host is required")
	}
	if conf.APIKey == "" {
		return fmt.Errorf("APIKey is required")
	}
	return nil
}

func (conf *Config) PhotosURL(date time.Time) string {
	u := url.URL{
		Scheme: conf.Scheme,
		Host:   conf.Host,
		Path:   conf.Endpoints.Photos.Path,
	}
	q := u.Query()
	q.Set("api_key", conf.APIKey)
	q.Set("earth_date", date.Format("2006-01-02"))
	u.RawQuery = q.Encode()
	return u.String()
}
