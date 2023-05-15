FROM alpine:3.18.0
COPY bin/* /usr/local/bin/

ENV VERSION=0.0.1
ENV GOARCH=amd64
ENV GOOS=linux
CMD ["ls", "-la", "/usr/local/bin"]
#CMD jsnog-bot-${VERSION}_${GOOS}_${GOARCH} -token ${TOKEN} -guild ${GUILD_ID}
