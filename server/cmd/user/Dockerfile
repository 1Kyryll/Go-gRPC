FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o user ./cmd/user/

FROM alpine:latest

COPY --from=build /app/user /user

EXPOSE 8083
EXPOSE 9001

CMD ["/user"]
