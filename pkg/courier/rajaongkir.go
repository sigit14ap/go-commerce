package courier

import (
	"encoding/json"
	"github.com/sigit14ap/go-commerce/internal/config"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
}

type Provider struct {
	cfg *config.Config
}

func NewCourierProvider(cfg *config.Config) *Provider {
	return &Provider{
		cfg: cfg,
	}
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

func GetProvinces() ([]domain.Province, error) {
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
			ID:   prov.ProvinceID,
			Name: prov.Province,
		})
	}

	return provinceList, err
}

func GetCities(cityListDTO dto.CityListDTO) ([]domain.City, error) {
	var data interface{}
	endpoint := "/city?province=" + strconv.Itoa(cityListDTO.ProvinceID)
	cities, err := call("GET", endpoint, data)

	if err != nil {
		return []domain.City{}, err
	}

	var response cityResponse
	err = json.Unmarshal(cities, &response)

	cityList := []domain.City{}

	for _, city := range response.Rajaongkir.Results {
		cityList = append(cityList, domain.City{
			ID:           city.CityID,
			Name:         city.City,
			ProvinceID:   city.ProvinceID,
			ProvinceName: city.Province,
		})
	}

	return cityList, err
}
