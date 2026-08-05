package main 

import (
	"log"
	"encoding/json"
	"net/http"
	"time"
)
func handler (w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json") //Arquivo JSON

	json.NewEncoder(w).Encode(map[string]string{
		"nome": "Projeto Korp",
		"horario": time.Now().UTC().Format(time.RFC3339),
	})
}
func main(){
	http.HandleFunc("/projeto-korp", handler) //Endpoint
	
	log.Println("Executando...")
	log.Fatal(http.ListenAndServe(":8080", nil)) //Servindo na porta 8080
}

//To-Do: implementar o prometheus
