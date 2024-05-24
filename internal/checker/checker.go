package checker

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/rs/zerolog/log"
)

const (
	requestTypeHTTP = "HTTP"
	requestTypePing = "PING"
)

const maxTimeout = time.Second * 15
const maxPingRequests = 3

func Online() bool {
	requestType := os.Getenv("REQUEST_TYPE")
	endpoint := os.Getenv("ENDPOINT")

	switch strings.ToUpper(requestType) {
	case requestTypeHTTP:
		return httpRequest(endpoint)
	case requestTypePing:
		return pingRequest(endpoint)
	}

	log.Error().Msg("invalid or not set request_type")
	return false
}

func httpRequest(endpoint string) bool {
	client := http.Client{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), maxTimeout)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Error().Err(err).Send()
		cancelCtx()
		return false
	}
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	if err != nil {
		log.Error().Err(err).Send()
		cancelCtx()
		return false
	}

	cancelCtx()
	return true
}

func pingRequest(endpoint string) bool {
	pinger, err := ping.NewPinger(endpoint)
	if err != nil {
		log.Error().Err(err).Send()
		return false
	}

	pinger.SetPrivileged(true)
	pinger.Count = maxPingRequests
	pinger.Timeout = maxTimeout

	err = pinger.Run()
	if err != nil {
		log.Error().Err(err).Send()
		return false
	}

	stat := pinger.Statistics()
	return stat.PacketsRecv > 0
}
