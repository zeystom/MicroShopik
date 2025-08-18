package domain

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type KeepAliveService struct {
	client    *http.Client
	serverURL string
	interval  time.Duration
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewKeepAliveService(serverURL string, interval time.Duration) *KeepAliveService {
	ctx, cancel := context.WithCancel(context.Background())

	return &KeepAliveService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		serverURL: serverURL,
		interval:  interval,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (k *KeepAliveService) Start() {
	log.Printf("Starting keep-alive service for %s with interval %v", k.serverURL, k.interval)

	go func() {
		ticker := time.NewTicker(k.interval)
		defer ticker.Stop()

		k.ping()

		for {
			select {
			case <-ticker.C:
				k.ping()
			case <-k.ctx.Done():
				log.Println("Keep-alive service stopped")
				return
			}
		}
	}()
}

func (k *KeepAliveService) Stop() {
	log.Println("Stopping keep-alive service...")
	k.cancel()
}

func (k *KeepAliveService) ping() {
	url := fmt.Sprintf("%s/health", k.serverURL)

	resp, err := k.client.Get(url)
	if err != nil {
		log.Printf("Keep-alive ping failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Keep-alive ping successful: %s", url)
	} else {
		log.Printf("Keep-alive ping returned status: %d", resp.StatusCode)
	}
}
