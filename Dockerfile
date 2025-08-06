# Go API için Dockerfile
FROM golang:1.24-alpine AS builder

# Gerekli paketleri yükle
RUN apk add --no-cache git ca-certificates tzdata

# Çalışma dizinini ayarla
WORKDIR /app

# Go modüllerini kopyala
COPY go.mod go.sum ./

# Bağımlılıkları indir
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Çalışma imajı
FROM alpine:latest

# Gerekli paketleri yükle
RUN apk --no-cache add ca-certificates tzdata

# Çalışma dizinini ayarla
WORKDIR /root/

# Derlenen uygulamayı kopyala
COPY --from=builder /app/main .

# Statik dosyaları kopyala
COPY --from=builder /app/static ./static

# Port'u aç
EXPOSE 8080

# Uygulamayı çalıştır
CMD ["./main"] 