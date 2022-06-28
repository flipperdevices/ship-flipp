package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

//go:embed index.tmpl
var f embed.FS
var cfg config
var latestStatus *easyShipStatus
var latestPosts []tgMessage
var latestTotalOffset int
var latestDeliveredOffset int

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	var err error
	latestStatus, err = getEasyShipStatus(cfg.EasyShipWebToken, cfg.EasyShipCompanyID)
	if err != nil {
		log.Fatalln("Status", err)
	}

	if cfg.TelegramChannel != "" {
		latestPosts, latestTotalOffset, latestDeliveredOffset, err = getFeedPosts(cfg.TelegramChannel)
		if err != nil {
			log.Fatalln("Telegram", err)
		}
	}

	latestStatus.Total += latestTotalOffset
	latestStatus.Delivered += latestDeliveredOffset

	go pollData()

	r := gin.New()
	funcMap := template.FuncMap{
		"unescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "index.tmpl"))
	r.SetHTMLTemplate(templ)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"status": latestStatus,
			"feed":   latestPosts,
		})
	})

	r.GET("/rss", func(c *gin.Context) {
		feed := &feeds.Feed{
			Title: "Flipper Zero Shipping Status",
			Link:  &feeds.Link{Href: "https://ship.flipp.dev"},
		}
		for _, p := range latestPosts {
			feed.Items = append(feed.Items, &feeds.Item{
				Content: p.Message,
				Created: p.Date.Time,
				Link:    &feeds.Link{Href: "https://ship.flipp.dev"},
			})
		}
		rss, err := feed.ToRss()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(200, "application/rss+xml", []byte(rss))
	})

	log.Println("Starting server...")
	log.Fatal(r.Run(":8080"))
}

func pollData() {
	for {
		if cfg.TelegramChannel != "" {
			posts, totalOffset, deliveredOffset, err := getFeedPosts(cfg.TelegramChannel)
			if err != nil {
				log.Println("Telegram", err)
			} else {
				latestPosts = posts
				latestTotalOffset = totalOffset
				latestDeliveredOffset = deliveredOffset
			}
		}

		status, err := getEasyShipStatus(cfg.EasyShipWebToken, cfg.EasyShipCompanyID)
		if err != nil {
			log.Println("Status", err)
		} else {
			status.Total += latestTotalOffset
			status.Delivered += latestDeliveredOffset
			log.Println("Status", "Total:", status.Total, "Delivered:", status.Delivered)
			latestStatus = status
		}

		time.Sleep(time.Minute * 10)
	}
}
