package gateaway

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alperdrsnn/clime"
	"github.com/gorilla/websocket"
	"github.com/hasan-kilici/chat/internal/gateaway/ws"
	"github.com/hasan-kilici/chat/internal/service/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  10 * 1024 * 1024, // 10 MB
}

func Start() error {
	repository.Connect()

	clime.Header("🚀 Chat Gateway Server")

	clime.InfoLine(fmt.Sprintf("WebSocket endpoint: %s", clime.BoldColor.Sprint("/ws")))
	clime.InfoLine(fmt.Sprintf("Static files served from: %s", clime.BoldColor.Sprint("/static")))

	clime.SuccessBanner("Starting server at " + clime.BoldColor.Sprint("127.0.0.1:5000") + " 🚀")

	clime.NewBox().
		WithTitle("Endpoints").
		WithBorderColor(clime.GreenColor).
		AddLine(fmt.Sprintf("%s  - WebSocket endpoint", clime.BoldColor.Sprint("/ws"))).
		AddLine(fmt.Sprintf("%s - Static file server", clime.BoldColor.Sprint("/static"))).
		Println()

	serverErr := make(chan error)

	go func() {
		http.HandleFunc("/ws", ws.WebSocketHandler)
		http.Handle("/", http.FileServer(http.Dir("./static")))
		serverErr <- http.ListenAndServe("127.0.0.1:5000", nil)
	}()

	time.Sleep(200 * time.Millisecond)

	err := <-serverErr
	if err != nil {
		clime.ErrorBanner("Failed to start server: " + err.Error())
		return err
	}

	return nil
}
