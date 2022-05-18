package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	m "github.com/alfianabdi/mutate-ingress-for-imported-service/pkg/mutate"
	"github.com/gorilla/mux"
)

func sendError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mutating ingress for services imported using multi cluster service discovery")
}

func mutateVirtualServerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		sendError(err, w)
	}

	mutated, err := m.MutateVirtualServer(body)
	if err != nil {
		sendError(err, w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func main() {
	wait := time.Second * 10
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/virtual-server", mutateVirtualServerHandler)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        r,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
	}

	sslCrtPath := "/app/ssl/tls.crt"
	sslKeyPath := "/app/ssl/tls.key"

	go func() {
		if err := s.ListenAndServeTLS(sslCrtPath, sslKeyPath); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	s.Shutdown(ctx)
	log.Println("Shutting down")
	os.Exit(0)

}
