# Go-libp2p PubSub Mesajlaşma Sistemi

Bu proje, Go-libp2p kullanarak peer-to-peer (P2P) mesajlaşma sistemi oluşturan bir örnek uygulamadır. Sistem, bir **Publisher** (yayıncı) ve bir veya daha fazla **Subscriber** (abone) arasında mesaj alışverişi sağlar.

## 🚀 Özellikler

- **P2P İletişim**: Merkezi sunucu olmadan peer-to-peer iletişim
- **PubSub Pattern**: Yayıncı-abone mesajlaşma modeli
- **GossipSub Protokolü**: Güvenilir ve ölçeklenebilir mesaj dağıtımı
- **Çoklu Bağlantı**: Publisher birden fazla subscriber'a bağlanabilir
- **Gerçek Zamanlı**: Anlık mesaj iletimi

## 📁 Proje Yapısı

```
go-api/
├── publisher.go    # Mesaj yayınlayan uygulama
├── subscriber.go   # Mesaj alan uygulama
├── go.mod          # Go modül bağımlılıkları
├── go.sum          # Bağımlılık hash'leri
└── README.md       # Bu dosya
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

## 🎯 Kullanım

### 1. Subscriber'ları Başlatın

İlk olarak, mesaj alacak subscriber'ları başlatın:

```bash
go run subscriber.go
```

Her subscriber çalıştığında şu bilgileri göreceksiniz:
```
Subscriber PeerID: 12D3KooW...
Address: /ip4/127.0.0.1/tcp/xxxxx/p2p/12D3KooW...
Subscribed to topic 'account-lookup'
```

**Önemli**: Subscriber'ların gösterdiği multiaddr'ları not edin, bunları publisher'da kullanacaksınız.

### 2. Publisher'ı Başlatın

Subscriber'lar çalıştıktan sonra publisher'ı başlatın:

```bash
go run publisher.go
```

Publisher başladığında:
```
Publisher PeerID: 12D3KooW...
Address: /ip4/127.0.0.1/tcp/xxxxx/p2p/12D3KooW...
Enter first subscriber multiaddr:
```

### 3. Bağlantıları Kurun

Publisher size iki subscriber multiaddr'ı soracak:

1. **İlk subscriber multiaddr'ını girin** (örnek):
   ```
   /ip4/127.0.0.1/tcp/60734/p2p/12D3KooWD5wDQYpjDKmZqAQKdbMufDXvM6JaQUm2ZH4WcsC8T1fd
   ```

2. **İkinci subscriber multiaddr'ını girin** (örnek):
   ```
   /ip4/127.0.0.1/tcp/60726/p2p/12D3KooWBRRAD6Pkq6xNm6sm5rpsherPAfEARa3GtnCoZyruhBLg
   ```

### 4. Mesaj Gönderin

Bağlantılar kurulduktan sonra, publisher size mesaj göndermenizi isteyecek:

```
Start sending messages. Type 'exit' to quit.
Enter account number to broadcast: 
```

Burada gönderdiğiniz her mesaj, bağlı tüm subscriber'lara iletilecektir.

## 🔧 Teknik Detaylar

### Kullanılan Teknolojiler

- **libp2p**: P2P ağ protokolü
- **GossipSub**: PubSub protokolü
- **Multiaddr**: Ağ adresi formatı
- **Peer ID**: Benzersiz peer tanımlayıcısı

### Mesajlaşma Akışı

1. **Bağlantı Kurma**: Publisher, subscriber'ların multiaddr'larını kullanarak bağlantı kurar
2. **Topic Aboneliği**: Tüm peer'lar "account-lookup" topic'ine abone olur
3. **Mesaj Yayını**: Publisher mesaj gönderdiğinde, GossipSub protokolü mesajı tüm abonelere dağıtır
4. **Mesaj Alma**: Subscriber'lar gelen mesajları alır ve ekrana yazdırır

### Protokol Detayları

- **Topic**: "account-lookup" (sabit topic adı)
- **Mesaj Formatı**: Plain text (byte array)
- **Bağlantı Tipi**: TCP üzerinden P2P
- **Keşif**: Manuel multiaddr girişi

## 🐛 Sorun Giderme

### Yaygın Hatalar

1. **"Invalid address" Hatası**
   - Multiaddr formatını kontrol edin
   - Boşluk veya özel karakter olmadığından emin olun

2. **"Failed to connect" Hatası**
   - Subscriber'ın çalıştığından emin olun
   - Firewall ayarlarını kontrol edin
   - Multiaddr'ın doğru olduğunu kontrol edin

3. **Mesaj Alınmıyor**
   - Topic adının aynı olduğunu kontrol edin
   - Bağlantının başarılı olduğunu kontrol edin

### Debug İpuçları

- Her peer'ın benzersiz bir PeerID'si vardır
- Multiaddr formatı: `/ip4/IP/tcp/PORT/p2p/PEER_ID`
- Subscriber'lar önce başlatılmalı, sonra publisher

## 📝 Örnek Kullanım Senaryosu

### Senaryo: Hesap Arama Sistemi

Bu sistem, finansal kurumlar arasında hesap bilgisi paylaşımı için kullanılabilir:

1. **Subscriber'lar**: Farklı bankalar (Bank A, Bank B)
2. **Publisher**: Merkezi hesap sorgulama sistemi
3. **Mesajlar**: Hesap numaraları
4. **Sonuç**: Tüm bankalar aynı anda hesap sorgularını alır

## 🔮 Gelecek Geliştirmeler

- [ ] Otomatik peer keşfi
- [ ] Şifreli mesajlaşma
- [ ] Mesaj doğrulama
- [ ] Web arayüzü
- [ ] Çoklu topic desteği
- [ ] Mesaj geçmişi

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