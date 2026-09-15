FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/PlanetsRadiusResearch \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -u 10001 app
COPY --from=build /out/server /out/migrate ./
COPY config ./config
COPY resources ./resources
USER app
EXPOSE 8080
ENTRYPOINT ["/app/server"]
