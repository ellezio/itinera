# build stage
FROM golang:1.25.7-alpine AS builder
WORKDIR /usr/src/app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download
COPY . .

ENV CGO_ENABLED=1
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go build -v -o /usr/local/bin/app ./cmd/itinera

# runtime stage
FROM alpine:3.22
RUN apk add --no-cache sqlite-libs
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
CMD ["app"]
