package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	listenAddr   string
	upgrader     websocket.Upgrader
	sseBroadcast chan []byte
	sseClients   map[*sseClient]bool
	sseMutex     sync.Mutex
}

type sseClient struct {
	writer http.ResponseWriter
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		sseBroadcast: make(chan []byte),
		sseClients:   make(map[*sseClient]bool),
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/ping", makeHTTPHandleFunc(s.handlePing)).Methods("GET")
	router.HandleFunc("/api/ws", makeHTTPHandleFunc(s.handleWs)).Methods("GET")
	router.HandleFunc("/api/sse", makeHTTPHandleFunc(s.handleSse)).Methods("GET")
	router.Use(s.loggingMiddleware)

	s.sseBroadcast = make(chan []byte)
	go func() {
		for {
			msg, ok := <-s.sseBroadcast
			if !ok {
				return
			}
			s.sseMutex.Lock()
			for client := range s.sseClients {
				_, err := client.writer.Write(msg)
				if err != nil {
					delete(s.sseClients, client) // Remove client on error
					continue
				}
				if flusher, ok := client.writer.(http.Flusher); ok {
					flusher.Flush()
				}
			}
			s.sseMutex.Unlock()
		}
	}()

	return http.ListenAndServe(s.listenAddr, router)
}

type MiddlewareFunc func(http.Handler) http.Handler

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("pong"))
	return nil
}

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) error {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			formattedMessage := fmt.Sprintf("data: %s\n\n", message)
			s.sseBroadcast <- []byte(formattedMessage)
		}
	}
	return nil
}

func (s *Server) handleSse(w http.ResponseWriter, r *http.Request) error {
	client := &sseClient{writer: w}

	s.sseMutex.Lock()
	s.sseClients[client] = true
	s.sseMutex.Unlock()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Send an initial message to confirm the connection
	fmt.Fprintf(w, "data: %s\n\n", "connected")
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	} else {
		log.Println("Unable to cast to Flusher")
	}

	// Keep the connection open until the client closes it
	<-r.Context().Done()
	s.sseMutex.Lock()
	delete(s.sseClients, client)
	s.sseMutex.Unlock()

	return nil
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
