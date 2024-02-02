FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /marcy-home-server ./cmd/app/main.go

EXPOSE 3001

# Run
CMD ["/marcy-home-server"]