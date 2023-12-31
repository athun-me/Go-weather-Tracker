package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Celsius float64 `json:"temp"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err

	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func checkWeather(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		City    string
		Kelvin  float64
		Celsius float64
	}

	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Println(err)
		}
		temp.Execute(w, nil)

	} else if r.Method == http.MethodPost {

		r.ParseForm()
		TempData := &Data{
			City: r.Form.Get("city"),
		}
		city := TempData.City

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		TempData.City = data.Name
		TempData.Kelvin = data.Main.Kelvin
		TempData.Celsius = data.Main.Kelvin - 273.15
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, TempData)
	}
}

func main() {

	http.HandleFunc("/weather", checkWeather)

	http.ListenAndServe(":8000", nil)
}
