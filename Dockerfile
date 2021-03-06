FROM golang:1.18beta1 as build-env
WORKDIR /go/src/app
ADD . /go/src/app
RUN go build ./cmd/server/main.go

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=build-env /go/src/app/main ./main
ENV GIN_MODE=release
ENV ENVIRONMENT=production
CMD ["/app/main"]
