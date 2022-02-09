FROM golang:alpine as builder
RUN mkdir /build
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /builded_app cmd/main/main.go
FROM scratch
COPY --from=builder builded_app /bin/builded_app
ENTRYPOINT ["/bin/cloudflare-ddns"]