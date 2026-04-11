FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o kitchen ./cmd/kitchen/

FROM alpine:latest

COPY --from=build /app/kitchen /kitchen

EXPOSE 8081

CMD ["/kitchen"]