FROM golang:1.24.4-alpine AS builder


WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o ./sso

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /build/sso ./sso

CMD ["/app/sso"]