package controller

import (
	"github.com/Caknoooo/go-pagination"
	"github.com/gofiber/fiber/v3"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user/dto"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user/query"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user/service"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user/validation"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/constants"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/utils"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type (
	UserController interface {
		Me(ctx fiber.Ctx) error
		GetAllUser(ctx fiber.Ctx) error
		Update(ctx fiber.Ctx) error
		Delete(ctx fiber.Ctx) error
	}

	userController struct {
		userService    service.UserService
		userValidation *validation.UserValidation
		db             *gorm.DB
	}
)

func NewUserController(injector *do.Injector, us service.UserService) UserController {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	userValidation := validation.NewUserValidation()
	return &userController{
		userService:    us,
		userValidation: userValidation,
		db:             db,
	}
}

func (c *userController) GetAllUser(ctx fiber.Ctx) error {
	var filter = &query.UserFilter{}
	filter.BindPaginationForFiber(ctx)

	if err := ctx.Bind().Query(filter); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	users, total, err := pagination.PaginatedQueryWithIncludable[query.User](c.db, filter)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	paginationResponse := pagination.CalculatePagination(filter.Pagination, total)
	response := pagination.NewPaginatedResponse(fiber.StatusOK, dto.MESSAGE_SUCCESS_GET_LIST_USER, users, paginationResponse)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *userController) Me(ctx fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)

	result, err := c.userService.GetUserById(ctx.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *userController) Update(ctx fiber.Ctx) error {
	var req dto.UserUpdateRequest
	if err := ctx.Bind().Body(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	if err := c.userValidation.ValidateUserUpdateRequest(req); err != nil {
		res := utils.BuildResponseFailed("Validation failed", err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	userId := ctx.Locals("user_id").(string)
	result, err := c.userService.Update(ctx.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *userController) Delete(ctx fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)

	if err := c.userService.Delete(ctx.Context(), userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	return ctx.Status(fiber.StatusOK).JSON(res)
}
