# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23.1 AS build-stage

WORKDIR /app

COPY *.mod *.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /mounting_analyzer

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /mounting_analyzer /mounting_analyzer

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/mounting_analyzer"]