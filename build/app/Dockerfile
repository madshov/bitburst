FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/bitburst ./cmd/app

ENTRYPOINT ["/go/bin/bitburst"]