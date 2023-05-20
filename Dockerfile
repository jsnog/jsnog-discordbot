FROM debian:11.7-slim

COPY bin/* /usr/local/bin/

RUN apt update && apt install -y ca-certificates

ENV VERSION=0.0.1
ENV GOARCH=amd64
ENV GOOS=linux

CMD jsnog-bot-${VERSION}_${GOOS}_${GOARCH} -token ${TOKEN} -guild ${GUILD_ID}
