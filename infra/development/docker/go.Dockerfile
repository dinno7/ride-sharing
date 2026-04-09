FROM golang:1.26-alpine AS builder

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

RUN echo -ne "nameserver 45.159.149.19\nnameserver 185.8.174.140" > /etc/resolv.conf;
RUN go env -w GOPROXY=https://go.devneeds.ir,direct && go env -w GOSUMDB=off;

# RUN go install github.com/cosmtrek/air@latest;

WORKDIR /data

COPY ./go.mod ./go.sum ./air ./.air.toml  ./
COPY ./services/$SERVICE_NAME ./services/$SERVICE_NAME
COPY ./shared ./shared

RUN  go mod tidy && go mod download;

EXPOSE 4000

CMD ["./air", "--build.cmd", "go build -o ./bin/main ./services/$SERVICE_NAME/cmd"]
