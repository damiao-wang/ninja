package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/health/service/ninja", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte{})
	})

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	sid := "ninja:1234"
	if err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      sid,
		Name:    "ninja",
		Port:    1234,
		Address: "127.0.0.1",
	}); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Agent().ServiceDeregister(sid); err != nil {
			log.Fatal(err)
		}
	}()

	http.ListenAndServe(":1234", r)
}
