package courier

import (
	"encoding/json"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type query map[string]interface{}

type status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Province struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

type provinceResponse struct {
	Rajaongkir struct {
		Query   []query    `json:"query"`
		Status  status     `json:"status"`
		Results []Province `json:"results"`
	} `json:"rajaongkir"`
}

type City struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	CityID     string `json:"city_id"`
	City       string `json:"city_name"`
}

type cityResponse struct {
	Rajaongkir struct {
		Query   query  `json:"query"`
		Status  status `json:"status"`
		Results []City `json:"results"`
	} `json:"rajaongkir"`
}

type CourierProvider interface {
	GetProvinces() ([]domain.Province, error)
	GetCities() ([]dto.ThirdPartyCityDTO, error)
	GetDeliveryCost() ([]dto.ThirdPartyCityDTO, error)
}

type Provider struct {
}

func NewCourierProvider() *Provider {
	return &Provider{}
}

func call(method string, endpoint string, data interface{}) ([]byte, error) {

	endpoint = os.Getenv("RAJAONGKIR_URL") + endpoint
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, endpoint, nil)

	var result []byte

	if err != nil {
		return result, err
	}

	req.Header.Set("key", os.Getenv("RAJAONGKIR_API_KEY"))

	response, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	return body, err
}

func (p *Provider) GetProvinces() ([]domain.Province, error) {
	var data interface{}
	province, err := call("GET", "/province", data)

	if err != nil {
		return []domain.Province{}, err
	}

	var response provinceResponse
	err = json.Unmarshal(province, &response)

	var provinceList []domain.Province

	for _, prov := range response.Rajaongkir.Results {
		provinceList = append(provinceList, domain.Province{
			ThirdPartyID: prov.ProvinceID,
			Name:         prov.Province,
		})
	}

	return provinceList, err
}

func (p *Provider) GetCities() ([]dto.ThirdPartyCityDTO, error) {
	var data interface{}
	cities, err := call("GET", "/city", data)

	if err != nil {
		return []dto.ThirdPartyCityDTO{}, err
	}

	var response cityResponse
	err = json.Unmarshal(cities, &response)

	cityList := []dto.ThirdPartyCityDTO{}

	for _, city := range response.Rajaongkir.Results {
		cityList = append(cityList, dto.ThirdPartyCityDTO{
			Name:       city.City,
			ProvinceID: city.ProvinceID,
			CityID:     city.CityID,
		})
	}

	return cityList, err
}

func (p *Provider) GetDeliveryCost() ([]dto.ThirdPartyCityDTO, error) {
	panic("implement me")
}
