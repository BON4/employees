package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BON4/employees/pkg/store"
	"github.com/go-redis/redis/v8"
	echo "github.com/labstack/echo/v4"
)

const (
	maxHeaderBytes    = 1 << 20
	ctxTimeout        = 5
	PORT              = ":8080"
	READ_TIMEOUT      = 5
	WRITE_TIMEOUT     = 5
	REDIS_CONN_STRING = "redis://admin:@localhost:6379/0"
)

type Server struct {
	e      *echo.Echo
	logger *log.Logger
	st     store.Store
}

func NewServer() (*Server, error) {
	e := echo.New()

	opt, err := redis.ParseURL(REDIS_CONN_STRING)
	if err != nil {
		return nil, err
	}

	customLogger := log.New(os.Stdout, "⇨ ", log.Flags())

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Server{
		e:      e,
		st:     store.NewRedisStore(rdb),
		logger: customLogger,
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
			s.logger.Fatalf("Error starting Server: %s", err.Error())
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

	if err := s.st.Close(); err != nil {
		s.logger.Printf("Server Exited with err: %s", err.Error())
	} else {
		s.logger.Printf("Server Exited Properly")
	}
	return s.e.Server.Shutdown(ctx)
}
