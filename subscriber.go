package main

import (
	"context"
	"fmt"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func main() {
	// Ana context oluştur - tüm işlemler için temel context
	ctx := context.Background()

	// Libp2p host oluştur - P2P ağında bu node'u temsil eder
	// Bu host, diğer peer'larla iletişim kurmak için gerekli
	host, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// Bu node'un benzersiz kimliğini (PeerID) yazdır
	// Bu ID, ağda bu node'u tanımlamak için kullanılır
	fmt.Println("Subscriber PeerID:", host.ID().String())
	
	// Sadece TCP adreslerini göster - diğer protokoller çok fazla çıktı veriyor
	// TCP en yaygın ve güvenilir protokoldür
	fmt.Println("Available TCP addresses:")
	for _, addr := range host.Addrs() {
		addrStr := addr.String()
		// Sadece TCP içeren adresleri filtrele
		if strings.Contains(addrStr, "/tcp/") {
			fmt.Printf("  %s/p2p/%s\n", addrStr, host.ID().String())
		}
	}

	// PubSub sistemi başlat - GossipSub protokolü kullanarak
	// Bu, mesajların peer'lar arasında yayılmasını sağlar
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	// Belirli bir topic'e (konu) katıl
	// "account-lookup" topic'i üzerinden mesaj alacağız
	topic, err := ps.Join("account-lookup")
	if err != nil {
		panic(err)
	}

	// Topic'e abone ol - bu topic'teki tüm mesajları alacağız
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	fmt.Println("\nSubscribed to topic 'account-lookup'")
	fmt.Println("Waiting for messages...")

	// Sonsuz döngü ile gelen mesajları dinle
	for {
		// Bir sonraki mesajı bekle
		msg, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		// Mesajı ve gönderen peer'ı yazdır
		fmt.Printf("Received Message from %s: %s\n", msg.ReceivedFrom.String(), string(msg.Data))
	}
}