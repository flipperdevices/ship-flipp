package main

import (
	"errors"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func getEasyShipStatus(webToken, companyID string) (*easyShipStatus, error) {
	req := &fasthttp.Request{}
	req.SetRequestURI("https://api.easyship.com/api/v2/companies/" + companyID + "/analytics?scopes=shipments_count_by_in_progress_status")
	req.Header.Set("Authorization", "Bearer "+webToken)
	res := &fasthttp.Response{}
	err := fasthttp.Do(req, res)
	if err != nil {
		return nil, err
	}
	shipping := fastjson.GetInt(res.Body(), "in_progress_shipments_count")
	delivered := fastjson.GetInt(res.Body(), "completed_shipments_count")
	if delivered == 0 {
		return nil, errors.New("unknown error: " + string(res.Body()))
	}
	return &easyShipStatus{
		Total:     shipping + delivered,
		Delivered: delivered,
		Date:      time.Now().UTC(),
	}, nil
}

type easyShipStatus struct {
	Total     int
	Delivered int
	Date      time.Time
}
