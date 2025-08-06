package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// P2PNode P2P ağ node'unu temsil eder
type P2PNode struct {
	Host     host.Host
	PubSub   *pubsub.PubSub
	Topic    *pubsub.Topic
	Sub      *pubsub.Subscription
	mu       sync.RWMutex
	peers    map[string]peer.AddrInfo
	messages []Message
}

// Message mesaj yapısını temsil eder
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	From      string    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
}

// APIResponse API yanıt yapısı
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ConnectRequest bağlantı isteği yapısı
type ConnectRequest struct {
	Multiaddr string `json:"multiaddr"`
}

// SendMessageRequest mesaj gönderme isteği yapısı
type SendMessageRequest struct {
	Content string `json:"content"`
}

var (
	p2pNode *P2PNode
	msgID   = 0
	msgMu   sync.Mutex
)

func main() {
	// P2P node'u başlat
	if err := initP2PNode(); err != nil {
		log.Fatal("P2P node başlatılamadı:", err)
	}

	// HTTP sunucusunu başlat
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	// Statik dosyaları serve et
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// API route'larını tanımla
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
		} else {
			homeHandler(w, r)
		}
	})
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/connect", connectHandler)
	http.HandleFunc("/api/peers", peersHandler)
	http.HandleFunc("/api/send", sendMessageHandler)
	http.HandleFunc("/api/messages", messagesHandler)
	http.HandleFunc("/api/subscribe", subscribeHandler)

	// Graceful shutdown için sinyal dinle
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API sunucusu başlatılıyor: http://localhost%s", port)
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal("HTTP sunucusu başlatılamadı:", err)
		}
	}()

	// Sinyal bekle
	<-sigChan
	log.Println("Sunucu kapatılıyor...")
	
	// P2P node'u temizle
	if p2pNode != nil {
		if p2pNode.Sub != nil {
			p2pNode.Sub.Cancel()
		}
		if p2pNode.Host != nil {
			p2pNode.Host.Close()
		}
	}
}

// initP2PNode P2P node'unu başlatır
func initP2PNode() error {
	ctx := context.Background()

	// Libp2p host oluştur
	host, err := libp2p.New()
	if err != nil {
		return fmt.Errorf("host oluşturulamadı: %v", err)
	}

	// PubSub başlat
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return fmt.Errorf("pubsub başlatılamadı: %v", err)
	}

	// Topic'e katıl
	topic, err := ps.Join("account-lookup")
	if err != nil {
		return fmt.Errorf("topic'e katılınamadı: %v", err)
	}

	// Topic'e abone ol
	sub, err := topic.Subscribe()
	if err != nil {
		return fmt.Errorf("topic'e abone olunamadı: %v", err)
	}

	p2pNode = &P2PNode{
		Host:     host,
		PubSub:   ps,
		Topic:    topic,
		Sub:      sub,
		peers:    make(map[string]peer.AddrInfo),
		messages: make([]Message, 0),
	}

	// Mesaj dinleme goroutine'ini başlat
	go listenForMessages(ctx)

	log.Printf("P2P Node başlatıldı - PeerID: %s", host.ID().String())
	return nil
}

// listenForMessages gelen mesajları dinler
func listenForMessages(ctx context.Context) {
	for {
		msg, err := p2pNode.Sub.Next(ctx)
		if err != nil {
			log.Printf("Mesaj dinleme hatası: %v", err)
			continue
		}

		// Mesajı kaydet
		message := Message{
			ID:        fmt.Sprintf("msg_%d", getNextMsgID()),
			Content:   string(msg.Data),
			From:      msg.ReceivedFrom.String(),
			Timestamp: time.Now(),
		}

		p2pNode.mu.Lock()
		p2pNode.messages = append(p2pNode.messages, message)
		p2pNode.mu.Unlock()

		log.Printf("Yeni mesaj alındı: %s (gönderen: %s)", message.Content, message.From)
	}
}

// getNextMsgID benzersiz mesaj ID'si oluşturur
func getNextMsgID() int {
	msgMu.Lock()
	defer msgMu.Unlock()
	msgID++
	return msgID
}

// HTTP Handler'ları

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "P2P API Sunucusu Çalışıyor",
		Data: map[string]interface{}{
			"peer_id": p2pNode.Host.ID().String(),
			"version": "1.0.0",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	p2pNode.mu.RLock()
	defer p2pNode.mu.RUnlock()

	// TCP adreslerini topla
	var tcpAddresses []string
	for _, addr := range p2pNode.Host.Addrs() {
		addrStr := addr.String()
		if strings.Contains(addrStr, "/tcp/") {
			tcpAddresses = append(tcpAddresses, fmt.Sprintf("%s/p2p/%s", addrStr, p2pNode.Host.ID().String()))
		}
	}

	response := APIResponse{
		Success: true,
		Message: "Node durumu",
		Data: map[string]interface{}{
			"peer_id":        p2pNode.Host.ID().String(),
			"tcp_addresses":  tcpAddresses,
			"connected_peers": len(p2pNode.peers),
			"message_count":   len(p2pNode.messages),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Sadece POST metodu desteklenir", http.StatusMethodNotAllowed)
		return
	}

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	if req.Multiaddr == "" {
		http.Error(w, "Multiaddr gerekli", http.StatusBadRequest)
		return
	}

	// Multiaddr'ı parse et
	maddr, err := ma.NewMultiaddr(req.Multiaddr)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: fmt.Sprintf("Geçersiz multiaddr: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Peer bilgilerini çıkar
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: fmt.Sprintf("Peer bilgileri çıkarılamadı: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Peer'a bağlan
	ctx := context.Background()
	if err := p2pNode.Host.Connect(ctx, *info); err != nil {
		response := APIResponse{
			Success: false,
			Message: fmt.Sprintf("Peer'a bağlanılamadı: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Peer'ı kaydet
	p2pNode.mu.Lock()
	p2pNode.peers[info.ID.String()] = *info
	p2pNode.mu.Unlock()

	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("Peer'a başarıyla bağlanıldı: %s", info.ID.String()),
		Data: map[string]interface{}{
			"peer_id": info.ID.String(),
			"address": req.Multiaddr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func peersHandler(w http.ResponseWriter, r *http.Request) {
	p2pNode.mu.RLock()
	defer p2pNode.mu.RUnlock()

	var peers []map[string]interface{}
	for id, info := range p2pNode.peers {
		peers = append(peers, map[string]interface{}{
			"peer_id": id,
			"addrs":   info.Addrs,
		})
	}

	response := APIResponse{
		Success: true,
		Message: "Bağlı peer'lar",
		Data:    peers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Sadece POST metodu desteklenir", http.StatusMethodNotAllowed)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Mesaj içeriği gerekli", http.StatusBadRequest)
		return
	}

	// Mesajı topic'e yayınla
	ctx := context.Background()
	if err := p2pNode.Topic.Publish(ctx, []byte(req.Content)); err != nil {
		response := APIResponse{
			Success: false,
			Message: fmt.Sprintf("Mesaj yayınlanamadı: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Mesajı kaydet
	message := Message{
		ID:        fmt.Sprintf("msg_%d", getNextMsgID()),
		Content:   req.Content,
		From:      p2pNode.Host.ID().String(),
		Timestamp: time.Now(),
	}

	p2pNode.mu.Lock()
	p2pNode.messages = append(p2pNode.messages, message)
	p2pNode.mu.Unlock()

	response := APIResponse{
		Success: true,
		Message: "Mesaj başarıyla yayınlandı",
		Data:    message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	p2pNode.mu.RLock()
	defer p2pNode.mu.RUnlock()

	response := APIResponse{
		Success: true,
		Message: "Alınan mesajlar",
		Data:    p2pNode.messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// Server-Sent Events için header'ları ayarla
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Mesaj kanalı oluştur
	msgChan := make(chan Message, 10)
	defer close(msgChan)

	// Yeni mesajları dinle
	go func() {
		for {
			time.Sleep(1 * time.Second)
			p2pNode.mu.RLock()
			if len(p2pNode.messages) > 0 {
				latestMsg := p2pNode.messages[len(p2pNode.messages)-1]
				msgChan <- latestMsg
			}
			p2pNode.mu.RUnlock()
		}
	}()

	// Mesajları client'a gönder
	for {
		select {
		case msg := <-msgChan:
			data, _ := json.Marshal(msg)
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
} 