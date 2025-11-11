# build
FROM golang:1.22 as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api

# run
FROM gcr.io/distroless/base-debian12
COPY --from=builder /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
