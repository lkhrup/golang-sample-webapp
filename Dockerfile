FROM golang:1.15-alpine3.12 as builder
RUN apk --no-cache add git
ARG GOPROXY
ENV GO111MODULE=on
WORKDIR /build/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.12
WORKDIR /root/
COPY --from=builder /build/app .
COPY views views
CMD ["./app"]
