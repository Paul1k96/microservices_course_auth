package mapper

import (
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	desc "github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserFromCreateRequest converts api request to user model.
func ToUserFromCreateRequest(request *desc.CreateRequest) *model.User {
	return &model.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     model.Role(request.Role),
	}
}

// ToCreateResponseFromUserService converts user id to api response.
func ToCreateResponseFromUserService(userID int64) *desc.CreateResponse {
	return &desc.CreateResponse{
		Id: userID,
	}
}

// ToGetResponseFromUserService converts user model to api response.
func ToGetResponseFromUserService(user *model.User) *desc.GetResponse {
	resp := &desc.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	if user.UpdatedAt != nil {
		resp.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return resp
}

// ToUserFromUpdateRequest converts api request to user model.
func ToUserFromUpdateRequest(request *desc.UpdateRequest) *model.User {
	user := &model.User{
		ID:   request.Id,
		Role: model.Role(request.Role),
	}

	if request.Name != nil {
		user.Name = request.Name.Value
	}

	if request.Email != nil {
		user.Email = request.Email.Value
	}

	return user
}
