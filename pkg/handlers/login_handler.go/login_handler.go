package login_handler

import (
	"errors"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/requests/user_request.go"
	"github.com/urosradivojevic/health/pkg/services/login"
)

type message struct {
	Message string `json:"message"`
}

// ShowAccount godoc
//
//		@Summary		  Login user
//		@Description	Login user
//		@Tags			    users
//		@Accept			  json
//		@Produce		  json
//		@Success		 200
//	 @Failure      	 422   {object}    message "Validation failed"
//	 @Failure      	 400   {object}    message "Invalid login credentials"
//	 @Param request body user_request.UserRequest true "User"
//		@Router			 /login [post]
func Handler(loginService login.Interface, validator *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request user_request.UserRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message{
				Message: err.Error(),
			})
		}

		if err := validator.Struct(request); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message{
				Message: "Validation failed",
			})
		}
		user, err := loginService.Login(request.Username, request.Password)
		if err != nil {
			if errors.Is(err, login.ErrUsernameNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(message{
					Message: err.Error(),
				})
			}
			if errors.Is(err, login.ErrInvalidCredentials) {
				return c.Status(fiber.StatusBadRequest).JSON(message{
					Message: err.Error(),
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
