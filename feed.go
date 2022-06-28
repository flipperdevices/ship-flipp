package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

const totalOffsetPrefix = "uwu "
const deliveredOffsetPrefix = "owo "

func getFeedPosts(channel string) (messages []tgMessage, totalOffset int, deliveredOffset int, err error) {
	_, res, err := fasthttp.Get(nil, fmt.Sprintf("https://tg.i-c-a.su/json/%s?limit=100", channel))
	if err != nil {
		return nil, 0, 0, err
	}
	var resp tgResponse
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(resp.Messages) == 0 {
		return nil, 0, 0, errors.New(string(res))
	}
	for _, m := range resp.Messages {
		if m.Media == nil && m.Message != "" {
			if totalOffset == 0 && strings.HasPrefix(m.Message, totalOffsetPrefix) {
				totalOffset, _ = strconv.Atoi(strings.TrimPrefix(m.Message, totalOffsetPrefix))
				continue
			}
			if deliveredOffset == 0 && strings.HasPrefix(m.Message, deliveredOffsetPrefix) {
				deliveredOffset, _ = strconv.Atoi(strings.TrimPrefix(m.Message, deliveredOffsetPrefix))
				continue
			}
			m.Date = timestamp{m.Date.UTC()}
			messages = append(messages, m)
		}
	}
	if len(messages) == 0 {
		return nil, 0, 0, errors.New("no text messages")
	}
	return
}

type tgResponse struct {
	Messages []tgMessage `json:"messages"`
}

type tgMessage struct {
	Date    timestamp   `json:"date"`
	Message string      `json:"message"`
	Media   interface{} `json:"media"`
}
