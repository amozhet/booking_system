package grpc

import (
	"booking/internal/domain/model"
	"booking/internal/service"
	pb "booking/proto"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingGRPCServer struct {
	pb.UnimplementedBookingServiceServer
	bookingService *service.BookingService
}

func NewBookingGRPCServer(bookingService *service.BookingService) *BookingGRPCServer {
	return &BookingGRPCServer{bookingService: bookingService}
}

func (s *BookingGRPCServer) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	booking := &model.Booking{
		ClientID:  req.ClientId,
		RoomID:    req.RoomId,
		StartDate: req.StartDate.AsTime(),
		EndDate:   req.EndDate.AsTime(),
		Status:    req.Status,
	}

	err := s.bookingService.CreateBooking(booking)
	if err != nil {
		return nil, err
	}

	return &pb.BookingResponse{
		Id:        booking.ID,
		ClientId:  booking.ClientID,
		RoomId:    booking.RoomID,
		StartDate: timestamppb.New(booking.StartDate),
		EndDate:   timestamppb.New(booking.EndDate),
		Status:    booking.Status,
	}, nil
}

func (s *BookingGRPCServer) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingResponse, error) {
	booking, err := s.bookingService.GetBookingByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.BookingResponse{
		Id:        booking.ID,
		ClientId:  booking.ClientID,
		RoomId:    booking.RoomID,
		StartDate: timestamppb.New(booking.StartDate),
		EndDate:   timestamppb.New(booking.EndDate),
		Status:    booking.Status,
	}, nil
}

func (s *BookingGRPCServer) UpdateBooking(ctx context.Context, req *pb.UpdateBookingRequest) (*pb.BookingResponse, error) {
	booking := &model.Booking{
		ID:        req.Id,
		ClientID:  req.ClientId,
		RoomID:    req.RoomId,
		StartDate: req.StartDate.AsTime(),
		EndDate:   req.EndDate.AsTime(),
		Status:    req.Status,
	}

	err := s.bookingService.UpdateBooking(booking)
	if err != nil {
		return nil, err
	}

	return &pb.BookingResponse{
		Id:        booking.ID,
		ClientId:  booking.ClientID,
		RoomId:    booking.RoomID,
		StartDate: timestamppb.New(booking.StartDate),
		EndDate:   timestamppb.New(booking.EndDate),
		Status:    booking.Status,
	}, nil
}

/*
func (s *BookingGRPCServer) DeleteBooking(ctx context.Context, req *pb.DeleteBookingRequest) (*pb.Empty, error) {
	err := s.bookingService.DeleteBooking(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
*/

func (s *BookingGRPCServer) ListBookings(ctx context.Context, req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error) {
	filters := make(map[string]interface{})
	for _, filter := range req.Filters {
		filters[filter.Key] = filter.Value
	}

	bookings, err := s.bookingService.ListBookings(int(req.Offset), int(req.Limit), filters, req.SortBy, req.SortOrder)
	if err != nil {
		return nil, err
	}

	var bookingResponses []*pb.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, &pb.BookingResponse{
			Id:        booking.ID,
			ClientId:  booking.ClientID,
			RoomId:    booking.RoomID,
			StartDate: timestamppb.New(booking.StartDate),
			EndDate:   timestamppb.New(booking.EndDate),
			Status:    booking.Status,
		})
	}

	return &pb.ListBookingsResponse{
		Bookings: bookingResponses,
	}, nil
}
