package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ConfigFileName has the default config file name
const ConfigFileName = "config.json"

// Provider interface
type Provider interface {
	GetWeatherData(string) (string, error)
}

// Config holds data for external connection
type Config struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// CurrentConfig is the configuration created on initialization
var CurrentConfig Config

// OpenWeatherProvider fetch weather information from openweather service
type OpenWeatherProvider struct {
	config *Config
}

// OpenWeatherResponse contains the data returned by openweather
type OpenWeatherResponse struct {
	LocationName string `json:"name"`
	Main         struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
	}
	Sys struct {
		Country string `json:"country"`
	}
}

// GetWeatherData returns data from service openweather
func (o *OpenWeatherProvider) GetWeatherData(query string) (string, error) {
	url := fmt.Sprintf(o.config.URL, query)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		return "", fmt.Errorf("Error getting information: %s", err)
	}

	var owResponse OpenWeatherResponse

	json.NewDecoder(resp.Body).Decode(&owResponse)
	return fmt.Sprintf("Location: %s - %s, temperature: %.2f, feels like %.2f\n", owResponse.LocationName, owResponse.Sys.Country, owResponse.Main.Temp, owResponse.Main.FeelsLike), nil
}

// WriteWeatherData to get the message to display to the end user.
func WriteWeatherData(w io.Writer, p Provider, query string) {
	result, err := p.GetWeatherData(query)
	if err != nil {
		result = err.Error()
	}

	w.Write([]byte(result))
}

func handler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Path[1:]
	provider := OpenWeatherProvider{config: &CurrentConfig}

	WriteWeatherData(w, &provider, query)
}

// LoadConfiguration loads the Configuration from a file
func LoadConfiguration(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&CurrentConfig)
}

func main() {
	LoadConfiguration(ConfigFileName)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
