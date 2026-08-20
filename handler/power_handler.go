package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"power-price-monitor/service"

	"github.com/labstack/echo/v4"
)

type PowerHandler struct {
	powerService *service.PowerPriceService
}

type optimalHoursRequest struct {
	maxHours  int
	threshold float64
	selection service.OptimalHoursOptions
}

func NewPowerHandler(powerService *service.PowerPriceService) *PowerHandler {
	return &PowerHandler{
		powerService: powerService,
	}
}

// GetOptimalHours handles GET /optimal-hours
func (h *PowerHandler) GetOptimalHours(c echo.Context) error {
	request, err := parseOptimalHoursRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get optimal hours from service
	result, err := h.powerService.GetOptimalHoursWithOptions(c.Request().Context(), request.maxHours, request.threshold, request.selection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get optimal hours: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func parseOptimalHoursRequest(c echo.Context) (optimalHoursRequest, error) {
	maxHoursStr := c.QueryParam("max_hours")
	if maxHoursStr == "" {
		return optimalHoursRequest{}, fmt.Errorf("max_hours parameter is required")
	}

	thresholdStr := c.QueryParam("threshold")
	if thresholdStr == "" {
		return optimalHoursRequest{}, fmt.Errorf("threshold parameter is required")
	}

	maxHours, err := strconv.Atoi(maxHoursStr)
	if err != nil || maxHours < 0 || maxHours > 24 {
		return optimalHoursRequest{}, fmt.Errorf("max_hours must be an integer between 0 and 24 (0 disables scheduling)")
	}

	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil || threshold < 0 {
		return optimalHoursRequest{}, fmt.Errorf("threshold must be a positive number")
	}

	consecutive := false
	if value := c.QueryParam("consecutive"); value != "" {
		consecutive, err = strconv.ParseBool(value)
		if err != nil {
			return optimalHoursRequest{}, fmt.Errorf("consecutive must be a boolean")
		}
	}

	startHour, err := parseOptionalHour(c.QueryParam("start_hour"), 0)
	if err != nil {
		return optimalHoursRequest{}, fmt.Errorf("start_hour must be an integer between 0 and 24")
	}

	endHour, err := parseOptionalHour(c.QueryParam("end_hour"), 24)
	if err != nil {
		return optimalHoursRequest{}, fmt.Errorf("end_hour must be an integer between 0 and 24")
	}
	if startHour > endHour {
		return optimalHoursRequest{}, fmt.Errorf("start_hour must not be greater than end_hour")
	}

	return optimalHoursRequest{
		maxHours:  maxHours,
		threshold: threshold,
		selection: service.OptimalHoursOptions{
			Consecutive: consecutive,
			StartHour:   startHour,
			EndHour:     endHour,
		},
	}, nil
}

func parseOptionalHour(value string, defaultHour int) (int, error) {
	if value == "" {
		return defaultHour, nil
	}

	hour, err := strconv.Atoi(value)
	if err != nil || hour < 0 || hour > 24 {
		return 0, fmt.Errorf("invalid hour")
	}
	return hour, nil
}

// HealthCheck handles GET /health
func (h *PowerHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
