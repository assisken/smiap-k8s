FROM golang:1.17-alpine

WORKDIR /usr/src/smiap-k8s

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin ./...

CMD ["reconcile"]
