FROM golang:1.24.0 AS builder
WORKDIR /build
COPY ./ ./
RUN go mod tidy
RUN CGO_ENABLED=0 go build ./main.go

FROM scratch
WORKDIR /app
COPY --from=builder /build .
#EXPOSE 80
ENTRYPOINT ["./main"]