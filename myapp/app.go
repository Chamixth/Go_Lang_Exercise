package main

//imports neccessary packages to handling json data, reading files, making http requests
//and manipulating strings
import (
	"encoding/json" //import packages for json encoding and decoding
	"fmt"           //import packages for formatted IO
	"io/ioutil"     //import packages for IO operations
	"net/http"      // import package for HTTP client and server functionality
	"strings"       //import package for string manipulation
)

// define the structure of the json data that used to store the API configuration
type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"` //Struct the field for storing api key
}

// define the structure of the json data that used to store the weather details
type weatherData struct {
	Name string   `json:"name"` //struct field for storing city name
	Main struct { //struct field for string weather details
		Temp     float64 `json:"temp"`     //Temperature
		Pressure float64 `json:"pressure"` //Pressure
		Humidity int     `json:"humidity"` //Humidity
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {

	//read the API configuration file
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}
	//unmarshal the Json data into the ApiConfigData struct
	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my weather channel!\n"))
}

func query(city string) (weatherData, error) {
	//Load the API configuration
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	//Construct the api for the API Request
	url := "http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city
	fmt.Println("API URL:", url) //print api url

	//Send the API get request to the API
	resp, err := http.Get(url)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return weatherData{}, err
	}
	fmt.Println("API Response:", string(body)) //print api response

	//Unmarshal the Json response int the weather data struct
	var d weatherData
	if err := json.Unmarshal(body, &d); err != nil {
		fmt.Println("JSON decoding error:", err) //print JSON decode error
		return weatherData{}, err
	}
	d.Main.Temp = d.Main.Temp - 273.15      //convert kelvin to celcius
	fmt.Println("Decoded weather data:", d) //Print decoded weather data

	return d, nil
}

// Setup the main function
func main() {

	//HTTP handler for welcome route
	http.HandleFunc("/hello", welcome)

	//HTTP handler for weather route
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2] //Extract the city name from the URL
			data, err := query(city)                      //Query weather data for the specified city
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf8.")
			json.NewEncoder(w).Encode(data) //Encode weather data as Json and send response
		})

	//start http server on port 8080
	http.ListenAndServe(":8080", nil)
}
