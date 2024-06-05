package grpc

import (
	"clientManage/internal/domain/model"
	"clientManage/internal/service"
	pb "clientManage/proto"
	"context"
)

type ClientGRPCServer struct {
	pb.UnimplementedClientManagementServiceServer
	clientService *service.ClientService
}

func NewClientGRPCServer(clientService *service.ClientService) *ClientGRPCServer {
	return &ClientGRPCServer{clientService: clientService}
}

func (s *ClientGRPCServer) CreateClient(ctx context.Context, req *pb.CreateClientRequest) (*pb.CreateClientResponse, error) {
	client := &model.Client{
		ID:        req.Client.Id,
		Name:      req.Client.Fname,
		Surname:   req.Client.Sname,
		Email:     req.Client.Email,
		Role:      req.Client.UserRole,
		Activated: req.Client.Activated,
		Version:   int(req.Client.Version),
	}

	err := s.clientService.CreateClient(client)
	if err != nil {
		return nil, err
	}

	return &pb.CreateClientResponse{
		Client: &pb.Client{
			Id:        client.ID,
			Fname:     client.Name,
			Sname:     client.Surname,
			Email:     client.Email,
			UserRole:  client.Role,
			Activated: client.Activated,
			Version:   int32(client.Version),
		},
	}, nil
}

func (s *ClientGRPCServer) GetClient(ctx context.Context, req *pb.GetClientRequest) (*pb.GetClientResponse, error) {
	client, err := s.clientService.GetClientByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetClientResponse{
		Client: &pb.Client{
			Id:        client.ID,
			Fname:     client.Name,
			Sname:     client.Surname,
			Email:     client.Email,
			UserRole:  client.Role,
			Activated: client.Activated,
			Version:   int32(client.Version),
		},
	}, nil
}

func (s *ClientGRPCServer) UpdateClient(ctx context.Context, req *pb.UpdateClientRequest) (*pb.UpdateClientResponse, error) {
	client := &model.Client{
		ID:        req.Client.Id,
		Name:      req.Client.Fname,
		Surname:   req.Client.Sname,
		Email:     req.Client.Email,
		Role:      req.Client.UserRole,
		Activated: req.Client.Activated,
		Version:   int(req.Client.Version),
	}

	err := s.clientService.UpdateClient(client)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateClientResponse{
		Client: &pb.Client{
			Id:        client.ID,
			Fname:     client.Name,
			Sname:     client.Surname,
			Email:     client.Email,
			UserRole:  client.Role,
			Activated: client.Activated,
			Version:   int32(client.Version),
		},
	}, nil
}

func (s *ClientGRPCServer) DeleteClient(ctx context.Context, req *pb.DeleteClientRequest) (*pb.DeleteClientResponse, error) {
	err := s.clientService.DeleteClient(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteClientResponse{
		Client: &pb.Client{
			Id: req.Id,
		},
	}, nil
}

func (s *ClientGRPCServer) ListClients(ctx context.Context, req *pb.ListClientsRequest) (*pb.ListClientsResponse, error) {
	filters := make(map[string]interface{})
	for _, filter := range req.Filters {
		filters[filter.Key] = filter.Value
	}

	clients, err := s.clientService.ListClients(int(req.Offset), int(req.Limit), filters, req.SortBy, req.SortOrder)
	if err != nil {
		return nil, err
	}

	var clientResponses []*pb.Client
	for _, client := range clients {
		clientResponses = append(clientResponses, &pb.Client{
			Id:        client.ID,
			Fname:     client.Name,
			Sname:     client.Surname,
			Email:     client.Email,
			UserRole:  client.Role,
			Activated: client.Activated,
			Version:   int32(client.Version),
		})
	}

	return &pb.ListClientsResponse{
		Clients: clientResponses,
	}, nil
}
