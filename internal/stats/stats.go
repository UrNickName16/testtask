package stats

import (
	"log"
	"net"
	"sync"
	"time"

	netstat "github.com/shirou/gopsutil/net"
)

type InterfaceConfig struct {
	HardwareAddr string
	MTU          int
	Flags        []string
	Addrs        []string
}

type InterfaceStats struct {
	Name            string
	Config          InterfaceConfig
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
	s.updateStats()
	go s.start()
	return s
}

func (s *Service) start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		s.updateStats()
	}
}

func getInterfaceConfig(ifaceName string) InterfaceConfig {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Printf("Error getting interface %s: %v", ifaceName, err)
		return InterfaceConfig{}
	}

	addrs, err := iface.Addrs()
	if err != nil {
		log.Printf("Error getting addresses for interface %s: %v", ifaceName, err)
		return InterfaceConfig{}
	}

	var addrStrings []string
	for _, addr := range addrs {
		addrStrings = append(addrStrings, addr.String())
	}

	var flags []string
	if iface.Flags&net.FlagUp != 0 {
		flags = append(flags, "up")
	}
	if iface.Flags&net.FlagBroadcast != 0 {
		flags = append(flags, "broadcast")
	}
	if iface.Flags&net.FlagLoopback != 0 {
		flags = append(flags, "loopback")
	}
	if iface.Flags&net.FlagPointToPoint != 0 {
		flags = append(flags, "pointtopoint")
	}
	if iface.Flags&net.FlagMulticast != 0 {
		flags = append(flags, "multicast")
	}

	return InterfaceConfig{
		HardwareAddr: iface.HardwareAddr.String(),
		MTU:          iface.MTU,
		Flags:        flags,
		Addrs:        addrStrings,
	}
}

func isInterfaceUp(name string) bool {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Printf("Error getting interface %s: %v", name, err)
		return false
	}
	return iface.Flags&net.FlagUp != 0
}

func (s *Service) updateStats() {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentStats := make(map[string]*InterfaceStats)

	netIOCounters, err := netstat.IOCounters(true)
	if err != nil {
		log.Printf("Error getting net IO counters: %v", err)
		return
	}

	for _, counter := range netIOCounters {
		ifaceName := counter.Name

		ifaceConfig := getInterfaceConfig(ifaceName)

		ifaceStat := &InterfaceStats{
			Name:            ifaceName,
			Config:          ifaceConfig,
			BytesReceived:   counter.BytesRecv,
			PacketsReceived: counter.PacketsRecv,
			BytesSent:       counter.BytesSent,
			PacketsSent:     counter.PacketsSent,
			LinkUp:          isInterfaceUp(ifaceName),
		}

		if prev, ok := s.prevStats[ifaceName]; ok {
			duration := 5.0 // Период обновления в секундах
			ifaceStat.SpeedReceived = float64(ifaceStat.BytesReceived-prev.BytesReceived) / duration
			ifaceStat.SpeedSent = float64(ifaceStat.BytesSent-prev.BytesSent) / duration
		}

		currentStats[ifaceName] = ifaceStat
	}

	s.prevStats = s.interfaces
	s.interfaces = currentStats

	// Логирование
	for name, iface := range s.interfaces {
		log.Printf("Interface %s: %+v", name, iface)
	}
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
