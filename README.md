# jira-bot

jira-bot is run peridocally to compile status on SPLAT epics. The bot runs hourly as a cron job on the vsphere02 cluster in the `vsphere-infra-helper` project. 

## Building
```
./hack/build.sh
```

## Running the cli

```sh
export JIRA_PROJECT=SPLAT
export JIRA_BOARD="SPLAT - Scrum Board"
export JIRA_PERSONAL_ACCESS_TOKEN=<your Jira personal token>
podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e  JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./splat-jira-bot help
```

## Transitioning an issue to a new state

```sh
podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e  JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./splat-jira-bot issue update-size-and-priority SPLAT-1450 --state="to do" —dry-run=false —override=true
```

## Assigning story points to an issue

```sh
podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e  JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./splat-jira-bot issue update-size-and-priority SPLAT-1450 --points=5 —dry-run=false —override=true
```

## Assigning priority to an issue:

```sh
podman run -e JIRA_PERSONAL_ACCESS_TOKEN=$JIRA_PERSONAL_ACCESS_TOKEN -e JIRA_PROJECT=$JIRA_PROJECT -e  JIRA_BOARD=$JIRA_BOARD quay.io/ocp-splat/jira-bot:latest ./splat-jira-bot issue update-size-and-priority SPLAT-1450 --priority="to do" —dry-run=false —override=true
```

