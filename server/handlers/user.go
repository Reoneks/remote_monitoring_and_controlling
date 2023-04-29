package handlers

import (
	"net/http"

	"remote_monitoring_and_controlling/structs"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) Register(ctx echo.Context) error {
	var req structs.Register
	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "Register").Err(err).Msg("Failed to bind user data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Error().Str("function", "Register").Err(err).Msg("Failed to validate user data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	image, err := h.user.Register(ctx.Request().Context(), &req)
	if err != nil {
		log.Error().Str("function", "Register").Err(err).Msg("Failed to register user")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(checkErr(err, ErrRegister)))
	}

	if len(image) > 0 {
		return ctx.Blob(http.StatusOK, "image/png", image)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) Login(ctx echo.Context) error {
	var req structs.Login
	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "Login").Err(err).Msg("Failed to bind user data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Error().Str("function", "Login").Err(err).Msg("Failed to validate user data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	token, twoFAEnabled, err := h.user.Login(ctx.Request().Context(), &req)
	if err != nil {
		log.Error().Str("function", "Login").Err(err).Msg("Failed to login user")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(checkErr(err, ErrLogin)))
	}

	if twoFAEnabled {
		return ctx.String(http.StatusAccepted, token)
	}

	return ctx.String(http.StatusOK, token)
}

func (h *Handler) AddAlternativeNumber(ctx echo.Context) error {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "AddAlternativeNumber").Err(err).Msg("Failed to bind alternative phone number")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	userID, ok := ctx.Get("userID").(string)
	if !ok {
		log.Error().Str("function", "AddAlternativeNumber").Msg("Failed to get userID")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	err := h.user.AddAlternativeNumber(ctx.Request().Context(), userID, req.Phone)
	if err != nil {
		log.Error().Str("function", "AddAlternativeNumber").Err(err).Msg("Failed to add alternative phone number")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrAddPhone))
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) TwoFA(ctx echo.Context) error {
	var req structs.TwoFA
	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "TwoFA").Err(err).Msg("Failed to bind 2fa data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Error().Str("function", "TwoFA").Err(err).Msg("Failed to validate 2fa data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	token, err := h.user.OTPCheck(ctx.Request().Context(), &req)
	if err != nil {
		log.Error().Str("function", "TwoFA").Err(err).Msg("Failed to login user")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(checkErr(err, ErrLogin)))
	}

	return ctx.String(http.StatusOK, token)
}

func (h *Handler) EnableTwoFA(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		log.Error().Str("function", "EnableTwoFA").Msg("Failed to get userID")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	token, ok := ctx.Get("token").(string)
	if !ok {
		log.Error().Str("function", "EnableTwoFA").Msg("Failed to get token")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	image, err := h.user.EnableTwoFA(ctx.Request().Context(), userID, token)
	if err != nil {
		log.Error().Str("function", "EnableTwoFA").Err(err).Msg("Failed to login user")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrLogin))
	}

	return ctx.Blob(http.StatusOK, "image/png", image)
}

func (h *Handler) DisableTwoFA(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		log.Error().Str("function", "DisableTwoFA").Msg("Failed to get userID")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	token, ok := ctx.Get("token").(string)
	if !ok {
		log.Error().Str("function", "DisableTwoFA").Msg("Failed to get token")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	err := h.user.DisableTwoFA(ctx.Request().Context(), userID, token)
	if err != nil {
		log.Error().Str("function", "DisableTwoFA").Err(err).Msg("Failed to login user")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrLogin))
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) Logout(ctx echo.Context) error {
	token, ok := ctx.Get("token").(string)
	if !ok {
		log.Error().Str("function", "Logout").Msg("Failed to get token")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	h.user.Logout(ctx.Request().Context(), token)
	return ctx.NoContent(http.StatusOK)
}
