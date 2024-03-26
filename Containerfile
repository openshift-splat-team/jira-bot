FROM golang:1.21

WORKDIR /usr/src/app
COPY . .
RUN ./hack/build.sh
RUN mv ./bin/jira-bot ./jira-bot
CMD ["./jira-bot"]
