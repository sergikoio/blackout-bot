package servertime

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type TimeApiResp struct {
	DateTime time.Time `json:"datetime"`
}

func GetKyivTimeNow() (time.Time, error) {
	client := http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		"https://www.worldtimeapi.org/api/timezone/Europe/Kiev",
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

	return timeApiResp.DateTime, nil
}
