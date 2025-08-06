package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	// Ana context oluştur - tüm işlemler için temel context
	ctx := context.Background()

	// Libp2p node oluştur - P2P ağında bu node'u temsil eder
	// Bu host, diğer peer'larla iletişim kurmak için gerekli
	host, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// Bu node'un benzersiz kimliğini (PeerID) yazdır
	// Bu ID, ağda bu node'u tanımlamak için kullanılır
	fmt.Println("Publisher PeerID:", host.ID().String())
	
	// Bu node'un tüm adreslerini yazdır
	// Diğer peer'lar bu adresler üzerinden bağlanabilir
	for _, addr := range host.Addrs() {
		fmt.Printf("Address: %s/p2p/%s\n", addr, host.ID().String())
	}

	// PubSub sistemi başlat - GossipSub protokolü kullanarak
	// Bu, mesajların peer'lar arasında yayılmasını sağlar
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	// Kullanıcıdan subscriber'ların adreslerini al
	// Bu adresler, mesaj göndereceğimiz peer'ların adresleri
	fmt.Println("Enter first subscriber multiaddr:")
	scanner := bufio.NewScanner(os.Stdin)
	
	// İlk subscriber'ın multiaddr'ını al ve bağlan
	scanner.Scan()
	line1 := strings.TrimSpace(scanner.Text())
	if line1 != "" {
		// Girilen string'i multiaddr formatına çevir
		maddr1, err := ma.NewMultiaddr(line1)
		if err != nil {
			fmt.Println("Invalid first address:", err)
		} else {
			// Multiaddr'dan peer bilgilerini çıkar
			info1, err := peer.AddrInfoFromP2pAddr(maddr1)
			if err != nil {
				fmt.Println("Invalid first peer info:", err)
			} else {
				// İlk peer'a bağlan
				if err := host.Connect(ctx, *info1); err != nil {
					fmt.Println("Failed to connect to first peer:", err)
				} else {
					fmt.Println("Connected to first peer:", info1.ID.String())
				}
			}
		}
	}
	
	// İkinci subscriber'ın multiaddr'ını al ve bağlan
	fmt.Println("Enter second subscriber multiaddr:")
	scanner.Scan()
	line2 := strings.TrimSpace(scanner.Text())
	if line2 != "" {
		// Girilen string'i multiaddr formatına çevir
		maddr2, err := ma.NewMultiaddr(line2)
		if err != nil {
			fmt.Println("Invalid second address:", err)
		} else {
			// Multiaddr'dan peer bilgilerini çıkar
			info2, err := peer.AddrInfoFromP2pAddr(maddr2)
			if err != nil {
				fmt.Println("Invalid second peer info:", err)
			} else {
				// İkinci peer'a bağlan
				if err := host.Connect(ctx, *info2); err != nil {
					fmt.Println("Failed to connect to second peer:", err)
				} else {
					fmt.Println("Connected to second peer:", info2.ID.String())
				}
			}
		}
	}

	// Belirli bir topic'e (konu) katıl
	// "account-lookup" topic'i üzerinden mesaj göndereceğiz
	topic, err := ps.Join("account-lookup")
	if err != nil {
		panic(err)
	}

	// Kullanıcıdan mesaj alıp yayınlama döngüsü
	fmt.Println("Start sending messages. Type 'exit' to quit.")
	for {
		fmt.Print("Enter account number to broadcast: ")
		scanner.Scan()
		msg := scanner.Text()
		
		// Eğer kullanıcı 'exit' yazarsa döngüden çık
		if strings.TrimSpace(msg) == "exit" {
			break
		}
		
		// Mesajı topic'e yayınla - tüm abonelere gönderilir
		topic.Publish(ctx, []byte(msg))
		
		// Bir saniye bekle - mesajın işlenmesi için zaman ver
		time.Sleep(1 * time.Second)
	}
}