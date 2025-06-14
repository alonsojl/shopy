package apigateway

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"shopy/internal/domain"
	"shopy/internal/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Service interface {
	LoginUser(ctx context.Context, email, password string) (string, error)
	AddUser(ctx context.Context, params domain.UserParams) (*models.User, error)
	DelUser(ctx context.Context, email string) error
}

type User struct {
	logger  *slog.Logger
	service Service
}

func NewUser(logger *slog.Logger, service Service) *User {
	return &User{
		logger:  logger,
		service: service,
	}
}

func (u *User) Router() APIGatewayFunc {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch event.HTTPMethod {
		case http.MethodPost:
			return u.HandleLoginUser(ctx, event)
		case http.MethodPut:
			return u.HandleAddUser(ctx, event)
		case http.MethodDelete:
			return u.HandleDelUser(ctx, event)
		}
		return events.APIGatewayProxyResponse{
			Body:       "method is not valid",
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
}

// @Summary 	Login user.
// @Description Login with user credentials.
// @Tags 		Users
// @Router 		/users [post]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param	    params body  UserAddRequest true "Credentials"
// @Success     200	{object} UserAuthorized "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (u *User) HandleLoginUser(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request UserCredentials

	if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
		u.logger.Error("invalid user body", "error", err)
		return Error(domain.ErrBodyRequest)
	}

	if err := request.Validate(); err != nil {
		u.logger.Error("invalid user params", "error", err)
		return Error(domain.ErrParams.Wrap(err))
	}

	token, err := u.service.LoginUser(ctx, request.Email, request.Password)
	if err != nil {
		u.logger.Error("error login user", "error", err)
		return Error(err)
	}

	var response = UserAuthorized{
		BaseResponse: NewBaseResponse(http.StatusOK),
		Token:        token,
	}

	return JSON(response, http.StatusOK)
}

// @Summary 	Add user.
// @Description Add new user.
// @Tags 		Users
// @Router 		/users [put]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param	    params body  UserAddRequest true "User"
// @Success     201	{object} UserAdded "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (u *User) HandleAddUser(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request UserAddRequest

	if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
		u.logger.Error("invalid user body", "error", err)
		return Error(domain.ErrBodyRequest)
	}

	if err := request.Validate(); err != nil {
		u.logger.Error("invalid user params", "error", err)
		return Error(domain.ErrParams.Wrap(err))
	}

	now := time.Now().UTC()
	user, err := u.service.AddUser(ctx, domain.UserParams{
		Email:     request.Email,
		Password:  request.Password,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		u.logger.Error("error adding user", "error", err)
		return Error(err)
	}

	var response = UserAdded{
		BaseResponse: NewBaseResponse(http.StatusCreated),
		User:         user,
	}

	return JSON(response, http.StatusCreated)
}

// @Summary 	Delete user.
// @Description Delete user profile.
// @Tags 		Users
// @Router 		/users/{email} [delete]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param       email path string true "Email"
// @Success     200	{object} UserDeleted "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (u *User) HandleDelUser(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := event.PathParameters["email"]
	if err := u.service.DelUser(ctx, email); err != nil {
		u.logger.Error("error deleting user", "error", err)
		return Error(err)
	}

	var response = UserDeleted{
		BaseResponse: NewBaseResponse(http.StatusOK),
		User:         "deleted",
	}

	return JSON(response, http.StatusOK)
}
