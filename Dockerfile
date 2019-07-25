###############################################################################
# BUILD STAGE

FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN apk --no-cache add git
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

###############################################################################
# PACKAGE STAGE

FROM scratch
COPY --from=builder /build/app /app/
WORKDIR /app
CMD ["./app"]
