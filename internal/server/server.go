package server

import (
	"net/http"

	"testtask/internal/server/graph"
	"testtask/internal/server/graph/generated"
	"testtask/internal/stats"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	resolver := &graph.Resolver{
		StatsService: stats.NewStatsService(),
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	return &Server{
		httpServer: &http.Server{
			Addr: ":8080",
		},
	}
}

func (s *Server) Run(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}
