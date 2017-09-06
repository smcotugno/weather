package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"text/template"
	"time"
)

const endpoint = "https://api.wunderground.com/api/e60f03e8db8ae4df/history_{{.Date}}/q/{{.State}}/{{.City}}.json"

const inLayout = "2006-01-02"
const outLayout = "20060102"

type Failure struct {
	Error string `json:"error"`
}

type Summary struct {
	MaxTemp string
	MinTemp string
	MaxPressure string
	MinPressure string
}

func formatDate(value string) (string, *Failure) {
	t, err := time.Parse(inLayout, value)
	if err != nil {
		return "", &Failure{err.Error()}
	}

	return t.Format(outLayout), nil
}

func formatCity(value string) string {
	reg, _ := regexp.Compile("\\s")

	return reg.ReplaceAllString(value, "_")
}

func callWeather(request string) (*Summary, *Failure) {
	resp, err := http.Get(request)
	if err != nil {
		return nil, &Failure{err.Error()}
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &Failure{err.Error()}
	}

	var hist WUHistory

	err = json.Unmarshal(payload, &hist)
	if err != nil {
		return nil, &Failure{err.Error()}
	}


	if len(hist.History.Dailysummary) < 1 {
		return nil, &Failure{"Unknown city"}
	}

	sum := Summary{MaxTemp: hist.History.Dailysummary[0].Maxtempi,
		MinTemp: hist.History.Dailysummary[0].Mintempi,
		MaxPressure: hist.History.Dailysummary[0].Maxpressurei,
		MinPressure: hist.History.Dailysummary[0].Minpressurei}

	return &sum, nil
}

func getWeather(params url.Values) (*Summary, *Failure) {
	var buffer bytes.Buffer

	tmpl, terr := template.New("request_params").Parse(endpoint)
	if terr != nil {
		return nil, &Failure{"Parameter parsing error."}
	}

	query := struct {
		Date string
		State string
		City string}{}

	date, err := formatDate(params.Get("date"))
	if err != nil {
		return nil, err
	}

	query.Date = date
	query.State = params.Get("state")
	query.City = formatCity(params.Get("city"))

	if query.Date == "" {
		return nil, &Failure{"Date is required."}
	} else if query.State == "" {
		return nil, &Failure{"State is required."}
	} else if query.City == "" {
		return nil, &Failure{"City is required."}
	}


	terr = tmpl.Execute(&buffer, query)

	payload, err := callWeather(buffer.String())
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
        var data []byte
	w.Header().Add("Content-Type", "application/json")

	payload, err := getWeather(r.URL.Query())
	if err == nil {
		resp, _ := json.Marshal(payload)
		data = resp
	} else {
		resp, _ := json.Marshal(err)
		data = resp
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/weather", handleWeather)

	server := http.Server{
		Addr: ":9393",
		Handler: mux,
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 30 * time.Second}

	log.Fatal(server.ListenAndServe())
}
