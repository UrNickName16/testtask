package stats

import (
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/net"
)

type InterfaceStats struct {
	Name            string
	Config          string
	LinkUp          bool
	PacketsSent     uint64
	PacketsReceived uint64
	BytesSent       uint64
	BytesReceived   uint64
	SpeedSent       float64
	SpeedReceived   float64
}

type Service struct {
	mu         sync.RWMutex
	interfaces map[string]*InterfaceStats
	prevStats  map[string]*InterfaceStats
}

func NewStatsService() *Service {
	s := &Service{
		interfaces: make(map[string]*InterfaceStats),
		prevStats:  make(map[string]*InterfaceStats),
	}
	go s.start()
	return s
}

func (s *Service) start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		s.updateStats()
		<-ticker.C
	}
}

func (s *Service) updateStats() {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentStats := make(map[string]*InterfaceStats)

	netIOCounters, err := net.IOCounters(true)
	if err != nil {
		log.Printf("Error getting net IO counters: %v", err)
		return
	}

	for _, counter := range netIOCounters {
		ifaceName := counter.Name

		ifaceStat := &InterfaceStats{
			Name:            ifaceName,
			BytesReceived:   counter.BytesRecv,
			PacketsReceived: counter.PacketsRecv,
			BytesSent:       counter.BytesSent,
			PacketsSent:     counter.PacketsSent,
			LinkUp:          true,
		}

		if prev, ok := s.prevStats[ifaceName]; ok {
			duration := 5.0
			ifaceStat.SpeedReceived = float64(counter.BytesRecv-prev.BytesReceived) / duration
			ifaceStat.SpeedSent = float64(counter.BytesSent-prev.BytesSent) / duration
		}

		currentStats[ifaceName] = ifaceStat
	}

	s.prevStats = s.interfaces
	s.interfaces = currentStats
}

func (s *Service) GetInterfaces() []*InterfaceStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*InterfaceStats
	for _, iface := range s.interfaces {
		result = append(result, iface)
	}
	return result
}
