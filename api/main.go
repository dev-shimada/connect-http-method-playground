package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	apiv1 "github.com/dev-shimada/connect-http-method-playground/api/gen/proto/api/v1"
	"github.com/dev-shimada/connect-http-method-playground/api/gen/proto/api/v1/apiv1connect"
	"github.com/google/uuid"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	Host = "0.0.0.0"
	Port = 8081
)

type server struct {
	*http.Server
	repository MessageRepository
}

// domain model
type Message struct {
	ID     string
	UserID string
	Text   string
}

// repository interface
type MessageRepository interface {
	Save(ctx context.Context, msg Message) error
	Get(ctx context.Context, id string) (*Message, error)
}

var repo = make(map[string]Message)

type Repository struct{}

func NewRepository() MessageRepository {
	return &Repository{}
}

func (r *Repository) Save(ctx context.Context, msg Message) error {
	repo[msg.ID] = msg
	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*Message, error) {
	message, exists := repo[id]
	if !exists {
		return nil, fmt.Errorf("message not found")
	}
	return &message, nil
}

func (s *server) Post(ctx context.Context, req *connect.Request[apiv1.PostRequest]) (*connect.Response[apiv1.PostResponse], error) {
	slog.Info("Received Post request", "user", req.Msg.UserId, "text", req.Msg.Text)
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate UUID: %w", err))
	}

	msg := Message{
		ID:     id.String(),
		UserID: req.Msg.UserId,
		Text:   req.Msg.Text,
	}

	if err := s.repository.Save(ctx, msg); err != nil {
		slog.Error("Failed to save message", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to save message: %w", err))
	}

	return &connect.Response[apiv1.PostResponse]{
		Msg: &apiv1.PostResponse{
			Id: id.String(),
		},
	}, nil
}

func (s *server) Get(ctx context.Context, req *connect.Request[apiv1.GetRequest]) (*connect.Response[apiv1.GetResponse], error) {
	slog.Info("Received Get request", "id", req.Msg.Id)

	msg, err := s.repository.Get(ctx, req.Msg.Id)
	if err != nil {
		slog.Error("Failed to get message", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("message not found: %w", err))
	}

	return &connect.Response[apiv1.GetResponse]{
		Msg: &apiv1.GetResponse{
			UserId: msg.UserID,
			Text:   msg.Text,
		},
	}, nil
}

func (s *server) Rest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		msg, err := s.repository.Get(r.Context(), r.URL.Query().Get("id"))
		if err != nil {
			slog.Error("Failed to get message", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		res := &apiv1.GetResponse{
			UserId: msg.UserID,
			Text:   msg.Text,
		}
		w.WriteHeader(http.StatusOK)
		if bytes, err := json.Marshal(res); err == nil {
			w.Write(bytes)
		} else {
			slog.Error("Failed to marshal response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var req apiv1.PostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("Failed to decode request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := uuid.NewString()
		if err := s.repository.Save(r.Context(), Message{
			ID:     id,
			UserID: req.UserId,
			Text:   req.Text,
		}); err != nil {
			slog.Error("Failed to save message", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := &apiv1.PostResponse{
			Id: id,
		}
		w.WriteHeader(http.StatusOK)
		if bytes, err := json.Marshal(res); err == nil {
			w.Write(bytes)
		} else {
			slog.Error("Failed to marshal response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s *server) PingPong(ctx context.Context, req *connect.Request[apiv1.PingPongRequest]) (*connect.Response[apiv1.PingPongResponse], error) {
	res := connect.NewResponse(&apiv1.PingPongResponse{
		UserId: req.Msg.UserId,
		Text:   req.Msg.Text,
	})
	res.Header().Set("Message-Version", "v1")
	return res, nil
}

func main() {
	// json logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(log.Writer(), nil)))

	mux := http.NewServeMux()

	// reflection
	reflector := grpcreflect.NewStaticReflector(
		apiv1connect.ApiServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// health check
	checker := grpchealth.NewStaticChecker(
		apiv1connect.ApiServiceName,
	)
	mux.Handle(grpchealth.NewHandler(checker))

	svc := &server{
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", Host, Port),
			Handler: h2c.NewHandler(mux, &http2.Server{}),
		},
		repository: NewRepository(),
	}

	// api
	path, handler := apiv1connect.NewApiServiceHandler(svc)
	mux.Handle(path, handler)

	// http/1.1 base endpoint
	mux.HandleFunc("/api/v1/ApiService", svc.Rest)

	// http/1.1 pingpong endpoint
	mux.HandleFunc("/api/v1/ApiService/PingPong", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			res := &apiv1.PingPongResponse{
				UserId: r.URL.Query().Get("user_id"),
				Text:   r.URL.Query().Get("text"),
			}
			if bytes, err := json.Marshal(res); err == nil {
				w.Write(bytes)
			} else {
				slog.Error("Failed to marshal response", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})

	// start server
	signalCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	slog.Info(fmt.Sprintf("Server is running at %s:%d", Host, Port))
	go func() {
		if err := svc.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("Server closed.")
			} else {
				slog.Error(fmt.Sprintf("Failed to serve: %v", err))
			}
		}
	}()
	<-signalCtx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	slog.Info("Shutting down server...")
	if err := svc.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("Failed to shutdown server: %v", err))
	}
}
