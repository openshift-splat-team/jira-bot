# jira-bot

jira-bot is run peridocally to compile status on SPLAT epics. The bot runs hourly as a cron job on the vsphere02 cluster in the `vsphere-infra-helper` project. 

Building:
```
./hack/build.sh
```

Testing:
```
export JIRA_PERSONAL_ACCESS_TOKEN=<your token>
./bin/splat-jira-bot
```

Deploying:
When a commit is pushed to the repo a github action picks it up, builds, and pushes an updated image to `quay.io/ocp-splat/jira-bot`.