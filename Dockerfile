FROM golang:1.23.3-bullseye AS build-stage

WORKDIR /app

COPY . ./

RUN go mod tidy

#Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./cmd/user ./cmd/main.go


# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/cmd/user /user

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/user"]