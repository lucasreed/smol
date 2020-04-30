FROM golang:1.14 AS build
LABEL maintainer="Luke Reed <luke@lreed.net>"
WORKDIR /go/src/github.com/lucasreed/smol
ADD . /go/src/github.com/lucasreed/smol
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" ./cmd/smolserv/

FROM gcr.io/distroless/base
COPY --from=build /go/src/github.com/lucasreed/smol/smolserv /smolserv
EXPOSE 8080
CMD ["./smolserv"]
