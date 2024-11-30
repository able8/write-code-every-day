package main

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 1 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		ResponseHeaderTimeout: 15 * time.Second,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   10,
	},
}

func main() {
	req, err := http.NewRequest("GET", "http://localhost:8888/", nil)
	if err != nil {
		panic(err)
	}
	requestFn := func() {
		resp, err := client.Do(req)
		if err != nil {
			slog.ErrorContext(req.Context(), "http get error", slog.Any("error", err))
			return
		}
		defer resp.Body.Close()
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			slog.ErrorContext(req.Context(), "http read error", slog.Any("error", err))
			return
		}
		slog.InfoContext(req.Context(), "response", slog.String("status", resp.Status))
	}
	for i := 0; i < 100; i++ {
		requestFn()
		time.Sleep(time.Second * 2)
	}
}
