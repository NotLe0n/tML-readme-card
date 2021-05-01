package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

//write an json error to w
func errorJson(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(`{"error":"` + msg + `"}`)
}

var wg sync.WaitGroup

var serverHandler *http.ServeMux
var server http.Server

func generateImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got a request on /")
	if r.Method != http.MethodGet {
		errorJson(w, "Method must be GET", http.StatusBadRequest)
		return
	}
	q := r.URL.Query()
	w.Header().Set("Content-Type", "image/png")
	err := generateImage(q.Get("steamid64"), w)
	if err != nil {
		log.Println(err.Error())
		errorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":3000", Handler: serverHandler}

	serverHandler.HandleFunc("/", generateImageHandler)

	wg.Add(1)
	go func() {
		defer wg.Done() //tell the waiter group that we are finished at the end
		cmdInterface()
		log.Println("cmd goroutine finished")
	}()

	log.Println("server starting on Port 3000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	} else if err == http.ErrServerClosed {
		log.Println("Server not listening anymore")
	}

	wg.Wait()
}

func cmdInterface() {
	for loop := true; loop; {
		var inp string
		if _, err := fmt.Scanln(&inp); err != nil {
			log.Println(err.Error())
		} else {
			switch inp {
			case "quit":
				log.Println("Attempting to shutdown server")
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatal("Error while trying to shutdown server: " + err.Error())
				}
				log.Println("Server was shutdown")
				loop = false
			default:
				fmt.Println("cmd not supported")
			}
		}
	}
}
