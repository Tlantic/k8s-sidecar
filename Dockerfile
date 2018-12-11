FROM golang:1.11.1 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/main /app/
WORKDIR /app
EXPOSE 50051
CMD ["./main"]