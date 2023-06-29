package handlers

import (
	"net/http"
	"remote_monitoring_and_controlling/internal/tasks"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetTasks(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		log.Error().Str("function", "GetTasks").Msg("Failed to get userID")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	tasks, err := h.tasks.GetTasks(ctx.Request().Context(), userID)
	if err != nil {
		log.Error().Str("function", "GetTasks").Err(err).Msg("Failed to get tasks")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrGetTasks))
	}

	return ctx.JSON(http.StatusOK, tasks)
}

func (h *Handler) GetVacations(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		log.Error().Str("function", "GetVacations").Msg("Failed to get userID")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	vacations, err := h.tasks.GetVacations(ctx.Request().Context(), userID)
	if err != nil {
		log.Error().Str("function", "GetVacations").Err(err).Msg("Failed to get vacations")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrGetVacations))
	}

	return ctx.JSON(http.StatusOK, vacations)
}

func (h *Handler) GetPayments(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		log.Error().Str("function", "GetPayments").Msg("Failed to get userID")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	payments, err := h.tasks.GetPayments(ctx.Request().Context(), userID)
	if err != nil {
		log.Error().Str("function", "GetPayments").Err(err).Msg("Failed to get payments")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrGetPayments))
	}

	return ctx.JSON(http.StatusOK, payments)
}

func (h *Handler) SaveTasks(ctx echo.Context) error {
	var req []tasks.Task
	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "SaveTasks").Err(err).Msg("Failed to bind tasks data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	if err := h.validator.Var(&req, "dive"); err != nil {
		log.Error().Str("function", "SaveTasks").Err(err).Msg("Failed to validate tasks data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	err := h.tasks.SaveTasks(ctx.Request().Context(), req)
	if err != nil {
		log.Error().Str("function", "SaveTasks").Err(err).Msg("Failed to save tasks data")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(checkErr(err, ErrSaveTasks)))
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) SaveVacations(ctx echo.Context) error {
	var req []tasks.Vacation
	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "SaveVacations").Err(err).Msg("Failed to bind vacations data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	if err := h.validator.Var(&req, "dive"); err != nil {
		log.Error().Str("function", "SaveVacations").Err(err).Msg("Failed to validate vacations data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	err := h.tasks.SaveVacations(ctx.Request().Context(), req)
	if err != nil {
		log.Error().Str("function", "SaveVacations").Err(err).Msg("Failed to save vacations data")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(checkErr(err, ErrSaveVacations)))
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) SavePayments(ctx echo.Context) error {
	var req []tasks.Payment
	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "SavePayments").Err(err).Msg("Failed to bind payments data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	if err := h.validator.Var(&req, "dive"); err != nil {
		log.Error().Str("function", "SavePayments").Err(err).Msg("Failed to validate payments data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	err := h.tasks.SavePayments(ctx.Request().Context(), req)
	if err != nil {
		log.Error().Str("function", "SavePayments").Err(err).Msg("Failed to save payments data")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(checkErr(err, ErrSavePayments)))
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) UpdateTaskStatus(ctx echo.Context) error {
	var req struct {
		TaskID string `json:"Uid"`
		Status bool   `json:"AgreeStatus"`
	}

	if err := ctx.Bind(&req); err != nil {
		log.Error().Str("function", "UpdateTaskStatus").Err(err).Msg("Failed to bind data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	if err := h.validator.Var(&req, "dive"); err != nil {
		log.Error().Str("function", "UpdateTaskStatus").Err(err).Msg("Failed to validate data")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	err := h.tasks.UpdateTaskStatus(ctx.Request().Context(), req.TaskID, req.Status)
	if err != nil {
		log.Error().Str("function", "UpdateTaskStatus").Err(err).Msg("Failed to update task status")
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(checkErr(err, ErrUpdateTaskStatus)))
	}

	return ctx.NoContent(http.StatusOK)
}
