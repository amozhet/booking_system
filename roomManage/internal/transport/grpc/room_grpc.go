package grpc

import (
	"context"
	"roomManage/internal/domain/model"
	"roomManage/internal/service"
	"roomManage/proto"
)

type RoomGRPCServer struct {
	service *service.RoomService
	proto.UnimplementedRoomServiceServer
}

func NewRoomGRPCServer(service *service.RoomService) *RoomGRPCServer {
	return &RoomGRPCServer{
		service: service,
	}
}

func (s *RoomGRPCServer) GetRooms(ctx context.Context, req *proto.GetRoomsRequest) (*proto.GetRoomsResponse, error) {
	rooms, err := s.service.GetAllRooms()
	if err != nil {
		return nil, err
	}
	var protoRooms []*proto.Room
	for _, room := range rooms {
		protoRooms = append(protoRooms, &proto.Room{
			Id:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Available:   room.Available,
		})
	}
	return &proto.GetRoomsResponse{Rooms: protoRooms}, nil
}

func (s *RoomGRPCServer) GetRoomByID(ctx context.Context, req *proto.GetRoomByIDRequest) (*proto.GetRoomResponse, error) {
	room, err := s.service.GetRoomByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.GetRoomResponse{Room: &proto.Room{
		Id:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Available:   room.Available,
	}}, nil
}

func (s *RoomGRPCServer) CreateRoom(ctx context.Context, req *proto.CreateRoomRequest) (*proto.RoomResponse, error) {
	room := &model.Room{
		ID:          req.Room.Id,
		Name:        req.Room.Name,
		Description: req.Room.Description,
		Available:   req.Room.Available,
	}
	err := s.service.CreateRoom(room)
	if err != nil {
		return nil, err
	}
	return &proto.RoomResponse{Room: req.Room}, nil
}

func (s *RoomGRPCServer) UpdateRoom(ctx context.Context, req *proto.UpdateRoomRequest) (*proto.RoomResponse, error) {
	room := &model.Room{
		ID:          req.Room.Id,
		Name:        req.Room.Name,
		Description: req.Room.Description,
		Available:   req.Room.Available,
	}
	err := s.service.UpdateRoom(req.Room.Id, room)
	if err != nil {
		return nil, err
	}
	return &proto.RoomResponse{Room: req.Room}, nil
}

func (s *RoomGRPCServer) DeleteRoom(ctx context.Context, req *proto.DeleteRoomRequest) (*proto.DeleteRoomResponse, error) {
	err := s.service.DeleteRoom(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteRoomResponse{Success: true}, nil
}
