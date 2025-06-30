package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

// Structs para resposta da API ViaCEP
type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro,omitempty"`
}

// Structs para resposta da API WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
	Error WeatherAPIError `json:"error,omitempty"`
}

type WeatherAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Struct para resposta final
type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

var (
	port          string
	weatherAPIKey string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env, usando variáveis de ambiente do sistema.")
	}

	port = os.Getenv("PORT")
	weatherAPIKey = os.Getenv("WEATHERAPI_KEY")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", weatherHandler)
	if port == "" {
		port = "8080"
	}
	log.Printf("~Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	cepRegex := regexp.MustCompile(`^\d{8}$`)
	if !cepRegex.MatchString(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := getCityByCEP(cep)
	if err == errNotFound {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	} else if err != nil {
		println("Error getting city by CEP:", city, err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	tempC, err := getTemperatureByCity(city)
	if err != nil {
		println("Error getting temperature by CEP:", err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := TemperatureResponse{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

var errNotFound = fmt.Errorf("not found")

func getCityByCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		println("Get viacep", cep, err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var viaCEP ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEP); err != nil {
		println("Unmarshal", cep, err)
		return "", err
	}
	if viaCEP.Erro {
		println("viaCEP", cep, err)
		return "", errNotFound
	}
	return viaCEP.Localidade, nil
}

func getTemperatureByCity(city string) (float64, error) {
	if weatherAPIKey == "" {
		return 0, fmt.Errorf("WEATHERAPI_KEY not set")
	}
	// Remove acentos e substitui ç por c
	normalize := func(s string) string {
		s = strings.ToLower(s)
		s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "+")
		replacements := []struct {
			old string
			new string
		}{
			{"á", "a"}, {"à", "a"}, {"ã", "a"}, {"â", "a"}, {"ä", "a"},
			{"é", "e"}, {"è", "e"}, {"ê", "e"}, {"ë", "e"},
			{"í", "i"}, {"ì", "i"}, {"î", "i"}, {"ï", "i"},
			{"ó", "o"}, {"ò", "o"}, {"õ", "o"}, {"ô", "o"}, {"ö", "o"},
			{"ú", "u"}, {"ù", "u"}, {"û", "u"}, {"ü", "u"},
			{"ç", "c"},
		}
		for _, r := range replacements {
			s = strings.ReplaceAll(s, r.old, r.new)
		}
		return s
	}
	cityNormalized := normalize(city)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", weatherAPIKey, cityNormalized)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var weather WeatherAPIResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return 0, err
	}

	if weather.Error.Message != "" {
		return 0, fmt.Errorf("weatherapi error: %s", weather.Error.Message)
	}

	return weather.Current.TempC, nil
}
