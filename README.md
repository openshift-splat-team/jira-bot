# jira-bot

jira-bot is run peridocally to compile status on SPLAT epics. The bot runs hourly as a cron job on the vsphere02 cluster in the `vsphere-infra-helper` project.

## Building
```
./hack/build.sh
```

## Running the CLI

The jira-bot now supports both **on-premise Jira** (e.g., issues.redhat.com) and **Atlassian Cloud** (e.g., yourcompany.atlassian.net).

### On-Premise Jira

```sh
export JIRA_PROJECT=SPLAT
export JIRA_BOARD="SPLAT - Scrum Board"
export JIRA_PERSONAL_ACCESS_TOKEN=<your Jira personal token>
# Optional: specify base URL (defaults to https://issues.redhat.com)
export JIRA_BASE_URL=https://issues.redhat.com

podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./jira-bot help
```

### Atlassian Cloud

For Atlassian Cloud, you need to generate an API token:
1. Visit: https://id.atlassian.com/manage-profile/security/api-tokens
2. Click "Create API token"
3. Give it a name and copy the token

```sh
export JIRA_PROJECT=SPLAT
export JIRA_BOARD="SPLAT - Scrum Board"
export JIRA_BASE_URL=https://yourcompany.atlassian.net
export JIRA_EMAIL=your.email@company.com
export JIRA_API_TOKEN=<your Atlassian Cloud API token>

podman run -e JIRA_BASE_URL=$JIRA_BASE_URL -e JIRA_EMAIL=$JIRA_EMAIL -e JIRA_API_TOKEN=$JIRA_API_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./jira-bot help
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `JIRA_PROJECT` | Yes | The Jira project key (e.g., SPLAT) |
| `JIRA_BOARD` | Yes | The name of the Scrum board |
| `JIRA_BASE_URL` | No | Base URL of your Jira instance (defaults to https://issues.redhat.com) |
| `JIRA_PERSONAL_ACCESS_TOKEN` | On-prem only | Bearer token for on-premise Jira |
| `JIRA_EMAIL` | Cloud only | Email address for Atlassian Cloud |
| `JIRA_API_TOKEN` | Cloud only | API token for Atlassian Cloud |

## Transitioning an issue to a new state

```sh
podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e  JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./jira-bot issue update-size-and-priority SPLAT-1450 --state="to do" —dry-run=false —override=true
```

## Assigning story points to an issue

```sh
podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e  JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./jira-bot issue update-size-and-priority SPLAT-1450 --points=5 —dry-run=false —override=true
```

## Assigning priority to an issue:

```sh
podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e  JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./jira-bot issue update-size-and-priority SPLAT-1450 --priority="to do" —dry-run=false —override=true
```

