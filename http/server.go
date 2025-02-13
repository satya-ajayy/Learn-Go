package http

import (
	// Go Internal Packages
	"context"
	"net/http"
	"time"

	// Local Packages
	errors "learn-go/errors"
	handlers "learn-go/http/handlers"
	resp "learn-go/http/response"
	xhandlers "learn-go/models/xhandlers"
	health "learn-go/services/health"

	// External Packages
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"moul.io/chizap"
)

// Server struct follow the alphabet order
type Server struct {
	health   *health.HealthCheckerService
	logger   *zap.Logger
	logs     *handlers.LogsHandler
	prefix   string
	students *handlers.StudentsHandler
}

func NewServer(
	prefix string,
	logger *zap.Logger,
	h *xhandlers.XHandlers,
	healthCheckService *health.HealthCheckerService,
) *Server {
	return &Server{
		prefix:   prefix,
		logger:   logger,
		students: h.StudentsHandlers,
		logs:     h.LogsHandlers,
		health:   healthCheckService,
	}
}

func (s *Server) Listen(ctx context.Context, addr string) error {
	r := chi.NewRouter()
	r.Use(chizap.New(s.logger, &chizap.Opts{
		WithReferer:   false,
		WithUserAgent: false,
	}))
	r.Use(middleware.Recoverer)

	r.Route(s.prefix, func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/health", s.HealthCheckHandler)

			r.Group(func(r chi.Router) {
				r.Route("/students", func(r chi.Router) {
					r.Get("/", s.ToHTTPHandlerFunc(s.students.GetAll))
					r.Get("/{rollNo}", s.ToHTTPHandlerFunc(s.students.GetOne))
					r.Post("/", s.ToHTTPHandlerFunc(s.students.Insert))
					r.Put("/{rollNo}", s.ToHTTPHandlerFunc(s.students.Update))
					r.Delete("/{rollNo}", s.ToHTTPHandlerFunc(s.students.Delete))
				})
				r.Route("/logs", func(r chi.Router) {
					r.Post("/", s.ToHTTPHandlerFunc(s.logs.Insert))
				})
			})
		})
	})

	errch := make(chan error)
	server := &http.Server{Addr: addr, Handler: r}
	go func() {
		s.logger.Info("Starting server", zap.String("addr", addr))
		errch <- server.ListenAndServe()
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}

// ToHTTPHandlerFunc converts a handler function to a http.HandlerFunc.
// This wrapper function is used to handle errors and respond to the client
func (s *Server) ToHTTPHandlerFunc(handler func(w http.ResponseWriter, r *http.Request) (any, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, status, err := handler(w, r)
		if err != nil {
			switch err := err.(type) {
			case *errors.Error:
				resp.RespondError(w, err)
			default:
				s.logger.Error("internal error", zap.Error(err))
				resp.RespondMessage(w, http.StatusInternalServerError, "internal error")
			}
			return
		}
		if response != nil {
			resp.RespondJSON(w, status, response)
		}
		if status >= 100 && status < 600 {
			w.WriteHeader(status)
		}
	}
}

// HealthCheckHandler returns the health status of the service
func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if ok := s.health.Health(r.Context()); !ok {
		resp.RespondMessage(w, http.StatusServiceUnavailable, "health check failed")
		return
	}
	resp.RespondMessage(w, http.StatusOK, "!!! We are RunninGoo !!!")
}
