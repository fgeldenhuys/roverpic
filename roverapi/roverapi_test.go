package roverapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPhotos(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"photos":[{"id":102685,"sol":1004,"camera":{"id":20,"name":"FHAZ","rover_id":5,"full_name":"Front Hazard Avoidance Camera"},"img_src":"http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/01004/opgs/edr/fcam/FLB_486615455EDR_F0481570FHAZ00323M_.JPG","earth_date":"2015-06-03","rover":{"id":5,"name":"Curiosity","landing_date":"2012-08-06","launch_date":"2011-11-26","status":"active"}}]}`)
	}))
	defer srv.Close()
	roverAPI := Init(testConfig())
	roverAPI.Config.Host = srv.URL[7:]

	photos, err := roverAPI.GetPhotos(time.Date(2020, time.May, 11, 0, 0, 0, 0, time.UTC))

	assert.NoError(t, err)
	assert.Len(t, photos, 1)
	assert.Equal(t, "http://mars.jpl.nasa.gov/msl-raw-images/proj/msl/redops/ods/surface/sol/01004/opgs/edr/fcam/FLB_486615455EDR_F0481570FHAZ00323M_.JPG", photos[0].ImgSrc)
}

func TestGetPhotosBadCode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	roverAPI := Init(testConfig())
	roverAPI.Config.Host = srv.URL[7:]

	_, err := roverAPI.GetPhotos(time.Date(2020, time.May, 11, 0, 0, 0, 0, time.UTC))

	assert.Error(t, err)
}
