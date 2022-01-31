FROM golang:1.17 as build-env
WORKDIR /go/src/app
ADD . /go/src/app
RUN go build ./cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/app/main ./main
ENV GIN_MODE=release
CMD ["/app/main"]
