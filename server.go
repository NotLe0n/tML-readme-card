package main

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/g4s8/hexcolor"
)

// write an json error to w
func errorJson(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

var wg sync.WaitGroup

var serverHandler *http.ServeMux
var server http.Server

type ImgConfig struct {
	textColor    color.Color
	bgColor      color.Color
	borderColor  color.Color
	borderWidth  uint64
	cornerRadius uint64
	version      string
	font         string
}

func generateImageHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	log.Println("Got a request: " + r.URL.RawQuery)

	if r.Method != http.MethodGet {
		errorJson(w, "Method must be GET", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=86400")

	textColor, err := hexcolor.Parse("#" + q.Get("text_color"))
	if err != nil {
		textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
	}

	bgColor, err := hexcolor.Parse("#" + q.Get("bg_color"))
	if err != nil {
		bgColor = color.RGBA{R: 25, G: 28, B: 30, A: 255} // dark gray
	}

	borderColor, err := hexcolor.Parse("#" + q.Get("border_color"))
	if err != nil {
		borderColor = color.RGBA{R: 35, G: 39, B: 42, A: 255} // light gray
	}

	borderWidth, err := strconv.ParseUint(q.Get("border_width"), 10, 32)
	if err != nil {
		borderWidth = 4
	}

	cornerRadius, err := strconv.ParseUint(q.Get("corner_radius"), 10, 32)
	if err != nil {
		cornerRadius = 15
	}

	font := q.Get("font")
	if font == "" || (font != "Andy" && font != "Sans") {
		font = "Andy"
	}

	version := q.Get("v")
	if version == "" || (version != "1.4" && version != "1.3") {
		version = "1.4"
	}

	config := ImgConfig{
		textColor:    textColor,
		bgColor:      bgColor,
		borderColor:  borderColor,
		borderWidth:  borderWidth,
		cornerRadius: cornerRadius,
		version:      version,
		font:         font,
	}

	img, err := generateImage(q.Get("steamid64"), config)
	if err != nil {
		log.Println(err.Error())
		errorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.Copy(w, bytes.NewReader(img))
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
