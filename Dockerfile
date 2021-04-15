FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir -p /opt/sprout_server

WORKDIR /opt/sprout_server

COPY go.mod /opt/sprout_server/
COPY go.sum /opt/sprout_server/
RUN go mod download

COPY . /opt/sprout_server
COPY ../files/sprout_server_config/config.yaml /opt/sprout_server/config.yaml
RUN go build -o sprout_server .

EXPOSE 8081

ENTRYPOINT ["bash", "/opt/sprout_server"]