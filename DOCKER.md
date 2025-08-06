# 🐳 Docker Kullanım Talimatları

Bu dosya, P2P API projesini Docker ile nasıl çalıştıracağınızı açıklar.

## 🚀 Hızlı Başlangıç

### 1. Docker Compose ile Çalıştırma

```bash
# Tüm servisleri başlat
docker-compose up -d

# Logları görüntüle
docker-compose logs -f

# Servisleri durdur
docker-compose down
```

### 2. Tek Container ile Çalıştırma

```bash
# İmajı build et
docker build -t p2p-api .

# Container'ı çalıştır
docker run -d --name p2p-api-1 -p 8080:8080 p2p-api

# İkinci container'ı çalıştır
docker run -d --name p2p-api-2 -p 8081:8080 p2p-api
```

## 📋 Docker Compose Servisleri

### Üretim Ortamı (`docker-compose.yml`)

- **p2p-api-1**: Port 8080'de çalışan ilk node
- **p2p-api-2**: Port 8081'de çalışan ikinci node  
- **p2p-api-3**: Port 8082'de çalışan üçüncü node (opsiyonel)

### Geliştirme Ortamı (`docker-compose.dev.yml`)

- **p2p-api-dev**: Hot reload ile geliştirme ortamı

## 🔧 Kullanım Örnekleri

### Çoklu Node Test

```bash
# Tüm servisleri başlat
docker-compose up -d

# Web arayüzlerine erişim
# http://localhost:8080 (API 1)
# http://localhost:8081 (API 2)
# http://localhost:8082 (API 3)

# Logları kontrol et
docker-compose logs p2p-api-1
docker-compose logs p2p-api-2
```

### Geliştirme Ortamı

```bash
# Geliştirme ortamını başlat
docker-compose -f docker-compose.dev.yml up -d

# Kod değişikliklerini otomatik yükle
# (Volume mount sayesinde)
```

### Manuel Test

```bash
# İlk container'ı başlat
docker run -d --name api1 -p 8080:8080 p2p-api

# İkinci container'ı başlat
docker run -d --name api2 -p 8081:8080 p2p-api

# Container'ları kontrol et
docker ps

# Logları görüntüle
docker logs api1
docker logs api2
```

## 🌐 API Test

### cURL ile Test

```bash
# Durum kontrolü
curl http://localhost:8080/api/status

# Peer bağlantısı (API 2'den API 1'e)
curl -X POST http://localhost:8081/api/connect \
  -H "Content-Type: application/json" \
  -d '{"multiaddr": "/ip4/172.17.0.2/tcp/8080/p2p/PEER_ID"}'

# Mesaj gönder
curl -X POST http://localhost:8081/api/send \
  -H "Content-Type: application/json" \
  -d '{"content": "Docker'dan merhaba!"}'
```

### Web Arayüzü

- **API 1**: http://localhost:8080
- **API 2**: http://localhost:8081
- **API 3**: http://localhost:8082

## 🔍 Debug ve Monitoring

### Container Durumu

```bash
# Çalışan container'ları listele
docker ps

# Container detayları
docker inspect p2p-api-1

# Container içine gir
docker exec -it p2p-api-1 sh
```

### Log Yönetimi

```bash
# Tüm logları görüntüle
docker-compose logs

# Belirli servisin logları
docker-compose logs p2p-api-1

# Canlı log takibi
docker-compose logs -f p2p-api-1
```

### Health Check

```bash
# Health check durumu
docker-compose ps

# Health check logları
docker inspect p2p-api-1 | grep Health -A 10
```

## 🛠️ Geliştirme

### Kod Değişiklikleri

```bash
# Geliştirme ortamında çalıştır
docker-compose -f docker-compose.dev.yml up -d

# Kod değişikliklerini otomatik yükle
# (Volume mount sayesinde hot reload)
```

### Yeni Versiyon Build

```bash
# Yeni imaj build et
docker-compose build

# Servisleri yeniden başlat
docker-compose up -d
```

## 🧹 Temizlik

```bash
# Container'ları durdur ve sil
docker-compose down

# İmajları da sil
docker-compose down --rmi all

# Volume'ları da sil
docker-compose down -v

# Tüm Docker kaynaklarını temizle
docker system prune -a
```

## ⚠️ Önemli Notlar

1. **Port Çakışması**: Eğer port 8080, 8081, 8082 kullanımdaysa, `docker-compose.yml` dosyasında port mapping'leri değiştirin.

2. **Network**: Tüm container'lar aynı Docker network'ünde çalışır, bu sayede birbirlerini bulabilirler.

3. **Persistent Data**: Container'lar silindiğinde veriler kaybolur. Kalıcı veri için volume kullanın.

4. **Resource Limits**: Büyük ölçekli deployment için resource limitleri ekleyin.

## 🚀 Production Deployment

### Kubernetes için

```yaml
# deployment.yaml örneği
apiVersion: apps/v1
kind: Deployment
metadata:
  name: p2p-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: p2p-api
  template:
    metadata:
      labels:
        app: p2p-api
    spec:
      containers:
      - name: p2p-api
        image: p2p-api:latest
        ports:
        - containerPort: 8080
```

### Environment Variables

```bash
# Production ortamı için
export GO_ENV=production
export PORT=8080

docker run -e GO_ENV=production -e PORT=8080 -p 8080:8080 p2p-api
``` 