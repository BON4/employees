package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echo "github.com/labstack/echo/v4"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
	PORT           = ":8080"
	READ_TIMEOUT   = 5
	WRITE_TIMEOUT  = 5
)

type Server struct {
	e      *echo.Echo
	logger *log.Logger
}

func NewServer() (*Server, error) {
	e := echo.New()
	return &Server{
		e:      e,
		logger: log.New(os.Stdout, "⇨ ", log.Flags()),
	}, nil
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           PORT,
		ReadTimeout:    time.Second * READ_TIMEOUT,
		WriteTimeout:   time.Second * WRITE_TIMEOUT,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Printf("Server is listening on PORT: %s", PORT)
		if err := s.e.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	if err := s.MapHandlers(s.e); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Printf("Server Exited Properly")
	return s.e.Server.Shutdown(ctx)
}
