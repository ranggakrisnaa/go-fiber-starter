package controller

import (
	"github.com/ranggakrisnaa/go-fiber-starter/modules/auth/dto"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/auth/service"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/auth/validation"
	userDto "github.com/ranggakrisnaa/go-fiber-starter/modules/user/dto"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/constants"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type (
	AuthController interface {
		Register(ctx fiber.Ctx) error
		Login(ctx fiber.Ctx) error
		RefreshToken(ctx fiber.Ctx) error
		Logout(ctx fiber.Ctx) error
		SendVerificationEmail(ctx fiber.Ctx) error
		VerifyEmail(ctx fiber.Ctx) error
		SendPasswordReset(ctx fiber.Ctx) error
		ResetPassword(ctx fiber.Ctx) error
	}

	authController struct {
		authService    service.AuthService
		authValidation *validation.AuthValidation
		db             *gorm.DB
	}
)

func NewAuthController(injector *do.Injector, as service.AuthService) AuthController {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	authValidation := validation.NewAuthValidation()
	return &authController{
		authService:    as,
		authValidation: authValidation,
		db:             db,
	}
}

func (c *authController) Register(ctx fiber.Ctx) error {
	var req userDto.UserCreateRequest
	if err := ctx.Bind().Body(&req); err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Validate request
	if err := c.authValidation.ValidateRegisterRequest(req); err != nil {
		res := utils.BuildResponseFailed("Validation failed", err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	result, err := c.authService.Register(ctx.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(userDto.MESSAGE_SUCCESS_REGISTER_USER, result)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) Login(ctx fiber.Ctx) error {
	var req userDto.UserLoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		response := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validate request
	if err := c.authValidation.ValidateLoginRequest(req); err != nil {
		res := utils.BuildResponseFailed("Validation failed", err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	result, err := c.authService.Login(ctx.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(userDto.MESSAGE_SUCCESS_LOGIN, result)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) RefreshToken(ctx fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := ctx.Bind().Body(&req); err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	result, err := c.authService.RefreshToken(ctx.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil)
		return ctx.Status(fiber.StatusUnauthorized).JSON(res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REFRESH_TOKEN, result)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) Logout(ctx fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)

	err := c.authService.Logout(ctx.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGOUT, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGOUT, nil)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) SendVerificationEmail(ctx fiber.Ctx) error {
	var req userDto.SendVerificationEmailRequest
	if err := ctx.Bind().Body(&req); err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	err := c.authService.SendVerificationEmail(ctx.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(userDto.MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS, nil)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) VerifyEmail(ctx fiber.Ctx) error {
	var req userDto.VerifyEmailRequest
	if err := ctx.Bind().Body(&req); err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	result, err := c.authService.VerifyEmail(ctx.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_VERIFY_EMAIL, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(userDto.MESSAGE_SUCCESS_VERIFY_EMAIL, result)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) SendPasswordReset(ctx fiber.Ctx) error {
	var req dto.SendPasswordResetRequest
	if err := ctx.Bind().Body(&req); err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	err := c.authService.SendPasswordReset(ctx.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SEND_PASSWORD_RESET, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_SEND_PASSWORD_RESET, nil)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) ResetPassword(ctx fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := ctx.Bind().Body(&req); err != nil {
		res := utils.BuildResponseFailed(userDto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	err := c.authService.ResetPassword(ctx.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_RESET_PASSWORD, err.Error(), nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_RESET_PASSWORD, nil)
	return ctx.Status(fiber.StatusOK).JSON(res)
}
