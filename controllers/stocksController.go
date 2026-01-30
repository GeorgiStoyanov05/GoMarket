package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// For dev. In prod, restrict origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

type finnhubSubMsg struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}

// GET /ws/trades?symbol=AAPL
func WSFinnhubTrades(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	token := os.Getenv("FINNHUB_API_KEY")
	if token == "" {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Upgrade browser -> your server
	clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer clientConn.Close()

	// Connect your server -> Finnhub WS
	fhURL := "wss://ws.finnhub.io?token=" + token
	fhConn, _, err := websocket.DefaultDialer.Dial(fhURL, nil)
	if err != nil {
		log.Println("Finnhub WS dial error:", err)
		return
	}
	defer fhConn.Close()

	// Subscribe
	if err := fhConn.WriteJSON(finnhubSubMsg{Type: "subscribe", Symbol: symbol}); err != nil {
		log.Println("Finnhub subscribe error:", err)
		return
	}

	// Keep-alive (helps avoid idle closes)
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = clientConn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(2*time.Second))
			case <-done:
				return
			}
		}
	}()

	// Forward Finnhub messages -> browser
	for {
		_, msg, err := fhConn.ReadMessage()
		if err != nil {
			close(done)
			return
		}

		// Optional: validate it's JSON (helps avoid weird frames)
		var tmp any
		if err := json.Unmarshal(msg, &tmp); err != nil {
			continue
		}

		if err := clientConn.WriteMessage(websocket.TextMessage, msg); err != nil {
			close(done)
			return
		}
	}
}


func GetSymbolDetailsPage(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.String(http.StatusBadRequest, "missing symbol")
		return
	}

	// HTMX request -> return partial
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "symbolDetails", middlewares.WithAuth(c, gin.H{
			"Symbol":     symbol,
			"DefaultRes": "5",
		}))
		return
	}

	// Normal navigation (typed URL) -> return shell + InitialPath
	c.HTML(http.StatusOK, "index.html", middlewares.WithAuth(c, gin.H{
		"InitialPath": "/details/" + symbol,
	}))
}
