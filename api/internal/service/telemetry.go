package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"spm-api/internal/websocket"
)

type ProcessInfo struct {
	PID                  int     `json:"pid"`
	ExecutableName       string  `json:"executable_name"`
	ResourceUtilization  float64 `json:"resource_utilization"`
}

type TelemetryPayload struct {
	AgentID            string        `json:"agent_id"`
	CPUUtilization     float64       `json:"cpu_utilization"`
	MemoryUtilization  float64       `json:"memory_utilization"`
	DiskIO             float64       `json:"disk_io"`
	NetworkIngress     float64       `json:"network_ingress"`
	NetworkEgress      float64       `json:"network_egress"`
	Temperature        float64       `json:"temperature"`
	Uptime             int           `json:"uptime"`
	Status             string        `json:"status"`
	Timestamp          string        `json:"timestamp,omitempty"`
	TopProcesses       []ProcessInfo `json:"top_processes,omitempty"`
}

type TelemetryService struct {
	db        *sql.DB
	redis     *redis.Client
	wsManager *websocket.Manager
}

func NewTelemetryService(db *sql.DB, redisClient *redis.Client, wsManager *websocket.Manager) *TelemetryService {
	return &TelemetryService{
		db:        db,
		redis:     redisClient,
		wsManager: wsManager,
	}
}

func (s *TelemetryService) ProcessTelemetry(payload TelemetryPayload) error {
	if s.db != nil {
		_, err := s.db.Exec(`
			INSERT INTO telemetry (agent_id, cpu_utilization, memory_utilization, disk_io, network_ingress, network_egress, temperature, uptime, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			payload.AgentID, payload.CPUUtilization, payload.MemoryUtilization, payload.DiskIO, payload.NetworkIngress, payload.NetworkEgress, payload.Temperature, payload.Uptime, payload.Status,
		)
		if err != nil {
			log.Printf("DB Insert Error: %v", err)
		} else {
			for _, p := range payload.TopProcesses {
				_, err = s.db.Exec(`
					INSERT INTO processes (agent_id, pid, executable_name, resource_utilization)
					VALUES ($1, $2, $3, $4)`,
					payload.AgentID, p.PID, p.ExecutableName, p.ResourceUtilization,
				)
				if err != nil {
					log.Printf("Process DB Insert Error: %v", err)
				}
			}
		}
	}

	payloadBytes, err := json.Marshal(payload)
	if err == nil {
		if s.redis != nil {
			ctx := context.Background()
			s.redis.Set(ctx, "agent:"+payload.AgentID+":latest", payloadBytes, 24*time.Hour)
		}
		s.wsManager.Broadcast(payloadBytes)
	}
	return nil
}

func (s *TelemetryService) GetAllAgents() ([]TelemetryPayload, error) {
	if s.redis == nil { return []TelemetryPayload{}, nil }
	ctx := context.Background()
	keys, err := s.redis.Keys(ctx, "agent:*:latest").Result()
	if err != nil { return nil, err }

	var agents []TelemetryPayload
	for _, key := range keys {
		val, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var payload TelemetryPayload
			if err := json.Unmarshal([]byte(val), &payload); err == nil {
				agents = append(agents, payload)
			}
		}
	}
	return agents, nil
}

func (s *TelemetryService) GetHistory(agentID string, limit int) []TelemetryPayload {
	var history []TelemetryPayload
	if s.db == nil { return history }
	
	rows, err := s.db.Query(`
		SELECT agent_id, cpu_utilization, memory_utilization, disk_io, network_ingress, network_egress, temperature, uptime, status, timestamp
		FROM telemetry WHERE agent_id = $1 ORDER BY timestamp DESC LIMIT $2
	`, agentID, limit)
	if err != nil { return history }
	defer rows.Close()
	
	for rows.Next() {
		var p TelemetryPayload
		var ts time.Time
		if err := rows.Scan(&p.AgentID, &p.CPUUtilization, &p.MemoryUtilization, &p.DiskIO, &p.NetworkIngress, &p.NetworkEgress, &p.Temperature, &p.Uptime, &p.Status, &ts); err == nil {
			p.Timestamp = ts.Format(time.RFC3339)
			history = append(history, p)
		}
	}
	return history
}

func (s *TelemetryService) GenerateBackupCSV() string {
	if s.db == nil { return "agent_id,cpu,memory,temperature,status,timestamp\n" }
	
	rows, err := s.db.Query(`SELECT agent_id, cpu_utilization, memory_utilization, temperature, status, timestamp FROM telemetry ORDER BY timestamp DESC LIMIT 5000`)
	if err != nil { return "Error generating backup\n" }
	defer rows.Close()
	
	var sb strings.Builder
	sb.WriteString("agent_id,cpu_utilization,memory_utilization,temperature,status,timestamp\n")
	
	for rows.Next() {
		var agentID, status string
		var cpu, mem, temp float64
		var ts time.Time
		if err := rows.Scan(&agentID, &cpu, &mem, &temp, &status, &ts); err == nil {
			sb.WriteString(fmt.Sprintf("%s,%.2f,%.2f,%.2f,%s,%s\n", agentID, cpu, mem, temp, status, ts.Format(time.RFC3339)))
		}
	}
	return sb.String()
}

func (s *TelemetryService) CheckHealth() (string, string) {
	dbStatus := "down"
	if s.db != nil && s.db.Ping() == nil {
		dbStatus = "up"
	}
	redisStatus := "down"
	if s.redis != nil && s.redis.Ping(context.Background()).Err() == nil {
		redisStatus = "up"
	}
	return dbStatus, redisStatus
}

func (s *TelemetryService) StartDatabasePruner() {
	if s.db == nil { return }
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			// Keep only latest 100 telemetry logs per agent
			_, err := s.db.Exec(`
				DELETE FROM telemetry 
				WHERE id IN (
					SELECT id FROM (
						SELECT id, ROW_NUMBER() OVER(PARTITION BY agent_id ORDER BY timestamp DESC) as rn
						FROM telemetry
					) t WHERE t.rn > 100
				)
			`)
			if err != nil {
				log.Printf("DB Prune Error (Telemetry): %v", err)
			}
			
			// Keep only latest 500 process logs per agent
			_, err = s.db.Exec(`
				DELETE FROM processes 
				WHERE id IN (
					SELECT id FROM (
						SELECT id, ROW_NUMBER() OVER(PARTITION BY agent_id ORDER BY timestamp DESC) as rn
						FROM processes
					) t WHERE t.rn > 500
				)
			`)
			if err != nil {
				log.Printf("DB Prune Error (Processes): %v", err)
			}
		}
	}()
}

