# Hcloud-Nuker

Inspired by other nuker programms out there, this is a dead-simple programm that nukes one or more Hcloud projects based on a YAML config.

## Usage

```
git clone https://github.com/alleaffengaffen/account-nuker.git
cd account-nuker/hcloud-nuker
go version # make sure it's installed
go run main.go
```

## Configuration:

An example config named `hcloud-nuke.yaml` with all options populated (must be in the same directoy as the programm): 

```yaml
---
projects:
- name: demo
  envSecretRef: demo_HCLOUD_TOKEN
  # secret: "hcloud_api_token" takes predecence over the envSecretRef if specified
```