FROM golang:1.19-alpine as builder
WORKDIR /app/
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/check ./check
RUN CGO_ENABLED=0 go build -o ./bin/out ./out
RUN CGO_ENABLED=0 go build -o ./bin/in ./in

FROM alpine
RUN apk update
RUN apk add openssh
RUN mkdir -p /root/.ssh/
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
COPY --from=builder /app/bin/* /opt/resource/
