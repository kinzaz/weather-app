package geo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type GeoData struct {
	City string `json:"city`
}

type CityPopulationResponse struct {
	Error bool `json:"error"`
}

func GetMyLocation(city string) (*GeoData, error) {
	if city != "" {
		isCity := checkCity(city)
		if !isCity {
			panic("Такого города нет")
		}
		return &GeoData{
			City: city,
		}, nil
	}
	response, err := http.Get("https://ipapi.co/json/")
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("NOT200")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var geo GeoData
	json.Unmarshal(body, &geo)
	return &geo, nil
}

func checkCity(city string) bool {
	postBody, _ := json.Marshal(map[string]string{"city": city})
	response, err := http.Post("https://countriesnow.space/api/v0.1/countries/population/cities", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		return false
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false
	}
	var populationResponse CityPopulationResponse
	json.Unmarshal(body, &populationResponse)

	return !populationResponse.Error
}
