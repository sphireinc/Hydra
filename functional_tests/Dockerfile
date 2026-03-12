FROM golang:1.26-bookworm

ENV CGO_ENABLED=1 \
    GO111MODULE=on \
    DEBIAN_FRONTEND=noninteractive

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        bash \
        build-essential \
        ca-certificates \
        curl \
        git \
        pkg-config \
        sqlite3 \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x /app/functional_tests/wait-and-test.sh

CMD ["/bin/bash", "/app/functional_tests/wait-and-test.sh"]