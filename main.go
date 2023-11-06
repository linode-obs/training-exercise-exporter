package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const Namespace = "helloworld"
const WeatherStation = "KBOS" // Boston Logan Airport

type WeatherResponse struct {
	Properties struct {
		Temperature struct {
			Value float64 `json:"value"`
		} `json:"temperature"`
	} `json:"properties"`
}

type metrics struct {
	temperature *prometheus.GaugeVec
	queries     *prometheus.CounterVec
}

func getWeather() (float64, error) {
	url := fmt.Sprintf("https://api.weather.gov/stations/%s/observations/latest", WeatherStation)
	httpClient := &http.Client{Timeout: time.Second * 30}
	res, err := httpClient.Get(url)

	if err != nil {
		log.Printf("ERROR: Could not fetch %s", url)
		return 0, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)

	if err != nil {
		log.Println("ERROR: Could not unmarshal json")
		return 0, err
	}

	temperature := weather.Properties.Temperature.Value
	return temperature, nil
}

func RegisterMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		temperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "outdoor_temperature_celsius",
				Help:      "Outdoor temperature reported by the NWS.",
			},
			[]string{"station"},
		),
		queries: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "nws_query_attempts_total",
				Help:      "Number of times we've queried the NWS API.",
			},
			[]string{"station"},
		),
	}
	reg.MustRegister(m.temperature)
	reg.MustRegister(m.queries)
	return m
}

func main() {

	reg := prometheus.NewRegistry()
	m := RegisterMetrics(reg)

	go func() {
		for {
			temperature, err := getWeather()
			m.queries.With(prometheus.Labels{"station": WeatherStation}).Inc()
			if err != nil {
				log.Print(err)
				log.Println("Not updating metrics on this run.")
				time.Sleep(time.Minute * 1)
				continue
			}
			m.temperature.With(prometheus.Labels{"station": WeatherStation}).Set(temperature)
			time.Sleep(time.Minute * 5)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":23456", nil))
}
