# 使用 Go 1.18 官方映像檔作為建置環境
FROM golang:1.18-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 檔案
COPY go.mod go.sum ./

# 下載依賴套件
RUN go mod download

# 複製原始碼
COPY . .

# 建置應用程式
RUN go build -o task-api .

# 使用輕量級的 alpine 映像檔作為執行環境
FROM alpine:latest

# 安裝 ca-certificates 以支援 HTTPS
RUN apk --no-cache add ca-certificates

# 設定工作目錄
WORKDIR /root/

# 從 builder 階段複製編譯好的二進位檔
COPY --from=builder /app/task-api .

# 暴露 8080 端口
EXPOSE 8080

# 執行應用程式
CMD ["./task-api"]