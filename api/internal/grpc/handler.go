package grpc

import (
	"context"
	"io"
	"log"
	
	"spm-api/internal/service"
	pb "spm-api/proto"
)

type TelemetryServer struct {
	pb.UnimplementedTelemetryServiceServer
	svc *service.TelemetryService
}

func NewTelemetryServer(svc *service.TelemetryService) *TelemetryServer {
	return &TelemetryServer{svc: svc}
}

func mapToServicePayload(req *pb.TelemetryPayload) service.TelemetryPayload {
	procs := make([]service.ProcessInfo, len(req.TopProcesses))
	for i, p := range req.TopProcesses {
		procs[i] = service.ProcessInfo{
			PID:                 int(p.Pid),
			ExecutableName:      p.ExecutableName,
			ResourceUtilization: p.ResourceUtilization,
		}
	}
	return service.TelemetryPayload{
		AgentID:           req.AgentId,
		CPUUtilization:    req.CpuUtilization,
		MemoryUtilization: req.MemoryUtilization,
		DiskIO:            req.DiskIo,
		NetworkIngress:    req.NetworkIngress,
		NetworkEgress:     req.NetworkEgress,
		Temperature:       req.Temperature,
		Uptime:            int(req.Uptime),
		Status:            req.Status,
		TopProcesses:      procs,
	}
}

func (s *TelemetryServer) SendTelemetry(ctx context.Context, req *pb.TelemetryPayload) (*pb.TelemetryResponse, error) {
	payload := mapToServicePayload(req)
	err := s.svc.ProcessTelemetry(payload)
	if err != nil {
		return &pb.TelemetryResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.TelemetryResponse{Success: true, Message: "Telemetry received via gRPC"}, nil
}

func (s *TelemetryServer) StreamTelemetry(stream pb.TelemetryService_StreamTelemetryServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.TelemetryResponse{Success: true, Message: "Stream closed"})
		}
		if err != nil {
			log.Printf("Error receiving stream: %v", err)
			return err
		}
		
		payload := mapToServicePayload(req)
		s.svc.ProcessTelemetry(payload)
	}
}

func (s *TelemetryServer) CheckHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	dbStatus, redisStatus := s.svc.CheckHealth()
	return &pb.HealthResponse{
		Database: dbStatus,
		Redis:    redisStatus,
	}, nil
}
