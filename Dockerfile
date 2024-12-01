FROM golang:alpine AS builder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -o /app/balance main.go

FROM scratch
ENV TZ Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/balance /app/balance
COPY --from=builder /build/balance.yaml /app/balance.yaml
CMD ["./balance"]