package main

type config struct {
	EasyShipWebToken  string `env:"EASYSHIP_WEB_TOKEN,required"`
	EasyShipCompanyID string `env:"EASYSHIP_COMPANY_ID,required"`
	TelegramChannel   string `env:"TG_CHANNEL"`
}
