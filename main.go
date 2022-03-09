package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
)

//go:embed index.tmpl
var f embed.FS
var cfg config
var latestStatus *shippingStatus

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	status, err := getEasyShipStatus(cfg.EasyShipWebToken, cfg.EasyShipCompanyID)
	if err != nil {
		log.Fatalln("Status", err)
	}
	latestStatus = status

	go pollStatus()

	r := gin.New()

	templ := template.Must(template.New("").ParseFS(f, "index.tmpl"))
	r.SetHTMLTemplate(templ)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"status": latestStatus,
		})
	})

	log.Println("Starting server...")
	log.Fatal(r.Run(":8080"))
}

func pollStatus() {
	for {
		status, err := getEasyShipStatus(cfg.EasyShipWebToken, cfg.EasyShipCompanyID)
		if err != nil {
			log.Println("Status", err)
		} else {
			latestStatus = status
		}
		time.Sleep(time.Minute * 10)
	}
}

type shippingStatus struct {
	Total     int
	Delivered int
	Date      time.Time
}
