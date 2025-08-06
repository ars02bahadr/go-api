# 🔗 P2P REST API - Go-libp2p

Bu proje, Go-libp2p kullanarak peer-to-peer (P2P) mesajlaşma sistemini REST API olarak sunan bir uygulamadır. Sistem, HTTP endpoint'leri üzerinden P2P işlemlerini yönetmenizi sağlar.

## 🚀 Özellikler

- **REST API**: HTTP endpoint'leri üzerinden P2P işlemleri
- **Web Arayüzü**: Kullanıcı dostu test arayüzü
- **P2P İletişim**: Merkezi sunucu olmadan peer-to-peer iletişim
- **PubSub Pattern**: Yayıncı-abone mesajlaşma modeli
- **GossipSub Protokolü**: Güvenilir ve ölçeklenebilir mesaj dağıtımı
- **Gerçek Zamanlı**: Server-Sent Events ile canlı mesaj dinleme
- **Çoklu Bağlantı**: Birden fazla peer'a bağlanabilme

## 📁 Proje Yapısı

```
go-api/
├── main.go              # Ana API sunucusu
├── static/
│   └── index.html       # Web test arayüzü
├── publisher.go         # Eski P2P publisher (referans)
├── subscriber.go        # Eski P2P subscriber (referans)
├── go.mod              # Go modül bağımlılıkları
├── go.sum              # Bağımlılık hash'leri
└── README.md           # Bu dosya
```

## 🛠️ Kurulum

### Gereksinimler

- Go 1.19 veya üzeri
- Git

### Kurulum Adımları

1. **Projeyi klonlayın:**
   ```bash
   git clone <repository-url>
   cd go-api
   ```

2. **Bağımlılıkları yükleyin:**
   ```bash
   go mod tidy
   ```

3. **API'yi başlatın:**
   ```bash
   go run main.go
   ```

4. **Web arayüzüne erişin:**
   ```
   http://localhost:8080
   ```

## 🎯 API Endpoint'leri

### 📊 Durum Kontrolü
```http
GET /api/status
```

**Yanıt:**
```json
{
  "success": true,
  "message": "Node durumu",
  "data": {
    "peer_id": "12D3KooW...",
    "tcp_addresses": ["/ip4/127.0.0.1/tcp/8080/p2p/..."],
    "connected_peers": 2,
    "message_count": 5
  }
}
```

### 🔌 Peer Bağlantısı
```http
POST /api/connect
Content-Type: application/json

{
  "multiaddr": "/ip4/127.0.0.1/tcp/8080/p2p/12D3KooW..."
}
```

**Yanıt:**
```json
{
  "success": true,
  "message": "Peer'a başarıyla bağlanıldı: 12D3KooW...",
  "data": {
    "peer_id": "12D3KooW...",
    "address": "/ip4/127.0.0.1/tcp/8080/p2p/..."
  }
}
```

### 👥 Bağlı Peer'lar
```http
GET /api/peers
```

**Yanıt:**
```json
{
  "success": true,
  "message": "Bağlı peer'lar",
  "data": [
    {
      "peer_id": "12D3KooW...",
      "addrs": ["/ip4/127.0.0.1/tcp/8080"]
    }
  ]
}
```

### 📤 Mesaj Gönderme
```http
POST /api/send
Content-Type: application/json

{
  "content": "Merhaba P2P dünyası!"
}
```

**Yanıt:**
```json
{
  "success": true,
  "message": "Mesaj başarıyla yayınlandı",
  "data": {
    "id": "msg_1",
    "content": "Merhaba P2P dünyası!",
    "from": "12D3KooW...",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### 📨 Mesaj Geçmişi
```http
GET /api/messages
```

**Yanıt:**
```json
{
  "success": true,
  "message": "Alınan mesajlar",
  "data": [
    {
      "id": "msg_1",
      "content": "Merhaba!",
      "from": "12D3KooW...",
      "timestamp": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 📡 Canlı Mesaj Dinleme
```http
GET /api/subscribe
```

**Server-Sent Events (SSE) ile gerçek zamanlı mesaj dinleme**

## 🌐 Web Arayüzü

API'yi test etmek için web arayüzü kullanabilirsiniz:

1. **Tarayıcıda açın:** `http://localhost:8080`
2. **Durumu kontrol edin:** API'nin çalışıp çalışmadığını görün
3. **Peer bağlantısı:** Multiaddr girerek peer'a bağlanın
4. **Mesaj gönderin:** Metin kutusuna mesaj yazıp gönderin
5. **Canlı dinleme:** Yeni mesajları gerçek zamanlı görün

## 🔧 Kullanım Örnekleri

### cURL ile Test

```bash
# Durum kontrolü
curl http://localhost:8080/api/status

# Peer bağlantısı
curl -X POST http://localhost:8080/api/connect \
  -H "Content-Type: application/json" \
  -d '{"multiaddr": "/ip4/127.0.0.1/tcp/8080/p2p/12D3KooW..."}'

# Mesaj gönderme
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d '{"content": "Test mesajı"}'

# Mesajları listele
curl http://localhost:8080/api/messages
```

### JavaScript ile Test

```javascript
// Durum kontrolü
const status = await fetch('/api/status').then(r => r.json());

// Peer bağlantısı
const connect = await fetch('/api/connect', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ multiaddr: '/ip4/127.0.0.1/tcp/8080/p2p/...' })
}).then(r => r.json());

// Mesaj gönderme
const send = await fetch('/api/send', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ content: 'Merhaba!' })
}).then(r => r.json());

// Canlı dinleme
const eventSource = new EventSource('/api/subscribe');
eventSource.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Yeni mesaj:', message);
};
```

## 🔮 Gelişmiş Özellikler

### Çoklu Node Test

1. **İlk API'yi başlatın:**
   ```bash
   PORT=8080 go run main.go
   ```

2. **İkinci API'yi başlatın:**
   ```bash
   PORT=8081 go run main.go
   ```

3. **Peer'ları bağlayın:**
   - İlk API'nin durumunu kontrol edin
   - İkinci API'den ilk API'ye bağlanın
   - Mesaj gönderin ve test edin

### Environment Variables

- `PORT`: API port'u (varsayılan: 8080)

## 🐛 Sorun Giderme

### Yaygın Hatalar

1. **"Port already in use"**
   ```bash
   # Farklı port kullanın
   PORT=8081 go run main.go
   ```

2. **"Failed to connect"**
   - Peer'ın çalıştığından emin olun
   - Multiaddr'ın doğru olduğunu kontrol edin
   - Firewall ayarlarını kontrol edin

3. **"Invalid multiaddr"**
   - Multiaddr formatını kontrol edin
   - Boşluk veya özel karakter olmadığından emin olun

### Debug İpuçları

- Her node'un benzersiz bir PeerID'si vardır
- Multiaddr formatı: `/ip4/IP/tcp/PORT/p2p/PEER_ID`
- Web arayüzünde tüm işlemleri test edebilirsiniz

## 📝 Teknik Detaylar

### Kullanılan Teknolojiler

- **Go**: Ana programlama dili
- **libp2p**: P2P ağ protokolü
- **GossipSub**: PubSub protokolü
- **HTTP**: REST API
- **Server-Sent Events**: Gerçek zamanlı iletişim
- **HTML/JavaScript**: Web arayüzü

### Mimari

```
┌─────────────────┐    HTTP    ┌─────────────────┐
│   Web Client    │ ────────── │   REST API      │
└─────────────────┘            └─────────────────┘
                                         │
                                         ▼
                                ┌─────────────────┐
                                │   P2P Node      │
                                │  (libp2p)       │
                                └─────────────────┘
                                         │
                                         ▼
                                ┌─────────────────┐
                                │   GossipSub     │
                                │   (PubSub)      │
                                └─────────────────┘
```

## 🔮 Gelecek Geliştirmeler

- [ ] WebSocket desteği
- [ ] Şifreli mesajlaşma
- [ ] Mesaj doğrulama
- [ ] Çoklu topic desteği
- [ ] Mesaj geçmişi veritabanı
- [ ] Docker desteği
- [ ] Kubernetes deployment
- [ ] Prometheus metrics
- [ ] JWT authentication

## 📄 Lisans

Bu proje MIT lisansı altında lisanslanmıştır.

## 🤝 Katkıda Bulunma

1. Fork yapın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Commit yapın (`git commit -m 'Add amazing feature'`)
4. Push yapın (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## 📞 İletişim

Sorularınız için issue açabilir veya iletişime geçebilirsiniz. 