FROM golang:1.23-bookworm AS build-stage

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY . .
RUN go mod download

RUN go build -o /nearbyassist cmd/main.go

FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY views/ views/
COPY public/ public/
COPY static/ static/
COPY --from=build-stage /nearbyassist /nearbyassist

# This container exposes port 3000 to the outside world
EXPOSE 3000

CMD [ "/nearbyassist" ]
