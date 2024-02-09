FROM quay.io/centos/centos:stream9-minimal
COPY bin/splat-jira-bot .

CMD splat-jira-bot

