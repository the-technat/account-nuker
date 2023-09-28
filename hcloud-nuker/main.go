package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/viper"
)

const (
	OWNER = "alleaffengaffen"
	REPO  = "account-nuker"
)

type Config struct {
	Projects []Project
}

type Project struct {
	Name         string
	EnvSecretRef string
	client       *hcloud.Client
}

var C Config

func main() {

	// Viper config settings
	viper.SetConfigName("hetzner-nuke-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	// Viper config magic
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// Run through all projects
	var wg sync.WaitGroup
	for i := 1; i <= len(C.Projects); i++ {
		go func() {
			defer wg.Done()
			err = NukeProject(&C.Projects[i])
			if err != nil {
				log.Printf("Project %s failed to nuke: %s", C.Projects[i].Name, err)
			}
		}()
	}
	wg.Wait()

}

func NukeProject(p *Project) error {

	err := p.SetupClient()
	if err != nil {
		return err
	}

	// First all servers
	servers, _, err := p.client.Server.List(context.TODO(), hcloud.ServerListOpts{})
	if err != nil {
		return err
	}
	for _, server := range servers {
		_, _, err := p.client.Server.DeleteWithResult(context.TODO(), server)
	}
	if err != nil {
		return err
	}
}

func (p *Project) SetupClient() error {
	// Grab the secret from env
	sec, found := os.LookupEnv(p.EnvSecretRef)
	if !found {
		return fmt.Errorf("env var %s for project %s not found", p.EnvSecretRef, p.Name)
	}

	// Create new client and check
	p.client = hcloud.NewClient(hcloud.WithToken(sec))
	_, _, err := p.client.Datacenter.List(context.TODO(), hcloud.DatacenterListOpts{})
	if err != nil {
		return fmt.Errorf("cannot proceed with project %s: %s", p.Name, err)
	}
}
