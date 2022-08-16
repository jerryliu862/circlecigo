# Stage 1: pull official base image and build executable file

FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . .

RUN apk --no-cache add git

# RUN go env -w GOPROXY=https://goproxy.io,direct

RUN go mod download

RUN go build -o app

# Stage 2: copy and run executable file

FROM alpine

WORKDIR /app

COPY --from=builder /app ./

CMD ["./app"]
