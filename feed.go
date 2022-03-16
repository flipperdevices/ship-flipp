package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/valyala/fasthttp"
)

func getFeedPosts(channel string) ([]tgMessage, error) {
	_, res, err := fasthttp.Get(nil, fmt.Sprintf("https://tg.i-c-a.su/json/%s?limit=100", channel))
	if err != nil {
		return nil, err
	}
	var resp tgResponse
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Messages) == 0 {
		return nil, errors.New(string(res))
	}
	var messages []tgMessage
	for _, m := range resp.Messages {
		if m.Media == nil && m.Message != "" {
			m.Date = timestamp{m.Date.UTC()}
			messages = append(messages, m)
		}
	}
	if len(messages) == 0 {
		return nil, errors.New("no text messages")
	}
	return messages, nil
}

type tgResponse struct {
	Messages []tgMessage `json:"messages"`
}

type tgMessage struct {
	Date    timestamp   `json:"date"`
	Message string      `json:"message"`
	Media   interface{} `json:"media"`
}
