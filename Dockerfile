FROM golang:1.27.2
LABEL authors="iyannsch"

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY api/ api/
COPY clientset/ clientset/
COPY *.go ./

RUN go build -o /garbage-collector

RUN chmod +x /garbage-collector

ENTRYPOINT ["/garbage-collector"]