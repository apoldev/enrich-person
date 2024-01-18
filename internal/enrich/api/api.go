package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Service обслуживает запросы к api
// Поскольку api оказались идентичными, то не требуется разделять логику запроса к api в этой реализации
type Service struct {
	// Транспорт общего клиента сможет переиспользовать запросы
	Client *http.Client
}

const (
	ErrApiLimit = "Request limit reached"

	BaseUrlAge         = "https://api.agify.io/?name="
	BaseUrlGender      = "https://api.genderize.io/?name="
	BaseUrlNationality = "https://api.nationalize.io/?name="
)

var (
	ErrNoData    = errors.New("no data")
	ErrDataLimit = errors.New("limit error")
)

func (s *Service) loadData(url string, data interface{}) (err error) {

	resp, err := s.Client.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	return decoder.Decode(&data)
}

func (s *Service) GetAge(name string) (age int, err error) {

	var data ageModel

	err = s.loadData(BaseUrlAge+url.QueryEscape(name), &data)
	if err != nil {
		return 0, err
	}

	if data.Error == ErrApiLimit {
		return 0, ErrDataLimit
	}

	return data.Age, err

}

func (s *Service) GetGender(name string) (gender string, err error) {

	var data genderModel

	err = s.loadData(BaseUrlGender+url.QueryEscape(name), &data)
	if err != nil {
		return "", err
	}

	if data.Error == ErrApiLimit {
		return "", ErrDataLimit
	}

	return data.Gender, nil

}

func (s *Service) GetNationality(name string) (nationality string, err error) {

	var data nationalityModel

	err = s.loadData(BaseUrlNationality+url.QueryEscape(name), &data)
	if err != nil {
		return
	}

	if data.Error == ErrApiLimit {
		return "", ErrDataLimit
	}

	if len(data.Country) == 0 {
		return "", ErrNoData
	}

	return data.Country[0].CountryID, nil

}

type ageModel struct {
	Age   int    `json:"age"`
	Error string `json:"error"`
}

type genderModel struct {
	Gender string `json:"gender"`
	Error  string `json:"error"`
}

type nationalityModelCountry struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type nationalityModel struct {
	Country []nationalityModelCountry `json:"country"`
	Error   string                    `json:"error"`
}
