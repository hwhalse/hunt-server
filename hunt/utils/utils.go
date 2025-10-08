package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

// SendApiRequest sends an http request either inside our cluster or to an external api
func SendApiRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Println("close body error ", err)
		}
	}()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

type Response struct {
	Data    any    `json:"data"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func SendApiFailure(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	log.Info().Str("message", msg).Msg("api failure")
	res := Response{
		Status:  "Failure",
		Message: msg,
	}
	err := json.NewEncoder(w).Encode(&res)
	if err != nil {
		log.Error().Err(err).Msg("unable to send failure response")
	}
}

func SendApiSuccess(w http.ResponseWriter, data any, msg string) {
	w.Header().Set("Content-Type", "application/json")
	res := Response{
		Status:  "Success",
		Message: msg,
		Data:    data,
	}
	err := json.NewEncoder(w).Encode(&res)
	if err != nil {
		log.Error().Err(err).Msg("unable to send success response")
	}
}
