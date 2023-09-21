# account-nuker
Nukes cloud accounts on a schedule

## aws-nuke

Based on [aws-nuke](https://github.com/rebuy-de/aws-nuke). Nukes the lab account every sunday & wednesday at 23:00 PM.

To setup nuking all we did was to create a new user named `nuker` which has the Administrator Role and a key-pair which was saved as Github action secrets

## azure-nuke

Based on [azure-nuke](https://github.com/ekristen/azure-nuke). Nukes the lab account every sunday & wednesday at 23:00 PM.

To setup nuking, all we did was pasting the following output into a github actions secret: ` az ad sp create-for-rbac --name "nuker" --role Contributor --scopes /subscriptions/559e87b7-6bd2-4c2a-a6f4-2c0e8b5b0edb --json-auth`

## To Do

I can imagine some improvements:
- [ ] Skip nuke if a message on slack/telegram or co has been received
- [ ] Push message if nuke failed
- [ ] Add Hetzner Cloud
- [ ] Add DigitalOcean Cloud
