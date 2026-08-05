package main 

import (
	"log"
	"encoding/json"
	"net/http"
	"time"
	// imports para criação das métricas
	// Fontes: https://prometheus.io/docs/guides/go-application/
	// 	   https://prometheus.io/docs/tutorials/instrumenting_http_server_in_go/
	// 	   https://dev.to/sunnynazar/the-complete-guide-to-prometheus-metric-types-promql-alerting-and-troubleshooting-5a69
	"github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

//Definindo variáveis das métricas

var (
	requestVolume = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Volume de requisições",
		},
		[]string{"path"},
	)
	serviceAvailability = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "service_available",
			Help: "Disponibilidade do serviço (1 up/0 down)",
		},	
	)
)	

func init() {
	prometheus.MustRegister(requestVolume)
	prometheus.MustRegister(serviceAvailability)
	// Inicializa a variável como disponível
	serviceAvailability.Set(1)
}

//Retorna estrutura JSON
func handler (w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json") //Arquivo JSON

	json.NewEncoder(w).Encode(map[string]string{
		"nome": "Projeto Korp",
		"horario": time.Now().UTC().Format(time.RFC3339), 
		//Para mais formatos - fonte: https://pkg.go.dev/time
	})
}
func main(){
	http.HandleFunc("/projeto-korp", handler) //Endpoint JSON
	http.Handle("/metrics", promhttp.Handler()) //Expondo métricas
	log.Println("Executando...")
	log.Fatal(http.ListenAndServe(":8080", nil)) //Servindo na porta 8080
}

//To-Do: implementar o prometheus - OK
