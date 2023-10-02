package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/viper"
)

// Config represents the structure of the yaml config we parse
type Config struct {
	Projects []Project
}

// Project represents a Hetzner project to nuke with it's configurations
type Project struct {
	Name                  string
	EnvSecretRef          string
	secret                string
	DisabledProducts      []string
	ExcludeServerSelector string // servers (and their related resources) that match this selector will not be nuked
	client                *hcloud.Client
}

var C Config

func main() {

	// Viper config settings
	viper.SetConfigName("hcloud-nuke-config")
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
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < len(C.Projects); i++ {
		wg.Add(1)
		project := &C.Projects[i]
		go func() {
			defer wg.Done()
			project.Nuke(ctx)
		}()
	}

	// some shutdown magic
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	taskCh := make(chan int, 1)
	go func() {
		wg.Wait()
		taskCh <- 1
	}()

	select {
	case <-taskCh:
		log.Print("Done")
	case <-sigCh:
		cancel()
		log.Print("Programm has been canceled")
	}

}

func (p *Project) Nuke(ctx context.Context) {

	err := p.SetupClient(ctx)
	if err != nil {
		log.Printf("couldn't setup client in project %s: %s", p.Name, err)
		return
	}

	p.DeleteServers(ctx)
	p.DeleteVolumes(ctx)
	p.DeleteImages(ctx)
	p.DeleteLoadBalancers(ctx)
	p.DeleteFloatingIPs(ctx)
	p.DeleteFirewalls(ctx)
	p.DeleteNetworks(ctx)
	p.DeletePrimaryIPs(ctx)
	p.DeletePlacementGroups(ctx)
	p.DeleteSSHKeys(ctx)
	p.DeleteCertificates(ctx)

}

// SetupClient makes sure the secret is set and the token is valid
func (p *Project) SetupClient(ctx context.Context) error {

	if p.secret == "" {
		// Grab the secret from env
		sec, found := os.LookupEnv(p.EnvSecretRef)
		if !found {
			return fmt.Errorf("env var %s for project %s not found", p.EnvSecretRef, p.Name)
		}
		p.secret = sec
	}

	// Create new client and check
	p.client = hcloud.NewClient(hcloud.WithToken(p.secret))
	_, _, err := p.client.Datacenter.List(ctx, hcloud.DatacenterListOpts{})
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) DeleteServers(ctx context.Context) {

	servers, _, err := p.client.Server.List(ctx, hcloud.ServerListOpts{})
	if err != nil {
		log.Printf("cannot list servers for project %s: %s", p.Name, err)
		return
	}

	for _, server := range servers {
		_, _, err := p.client.Server.DeleteWithResult(ctx, server)
		if err != nil {
			log.Printf("cannot delete server %d in project %s: %s", server.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteVolumes(ctx context.Context) {
	volumes, _, err := p.client.Volume.List(ctx, hcloud.VolumeListOpts{})
	if err != nil {
		log.Printf("cannot list volumes for project %s: %s", p.Name, err)
		return
	}

	for _, volume := range volumes {
		_, err = p.client.Volume.Delete(ctx, volume)
		if err != nil {
			log.Printf("cannot delete volume %d in project %s: %s", volume.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteSSHKeys(ctx context.Context) {
	keys, _, err := p.client.SSHKey.List(ctx, hcloud.SSHKeyListOpts{})
	if err != nil {
		log.Printf("cannot list ssh keys for project %s: %s", p.Name, err)
		return
	}
	filteredKeys := []*hcloud.SSHKey{}
	for _, key := range keys {
		if key.Name != "rootkey" && key.Name != "yubikey" {
			filteredKeys = append(filteredKeys, key)
		}
	}

	for _, key := range filteredKeys {
		_, err := p.client.SSHKey.Delete(ctx, key)
		if err != nil {
			log.Printf("cannot delete ssh key %d in project %s: %s", key.ID, p.Name, err)

		}
	}

}

func (p *Project) DeletePrimaryIPs(ctx context.Context) {
	pips, _, err := p.client.PrimaryIP.List(ctx, hcloud.PrimaryIPListOpts{})
	if err != nil {
		log.Printf("cannot list primary ips for project %s: %s", p.Name, err)
		return
	}
	for _, ip := range pips {
		_, err := p.client.PrimaryIP.Delete(ctx, ip)
		if err != nil {
			log.Printf("cannot delete primary ip %d in project %s: %s", ip.ID, p.Name, err)
		}
	}
}

func (p *Project) DeletePlacementGroups(ctx context.Context) {
	pgs, _, err := p.client.PlacementGroup.List(ctx, hcloud.PlacementGroupListOpts{})
	if err != nil {
		log.Printf("cannot list placement groups for project %s: %s", p.Name, err)
		return
	}
	for _, pg := range pgs {
		_, err := p.client.PlacementGroup.Delete(ctx, pg)
		if err != nil {
			log.Printf("cannot delete placement group %d in project %s: %s", pg.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteLoadBalancers(ctx context.Context) {
	loadbalancers, _, err := p.client.LoadBalancer.List(ctx, hcloud.LoadBalancerListOpts{})
	if err != nil {
		log.Printf("cannot list load balancers for project %s: %s", p.Name, err)
		return
	}
	for _, loadbalancer := range loadbalancers {
		_, err := p.client.LoadBalancer.Delete(ctx, loadbalancer)
		if err != nil {
			log.Printf("cannot delete load balancer %d in project %s: %s", loadbalancer.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteImages(ctx context.Context) {
	images, _, err := p.client.Image.List(ctx, hcloud.ImageListOpts{})
	filteredImages := []*hcloud.Image{}

	for _, img := range images {
		if img.Type == "snapshot" || img.Type == "backup" {
			filteredImages = append(filteredImages, img)
		}
	}

	if err != nil {
		log.Printf("cannot list images for project %s: %s", p.Name, err)
		return
	}
	for _, img := range filteredImages {
		_, err := p.client.Image.Delete(ctx, img)
		if err != nil {
			log.Printf("cannot delete image %d in project %s: %s", img.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteFloatingIPs(ctx context.Context) {
	floats, _, err := p.client.FloatingIP.List(ctx, hcloud.FloatingIPListOpts{})
	if err != nil {
		log.Printf("cannot list floating ips for project %s: %s", p.Name, err)
		return
	}
	for _, float := range floats {
		_, err := p.client.FloatingIP.Delete(ctx, float)
		if err != nil {
			log.Printf("cannot delete floating ip %d in project %s: %s", float.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteFirewalls(ctx context.Context) {
	firewalls, _, err := p.client.Firewall.List(ctx, hcloud.FirewallListOpts{})
	if err != nil {
		log.Printf("cannot list firewalls for project %s: %s", p.Name, err)
		return
	}
	for _, fw := range firewalls {
		_, err := p.client.Firewall.Delete(ctx, fw)
		if err != nil {
			log.Printf("cannot delete firewall %d in project %s: %s", fw.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteCertificates(ctx context.Context) {
	certs, _, err := p.client.Certificate.List(ctx, hcloud.CertificateListOpts{})
	if err != nil {
		log.Printf("cannot list certificates for project %s: %s", p.Name, err)
		return
	}
	for _, cert := range certs {
		_, err := p.client.Certificate.Delete(ctx, cert)
		if err != nil {
			log.Printf("cannot delete certificate %d in project %s: %s", cert.ID, p.Name, err)
		}
	}
}

func (p *Project) DeleteNetworks(ctx context.Context) {
	networks, _, err := p.client.Network.List(ctx, hcloud.NetworkListOpts{})
	if err != nil {
		log.Printf("cannot list networks for project %s: %s", p.Name, err)
		return
	}
	for _, net := range networks {
		_, err := p.client.Network.Delete(ctx, net)
		if err != nil {
			log.Printf("cannot delete network %d in project %s: %s", net.ID, p.Name, err)
		}
	}
}
