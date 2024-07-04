FROM golang:1.18-buster AS module-builder

RUN apt-get update && apt-get install -y gcc-arm-linux-gnueabihf

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

ARG APP_NAME
ARG GITHUB_TOKEN

RUN go env -w GOPRIVATE=github.com/NubeIO
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/NubeIO".insteadOf "https://github.com/NubeIO"

RUN go mod download

COPY . .

RUN env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ go build -o ${APP_NAME}
