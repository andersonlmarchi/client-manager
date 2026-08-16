# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
COPY packages/ packages/
COPY services/identity/ services/identity/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/identity ./services/identity/cmd/identity

FROM alpine:3.20
RUN apk add --no-cache ca-certificates netcat-openbsd \
	&& adduser -D -H -u 10001 appuser
WORKDIR /app
COPY --from=build /out/identity /app/identity
COPY docker/bin/service-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /app/identity
USER appuser
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/identity"]
