package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/NotLe0n/tML-readme-card/src/widgets"
	"github.com/g4s8/hexcolor"
	"github.com/spf13/viper"
)

var wg sync.WaitGroup

var serverHandler *http.ServeMux
var server http.Server

func main() {
	setup_config()

	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":" + viper.GetString("port"), Handler: serverHandler}

	serverHandler.HandleFunc("/", generateImageHandler)

	wg.Add(1)
	go func() {
		defer wg.Done() //tell the waiter group that we are finished at the end
		cmdInterface()
		log.Println("cmd goroutine finished")
	}()

	log.Println("server starting on Port " + viper.GetString("port"))
	var err error
	if viper.GetBool("useHTTPS") {
		if viper.GetString("certPath") == "" || viper.GetString("keyPath") == "" {
			log.Fatal("'certPath' or 'keyPath' cannot be empty")
		}
		err = server.ListenAndServeTLS(viper.GetString("certPath"), viper.GetString("keyPath"))
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	} else if err == http.ErrServerClosed {
		log.Println("Server not listening anymore")
	}

	wg.Wait()
}

func setup_config() {
	viper.SetDefault("port", "8005")
	viper.SetDefault("useHTTPS", false)
	viper.SetDefault("certPath", "")
	viper.SetDefault("certKey", "")

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}

	json_string, err := json.MarshalIndent(viper.AllSettings(), "", "\t")
	if err == nil {
		log.Printf("using config: \n%s\n", json_string)
	}
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

	config := widgets.ImgConfig{
		TextColor:    textColor,
		BgColor:      bgColor,
		BorderColor:  borderColor,
		BorderWidth:  borderWidth,
		CornerRadius: cornerRadius,
		Version:      version,
		Font:         font,
	}

	if q.Has("steamid64") {
		img, err := widgets.GenerateAuthorWidget(q.Get("steamid64"), config)
		if err != nil {
			log.Println(err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		io.Copy(w, bytes.NewReader(img))
	} else if q.Has("modname") {
		img, err := widgets.GenerateModWidget(q.Get("modname"), config)
		if err != nil {
			log.Println(err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		io.Copy(w, bytes.NewReader(img))
	}

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

// write a json error to w
func errorJson(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
