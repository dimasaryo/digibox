package main

import (
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

const (
	Token = "my token"
)

type TokenSource struct {
	AccessToken string
}

type DigitalOceanClient struct {
	Client *godo.Client
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}

	return token, nil
}

func (d *DigitalOceanClient) Start() error {
	// Get snapshot
	// Create droplets based on the snapshot
	// Delete snapshot
	return nil
}

func (d *DigitalOceanClient) Stop() error {
	// Power off droplets
	// Create snapshot from current droplets
	// Delete droplets
	return nil
}
