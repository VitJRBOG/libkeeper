FROM golang:1.19
WORKDIR /api
COPY libkeeper-api/go.mod libkeeper-api/go.sum libkeeper-api/main.go ./
COPY libkeeper-api/cmd/ ./cmd
COPY libkeeper-api/internal/ ./internal
RUN /bin/bash -c "go mod download && go build -o ./api && find . -not -name 'api' -delete"