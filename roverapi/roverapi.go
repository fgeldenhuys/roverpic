package roverapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RoverAPI struct {
	Config
}

func Init(conf Config) *RoverAPI {
	return &RoverAPI{
		Config: conf,
	}
}

type Photo struct {
	ID     int `json:"id"`
	Sol    int `json:"sol"`
	Camera struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		RoverID  int    `json:"rover_id"`
		FullName string `json:"full_name"`
	} `json:"camera"`
	ImgSrc    string `json:"img_src"`
	EarthDate string `json:"earth_date"`
	Rover     struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		LandingDate string `json:"landing_date"`
		LaunchDate  string `json:"launch_date"`
		Status      string `json:"status"`
	} `json:"rover"`
}

func (api *RoverAPI) GetPhotos(date time.Time) ([]Photo, error) {
	url := api.Config.PhotosURL(date)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to query API: %w", err)
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("Failed to query API: no body returned in response")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to query API: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body during API query: %w", err)
	}

	var photos struct {
		Photos []Photo
	}
	err = json.Unmarshal(body, &photos)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse response body during API query: %w", err)
	}

	return photos.Photos, nil
}
