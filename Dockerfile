FROM golang:1.17.8-alpine AS build

WORKDIR /app

RUN apk add --no-cache gcc g++

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# RUN go build -o srv main.go && chmod +x ./srv
RUN go build -o srv 

FROM alpine
WORKDIR /app
COPY --from=build /app/srv /app
CMD ["./srv"]
