FROM golang:1.19-alpine as builder
WORKDIR /app/
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/check ./check
RUN CGO_ENABLED=0 go build -o ./bin/out ./out
RUN CGO_ENABLED=0 go build -o ./bin/in ./in

FROM alpine
COPY --from=builder /app/bin/* /opt/resource/