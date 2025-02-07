ARG FOLDERNAME=monitor
ARG FINALIMAGE=scratch
FROM golang:1.23.5 AS builder
WORKDIR /build
COPY ./ ./
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o ./${FOLDERNAME}/main

FROM ${FINALIMAGE}
WORKDIR /app
COPY --from=builder /build/${FOLDERNAME}/main ./main
#EXPOSE 80
ENTRYPOINT ["./main"]