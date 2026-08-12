# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY packages/ packages/
COPY services/configuration/ services/configuration/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/configuration ./services/configuration/cmd/configuration

FROM alpine:3.20
RUN apk add --no-cache ca-certificates netcat-openbsd \
	&& adduser -D -H -u 10001 appuser
WORKDIR /app
COPY --from=build /out/configuration /app/configuration
COPY docker/bin/service-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /app/configuration
USER appuser
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/configuration"]
