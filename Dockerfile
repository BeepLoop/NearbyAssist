FROM golang:1.23-bookworm AS build-stage

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY . .
RUN go mod download

RUN go build -o /nearbyassist cmd/main.go

# FROM gcr.io/distroless/base-debian12 AS release-stage
FROM debian:bookworm-slim AS release-stage

RUN apt-get update && apt-get install -y \
    chromium \
    fonts-liberation \
    libnss3 \
    libxss1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libgtk-3-0 \
    libgbm1 \
    libxshmfence1 \
    ca-certificates \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /

COPY views/ views/
COPY public/ public/
COPY static/ static/
COPY --from=build-stage /nearbyassist /nearbyassist

# Set chrome binary path for chromedp
ENV CHROME_BIN=/usr/bin/chromium
ENV PATH="$PATH:/usr/bin"

# This container exposes port 3000 to the outside world
EXPOSE 3000

CMD [ "/nearbyassist" ]
