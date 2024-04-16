FROM golang:1.21-bookworm AS build-stage

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY . .
RUN go mod download

RUN go build -o /nearbyassist


FROM gcr.io/distroless/base-debian12 as release-stage

WORKDIR /

COPY --from=build-stage /nearbyassist /nearbyassist

# This container exposes port 3000 to the outside world
EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT [ "/nearbyassist" ]
