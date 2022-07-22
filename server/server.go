package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
}

type Greet struct {
	Greeting string `json:"greeting"`
}

var signalChannel chan os.Signal

var greeting Greet

func (s *Server) StartServer() {
	fmt.Println("Starting Server...")
	signalChannel = make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	r := s.CreateRoutes()
	s.server = &http.Server{Addr: ":8090", Handler: r}
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal(err)
			return
		}
	}()
	fmt.Println("Served started")

	<-signalChannel
	fmt.Println("Interrupt Detected")
	s.StopServer()
}

func (s *Server) StopServer() {
	fmt.Println("Shutting down server gracefully...")
	context := context.Background()
	if err := s.server.Shutdown(context); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server Stopped!")
}

func (s *Server) CreateRoutes() *mux.Router {
	r := mux.NewRouter()

	r = r.PathPrefix("/greet").Subrouter()
	r.HandleFunc("/", s.goServer)

	r.HandleFunc("/getGreeting", s.getGreeting)
	r.HandleFunc("/setGreeting", s.setGreeting)
	return r
}

func (s *Server) goServer(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "This is greet"); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) getGreeting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/greet/getGreeting" {
		if _, err := fmt.Fprintf(w, "invalid path"); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		if _, err := fmt.Fprintf(w, "invalid method"); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := fmt.Fprintf(w, " Greeting: %v", greeting); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) setGreeting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/greet/setGreeting" {
		if _, err := fmt.Fprintf(w, "invalid path"); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		if _, err := fmt.Fprintf(w, "invalid method"); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&greeting); err != nil {
		log.Fatal(err)
	}
	if _, err := fmt.Fprintf(w, "Greeting message set to: %s", greeting.Greeting); err != nil {
		log.Fatal(err)
	}
}
