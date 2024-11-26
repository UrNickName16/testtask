package graph

import "testtask/internal/stats"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	StatsService *stats.Service
}
