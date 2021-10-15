FROM golang:1.17-alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download -modcacherw

COPY --chown=nobody:nogroup ./main.go ./main.go

RUN go build -trimpath -o readelf-sections .

FROM alpine:3.14

USER nobody

COPY --chown=0:0 --from=builder /app/readelf-sections /readelf-sections

CMD ["/readelf-sections"]
