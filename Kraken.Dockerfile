FROM golang:1.23.5 AS builder
WORKDIR /build
COPY ./ ./
RUN go mod tidy
RUN CGO_ENABLED=0 go build ./main.go

FROM scratch
WORKDIR /app
COPY --from=builder /build .
#EXPOSE 80
ENTRYPOINT ["./main"]