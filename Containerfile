FROM quay.io/centos/centos:stream9-minimal
COPY bin/jira-bot .

CMD ./jira-bot

