# Compile stage
# FROM golang:1.10.1-alpine3.7 AS build-env
FROM golang:1.10.1-alpine3.7

# Cgo enables the creation of Go packages that call C code, set to 0 to disable it
ENV CGO_ENABLED 0

# Compile Delve
RUN apk add --no-cache git
RUN go get github.com/derekparker/delve/cmd/dlv

ADD . /go/src/remote-docker-hello

# The -gcflags "all=-N -l" flag helps us get a better debug experience
# RUN go build -gcflags "all=-N -l" -o /remote-docker-hello remote-docker-hello

# Final stage
# FROM alpine:3.7
 
# Port 8080 belongs to our application, 2345 belongs to Delve
EXPOSE 8080 2345

# Allow delve to run on Alpine based containers.
# RUN apk add --no-cache libc6-compat

WORKDIR /go/src/remote-docker-hello

# COPY --from=build-env /remote-docker-hello /
# COPY --from=build-env /go/bin/dlv /

# Run delve
CMD ["/go/bin/dlv", "debug", "--listen=0.0.0.0:2345", "--headless=true", "--api-version=2", "--log"]