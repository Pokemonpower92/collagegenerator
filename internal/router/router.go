package router

import (
	"log/slog"
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/logger"
	"github.com/pokemonpower92/collagegenerator/internal/middleware"
)

type ApiFunc func(http.ResponseWriter, *http.Request, *slog.Logger) error

type ApiError struct {
	Error string
}

func makeHttpHandler(h ApiFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, _ := logger.GetRequestLogger(r)
		if err := h(w, r, l); err != nil {
			panic(err)
		}
	})
}

type Router struct {
	Mux *http.ServeMux
}

func NewRouter() *Router {
	sm := http.NewServeMux()
	return &Router{Mux: sm}
}

func (r *Router) RegisterRoute(path string, handler ApiFunc) {
	handlerFunc := makeHttpHandler(handler)
	stdMiddleware := middleware.New(
		middleware.Context(),
		middleware.Logger(),
		middleware.Error(),
	)
	handlerFunc = stdMiddleware.Use(handlerFunc)
	r.Mux.HandleFunc(path, handlerFunc.ServeHTTP)
}
