package servertime

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type TimeApiResp struct {
	DateTime string `json:"dateTime"`
}

func GetKyivTimeNow() (time.Time, error) {
	client := http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		"https://www.timeapi.io/api/Time/current/zone?timeZone=Europe/Kiev",
		nil,
	)
	if err != nil {
		return time.Time{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}

	var timeApiResp TimeApiResp
	err = json.Unmarshal(respBody, &timeApiResp)
	if err != nil {
		return time.Time{}, err
	}

	timeKyivNow, err := time.Parse("2006-01-02T15:04:05", timeApiResp.DateTime)
	if err != nil {
		return time.Time{}, err
	}

	return timeKyivNow, nil
}
