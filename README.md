# account-nuker

Nukes cloud accounts on a schedule

## aws-nuke

Based on [aws-nuke](https://github.com/rebuy-de/aws-nuke). Nukes the lab account every sunday & wednesday at 23:00 PM.

Auth is using an IAM role and the Github Actions OIDC provider.

## azure-nuke

Based on [azure-nuke](https://github.com/ekristen/azure-nuke). Nukes the lab account every sunday & wednesday at 23:00 PM.

Auth is using a simple service principal with Role assignments to the subscription and the Graph API (ReadWriteAll to User,Group,Application).

