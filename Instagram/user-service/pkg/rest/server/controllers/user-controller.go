package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/Rohithknaidu/Instagram/user-service/pkg/rest/server/daos/clients/sqls"
	"github.com/Rohithknaidu/Instagram/user-service/pkg/rest/server/models"
	"github.com/Rohithknaidu/Instagram/user-service/pkg/rest/server/services"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() (*UserController, error) {
	userService, err := services.NewUserService()
	if err != nil {
		return nil, err
	}
	return &UserController{
		userService: userService,
	}, nil
}

func (userController *UserController) CreateUser(context echo.Context) error {
	// validate input
	var input models.User
	if err := context.Bind(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"error": err.Error()})
		return err
	}

	// trigger user creation
	userCreated, err := userController.userService.CreateUser(&input)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return err
	}

	context.JSON(http.StatusCreated, userCreated)
	return nil
}

func (userController *UserController) ListUsers(context echo.Context) error {
	// trigger all users fetching
	users, err := userController.userService.ListUsers()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return err
	}

	context.JSON(http.StatusOK, users)
	return nil
}

func (userController *UserController) FetchUser(context echo.Context) error {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return err
	}

	// trigger user fetching
	user, err := userController.userService.GetUser(id)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()})
			return err
		}
		context.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return err
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request().Context())
		currentSpan.SetAttributes(attribute.String("user.id", strconv.FormatInt(user.Id, 10)))
	}

	context.JSON(http.StatusOK, user)
	return nil
}

func (userController *UserController) UpdateUser(context echo.Context) error {
	// validate input
	var input models.User
	if err := context.Bind(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"error": err.Error()})
		return err
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return err
	}

	// trigger user update
	if _, err := userController.userService.UpdateUser(id, &input); err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()})
			return err
		}
		context.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return err
	}

	context.JSON(http.StatusNoContent, map[string]interface{}{})
	return nil
}

func (userController *UserController) DeleteUser(context echo.Context) error {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return err
	}

	// trigger user deletion
	if err := userController.userService.DeleteUser(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return err
	}

	context.JSON(http.StatusNoContent, map[string]interface{}{})
	return nil
}
