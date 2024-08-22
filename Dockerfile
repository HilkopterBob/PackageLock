# syntax=docker/dockerfile:1

FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Download Go modules

ADD . /app

RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy


# Build
ARG APP_VERSION="v0.1.0+hotfixes"
RUN \
  CGO_ENABLED=0 GOOS=linux go build -ldflags "-X 'main.AppVersion=$APP_VERSION'" -o /packagelock

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["packagelock start"]
