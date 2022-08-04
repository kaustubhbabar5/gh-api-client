FROM golang:1.18-alpine as builder 

WORKDIR /build/

COPY /go.mod /go.mod

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go

FROM golang:1.18-alpine


RUN apk update && apk add ca-certificates
RUN adduser -D -g '' admin

WORKDIR /app/

COPY --from=builder /build/main ./api
USER admin

ENTRYPOINT ["/app/api"]
