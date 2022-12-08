package checker

import (
	"context"
	"net/http"
	"os"
	"time"
)

const maxTimeout = time.Second * 15

func Online() bool {
	endpoint := os.Getenv("ENDPOINT")

	client := http.Client{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), maxTimeout)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		cancelCtx()
		return false
	}
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	if err != nil {
		cancelCtx()
		return false
	}

	cancelCtx()

	return true
}
