package service

import (
	"context"
	"time"

	pb "github.com/SaidovZohid/note_user_service/genproto/user_service"
	"github.com/SaidovZohid/note_user_service/storage"
	"github.com/SaidovZohid/note_user_service/storage/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
}

func NewUserService(strg storage.StorageI, inMemory storage.InMemoryStorageI) *UserService {
	return &UserService{
		storage:  strg,
		inMemory: inMemory,
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Create(parseUser(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		ImageUrl:    user.ImageUrl,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) Get(ctx context.Context, req *pb.IdRequest) (*pb.User, error) {
	user, err := s.storage.User().Get(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		ImageUrl:    user.ImageUrl,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) GetByEmail(ctx context.Context, req *pb.EmailRequest) (*pb.User, error) {
	user, err := s.storage.User().GetByEmail(req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		ImageUrl:    user.ImageUrl,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) Update(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Update(&repo.User{
		ID:          req.Id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		ImageUrl:    req.ImageUrl,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &pb.User{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		ImageUrl:    user.ImageUrl,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) Delete(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	err := s.storage.User().Delete(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) GetAll(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.storage.User().GetAll(&repo.GetAllUsersParams{
		Limit:  req.Limit,
		Page:   req.Page,
		Search: req.Search,
		SortBy: req.SortBy,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}
	res := pb.GetAllUsersResponse{
		Count: users.Count,
		Users: make([]*pb.User, 0),
	}
	for _, user := range users.Users {
		u := pb.User{
			Id:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Username:    user.Username,
			Email:       user.Email,
			Password:    user.Password,
			PhoneNumber: user.PhoneNumber,
			ImageUrl:    user.ImageUrl,
			Type:        user.Type,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		}
		res.Users = append(res.Users, &u)
	}

	return &res, nil
}

func parseUser(user *pb.User) *repo.User {
	return &repo.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		ImageUrl:    user.ImageUrl,
		Type:        user.Type,
	}
}
