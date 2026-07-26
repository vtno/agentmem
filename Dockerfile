# syntax=docker/dockerfile:1

# Production image for the agentmem docs site.
# Build:  docker build -t agentmem-site .
# Run:    docker run --rm -p 8080:80 agentmem-site

FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
  -ldflags "-s -w -X main.version=${VERSION}" \
  -o /out/agentmem-site .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/agentmem-site /agentmem-site

EXPOSE 80

USER nonroot:nonroot

# Kamal proxy expects the app on port 80 by default.
CMD ["/agentmem-site", "-addr", ":80"]
