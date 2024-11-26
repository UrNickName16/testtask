package graph

import (
	"context"

	"testtask/internal/server/graph/generated"
	"testtask/internal/server/graph/model"
	"testtask/internal/stats"
)

type Resolver struct {
	StatsService *stats.Service
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) NetworkInterfaces(ctx context.Context) ([]*model.NetInterface, error) {
	interfaces := r.StatsService.GetInterfaces()
	var result []*model.NetInterface
	for _, iface := range interfaces {
		result = append(result, &model.NetInterface{
			Name:            iface.Name,
			Config:          &iface.Config,
			LinkUp:          iface.LinkUp,
			PacketsSent:     int(iface.PacketsSent),
			PacketsReceived: int(iface.PacketsReceived),
			BytesSent:       int(iface.BytesSent),
			BytesReceived:   int(iface.BytesReceived),
			SpeedSent:       iface.SpeedSent,
			SpeedReceived:   iface.SpeedReceived,
		})
	}
	return result, nil
}
