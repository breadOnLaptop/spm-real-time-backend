package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"io"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/proto"

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
	go wsManager.Run()

	telemetryService := service.NewTelemetryService(dbConn, redisClient, wsManager)
	telemetryService.StartDatabasePruner()

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

	http.HandleFunc("/api/telemetry/binary", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		if r.Method == "OPTIONS" { return }
		
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		var protoPayload pb.TelemetryPayload
		if err := proto.Unmarshal(body, &protoPayload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payload := grpc_internal.MapToServicePayload(&protoPayload)
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

	// Start Multiplexed gRPC & HTTP Server
	grpcServer := grpc.NewServer()
	grpcHandler := grpc_internal.NewTelemetryServer(telemetryService)
	pb.RegisterTelemetryServiceServer(grpcServer, grpcHandler)

	mux := http.DefaultServeMux
	mixedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s | Proto: %d.%d | Content-Type: %s", r.Method, r.URL.Path, r.ProtoMajor, r.ProtoMinor, r.Header.Get("Content-Type"))
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})

	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: h2c.NewHandler(mixedHandler, h2s),
	}

	log.Printf("Multiplexed Server (HTTP/gRPC) starting on port %s", cfg.Port)
	if err := h1s.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func setupCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
