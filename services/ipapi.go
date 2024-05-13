package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type IpDetails struct {
	Location struct {
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"location"`
}

func FindIpLocation(ip string) (ipd IpDetails, err error) {
	url := fmt.Sprintf("https://api.ipapi.is/?q=%s", ip)
	response, err := http.Get(url)
	if err != nil {
		return ipd, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ipd, err
	}

	if err := json.Unmarshal(data, &ipd); err != nil {
		return ipd, err
	}

	return ipd, nil
}
