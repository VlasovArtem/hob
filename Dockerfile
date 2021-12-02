FROM golang:1.17-alpine3.13 AS build

WORKDIR /hob
COPY go.mod go.sum main.go Makefile ./
COPY src src
RUN apk --no-cache add make git gcc libc-dev curl && make build

# -----------------------------------------------------------------------------
# Build the final Docker image

FROM alpine:3.14.3

COPY --from=build /hob/execs/hob /bin/hob
COPY content content
ENV COUNTRIES_DIR=/

ENTRYPOINT [ "/bin/hob", "-v", "a" ]
