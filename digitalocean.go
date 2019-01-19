package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
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

func NewDigitalOceanClient(tokenSource *TokenSource) *DigitalOceanClient {
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	do := &DigitalOceanClient{
		Client: godo.NewClient(oauthClient),
	}
	return do
}

func (d *DigitalOceanClient) Start(name string) error {
	// Get snapshot
	log.Printf("Get snapshot of %s", name)
	images, _, err := d.Client.Images.ListUser(context.Background(), nil)
	if err != nil {
		return err
	}

	var devboxImage *godo.Image
	for _, image := range images {
		if strings.Contains(image.Name, name) {
			devboxImage = &image
		}
	}

	if devboxImage == nil {
		err = errors.New(fmt.Sprintf("Could not found snapshot from %s", name))
		return err
	}

	log.Printf("Found snapshot %s", devboxImage.Name)

	// Create droplets based on the snapshot
	log.Printf("Restoring snapshot %s", devboxImage.Name)
	dropletOptions := &godo.DropletCreateRequest{
		Name:   name,
		Region: devboxImage.Regions[0],
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			ID:   devboxImage.ID,
			Slug: devboxImage.Slug,
		},
		SSHKeys:           nil,
		Backups:           false,
		IPv6:              true,
		PrivateNetworking: false,
		Monitoring:        true,
		UserData:          "",
		Volumes:           nil,
		Tags:              nil,
	}
	droplet, _, err := d.Client.Droplets.Create(context.Background(), dropletOptions)
	if err != nil {
		return err
	}
	log.Printf("Droplet has been restored %s", droplet.Name)

	// Delete snapshot
	log.Printf("Delete snapshot %s", devboxImage.Name)
	_, err = d.Client.Snapshots.Delete(context.Background(), strconv.Itoa(devboxImage.ID))
	if err != nil {
		return err
	}
	log.Printf("Snapshot %s has been deleted", devboxImage.Name)
	return nil
}

func (d *DigitalOceanClient) Stop() error {
	// Power off droplets
	// Create snapshot from current droplets
	// Delete droplets
	return nil
}
