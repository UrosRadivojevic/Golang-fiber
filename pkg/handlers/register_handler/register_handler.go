package register_handler

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/message"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories/user_repository"
	"github.com/urosradivojevic/health/pkg/requests/user_request.go"
	"github.com/urosradivojevic/health/pkg/services/hasher"
)

// ShowAccount godoc
//
//		@Summary		  Register user
//		@Description	Register user in database
//		@Tags			    users
//		@Accept			  json
//		@Produce		  json
//		@Success		 201
//	 @Failure      	 422   {object}    message.Msg "Validation failed"
//	 @Param request body user_request.UserRequest true "User"
//		@Router			 /register [post]
func Handler(repo user_repository.Interface, hash hasher.Interface, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request user_request.UserRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message.Msg{
				Message: err.Error(),
			})
		}

		err := validate.Struct(request)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message.Msg{
				Message: "Validation failed",
			})
		}
		hashPassword, err := hash.Hash(request.Password)
		if err != nil {
			return err
		}
		user := model.User{
			Username:  request.Username,
			Password:  string(hashPassword),
			Firstname: request.Firstname,
		}
		_, err = repo.Register(c.UserContext(), user)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}
