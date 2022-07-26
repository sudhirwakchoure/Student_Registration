FROM golang:1.18-alpine AS builder

WORKDIR /test
COPY . .
ENV CGO_ENABLED=0 
RUN go build -o main .

FROM alpine:latest
COPY --from=builder /test/config.properties config.properties
COPY --from=builder /test/main main

EXPOSE 8080
ENTRYPOINT ["./main"]


# FROM golang:1.18-alpine 
# RUN mkdir /restapi
# WORKDIR /restapi
# COPY  . /restapi/
# RUN go mod download
# RUN go build -o main
# EXPOSE 3000
# CMD ["/restapi/main"]