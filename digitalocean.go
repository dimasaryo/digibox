package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// TokenSource Token source
type TokenSource struct {
	AccessToken string
}

// DigitalOceanClient Digitalocean client
type DigitalOceanClient struct {
	Client *godo.Client
}

// Token Get token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}

	return token, nil
}

// NewDigitalOceanClient Create new client
func NewDigitalOceanClient(tokenSource *TokenSource) *DigitalOceanClient {
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	do := &DigitalOceanClient{
		Client: godo.NewClient(oauthClient),
	}
	return do
}

// Start Start remote development server
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

// Stop Stop remote development server
func (d *DigitalOceanClient) Stop(name string) error {
	// Get droplet
	log.Printf("Get droplet %s", name)
	droplets, _, err := d.Client.Droplets.List(context.Background(), nil)
	if err != nil {
		return err
	}

	var devbox *godo.Droplet
	for _, droplet := range droplets {
		if droplet.Name == name {
			devbox = &droplet
		}
	}

	if devbox == nil {
		err = fmt.Errorf("Could not find any droplet for %s", name)
		return err
	}

	// Power off droplets
	log.Printf("Power off droplet %s", name)
	action, _, err := d.Client.DropletActions.PowerOff(context.Background(), devbox.ID)
	if err != nil {
		return err
	}

	err = d.waitForAction(devbox.ID, action.ID)
	if err != nil {
		return err
	}

	// Create snapshot from current droplets
	log.Printf("Create snapshot from droplet %s", name)
	action, _, err = d.Client.DropletActions.Snapshot(context.Background(), devbox.ID, name)
	if err != nil {
		return err
	}

	err = d.waitForAction(devbox.ID, action.ID)
	if err != nil {
		return err
	}

	// Delete droplets
	log.Printf("Delete droplet %s", name)
	_, err = d.Client.Droplets.Delete(context.Background(), devbox.ID)
	if err != nil {
		return err
	}
	log.Printf("Droplet %s was deleted", name)
	return nil
}

func (d *DigitalOceanClient) waitForAction(devboxID, actionID int) error {
	waitTimeout := time.NewTimer(time.Duration(10) * time.Minute)
	defer waitTimeout.Stop()
	tick := time.Tick(15 * time.Second)
	for {
		select {
		case <-waitTimeout.C:
			return fmt.Errorf("Timeout after 10 minutes")
		case <-tick:
			action, _, _ := d.Client.DropletActions.Get(context.Background(), devboxID, actionID)

			log.Printf("Status for action %s is %s", action.Type, action.Status)
			if action.Status == "errored" {
				return fmt.Errorf("Error occured when doing action %s", action.Type)
			}

			if action.Status == "completed" {
				log.Printf("Action %s is completed at %s", action.Type, action.CompletedAt.String())
				return nil
			}
		}
	}
}
