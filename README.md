# account-nuker
Nukes cloud accounts on a schedule

## aws-nuke

Based on [aws-nuke](https://github.com/rebuy-de/aws-nuke). Nukes the lab account every sunday & wednesday at 23:00 PM.

To setup nuking all we did was to create a new user named `nuker` which has the Administrator Role and a key-pair which was saved as Github action secrets

## azure-nuke

Based on [azure-nuke](https://github.com/ekristen/azure-nuke). Nukes the lab account every sunday & wednesday at 23:00 PM.

Add `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID` and `AZURE_SUBSCRIPTION_ID` from this as secrets : ` az ad sp create-for-rbac --name "nuker" --role Contributor --scopes /subscriptions/559e87b7-6bd2-4c2a-a6f4-2c0e8b5b0edb`

## Hetzner Nuking

Since I couldn' find a nuker for Hetzner Cloud, the [hcloud-nuker](./hcloud-nuker/) folder contains my own implementation using [hcloud-go](https://github.com/hetznercloud/hcloud-go). If you look at the configuration file [hetzner-nuke-config.yaml] you will likely see how it works and how you can add a new project to the list.