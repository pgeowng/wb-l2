package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Context struct {
	Res http.ResponseWriter
	Req *http.Request
}

func (ctx *Context) SendError(statusCode int) {
	ctx.Res.WriteHeader(statusCode)
}

func (ctx *Context) SendJSON(statusCode int, data H) {
	ctx.Res.WriteHeader(statusCode)
	ctx.Res.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("json marshal err. Err: %s, for request %v", err, ctx.Req.RequestURI)
	}
	ctx.Res.Write(bytes)
}

type Handler = func(Context)
type Middleware = func(Handler) Handler
type H = map[string]interface{}

func LoggerMW(next func(Context)) func(Context) {
	return func(ctx Context) {
		fmt.Println("req:", ctx.Req.URL)
		next(ctx)
	}
}

type Server struct {
	http  *http.Server
	mux   *http.ServeMux
	paths map[string]map[string]Handler // exactpath/method/handler
}

func New(address string) *Server {
	h := &http.Server{Addr: address}
	mux := http.NewServeMux()
	h.Handler = mux

	srv := &Server{
		http:  h,
		mux:   mux,
		paths: map[string]map[string]Handler{},
	}

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		srv.matchPath(Context{Res: rw, Req: r})
	})

	return srv
}

func (s *Server) Listen() error {

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) matchPath(ctx Context) {
	handlers, ok := s.paths[ctx.Req.URL.Path]
	if !ok {
		ctx.SendError(http.StatusNotFound)
		return
	}

	handler, ok := handlers[ctx.Req.Method]
	if !ok {
		ctx.SendError(http.StatusNotFound)
		return
	}

	if err := ctx.Req.ParseForm(); err != nil {
		fmt.Println(err)
		ctx.SendError(http.StatusInternalServerError)
		return
	}

	handler(ctx)
}

func (s *Server) Get(path string, end Handler, mw ...Middleware) {
	for idx := len(mw) - 1; idx >= 0; idx-- {
		end = mw[idx](end)
	}

	if s.paths[path] == nil {
		s.paths[path] = map[string]Handler{}
	}
	s.paths[path]["GET"] = end
}

func (s *Server) Post(path string, end Handler, mw ...Middleware) {
	for idx := len(mw) - 1; idx >= 0; idx-- {
		end = mw[idx](end)
	}

	if s.paths[path] == nil {
		s.paths[path] = map[string]Handler{}
	}
	s.paths[path]["POST"] = end
}
