package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net"

	"google.golang.org/grpc"
	grpc_internal "spm-api/internal/grpc"
	pb "spm-api/proto"

	"spm-api/internal/config"
	"spm-api/internal/db"
	"spm-api/internal/service"
	"spm-api/internal/websocket"
)

func main() {
	cfg := config.Load()

	dbConn, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Printf("DB connection failed: %v", err)
	}
	
	redisClient, err := db.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Printf("Redis connection failed: %v", err)
	}

	wsManager := websocket.NewManager()
	telemetryService := service.NewTelemetryService(dbConn, redisClient, wsManager)

	// WebSockets
	http.HandleFunc("/api/ws", wsManager.HandleConnection)

	// Initialize HTTP Handlers
	http.HandleFunc("/api/telemetry", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		if r.Method == "OPTIONS" { return }
		var payload service.TelemetryPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		telemetryService.ProcessTelemetry(payload)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/agents", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		if r.Method == "OPTIONS" { return }
		agents, _ := telemetryService.GetAllAgents()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(agents)
	})

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		if r.Method == "OPTIONS" { return }
		dbStatus := "down"
		if dbConn != nil && dbConn.Ping() == nil { dbStatus = "up" }
		redisStatus := "down"
		if redisClient != nil && redisClient.Ping(context.Background()).Err() == nil { redisStatus = "up" }
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"database": dbStatus, "redis": redisStatus})
	})

	http.HandleFunc("/api/history", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		if r.Method == "OPTIONS" { return }
		agentID := r.URL.Query().Get("agent_id")
		history := telemetryService.GetHistory(agentID, 50)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)
	})

	http.HandleFunc("/api/backup", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		if r.Method == "OPTIONS" { return }
		csv := telemetryService.GenerateBackupCSV()
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=telemetry_backup.csv")
		w.Write([]byte(csv))
	})

	http.HandleFunc("/ws", wsManager.HandleConnection)

	// Start gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on 50051: %v", err)
	}
	grpcServer := grpc.NewServer()
	grpcHandler := grpc_internal.NewTelemetryServer(telemetryService)
	pb.RegisterTelemetryServiceServer(grpcServer, grpcHandler)
	
	go func() {
		log.Println("gRPC Server starting on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP Server
	log.Printf("HTTP Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func setupCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
