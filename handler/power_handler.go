package handler

import (
	"net/http"
	"strconv"

	"power-price-monitor/service"

	"github.com/labstack/echo/v4"
)

type PowerHandler struct {
	powerService *service.PowerPriceService
}

func NewPowerHandler(powerService *service.PowerPriceService) *PowerHandler {
	return &PowerHandler{
		powerService: powerService,
	}
}

// GetOptimalHours handles GET /optimal-hours
func (h *PowerHandler) GetOptimalHours(c echo.Context) error {
	// Parse query parameters
	maxHoursStr := c.QueryParam("max_hours")
	thresholdStr := c.QueryParam("threshold")

	// Validate required parameters
	if maxHoursStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "max_hours parameter is required",
		})
	}

	if thresholdStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "threshold parameter is required",
		})
	}

	// Parse max_hours
	maxHours, err := strconv.Atoi(maxHoursStr)
	if err != nil || maxHours < 0 || maxHours > 24 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "max_hours must be an integer between 0 and 24 (0 disables scheduling)",
		})
	}

	// Parse threshold
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil || threshold < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "threshold must be a positive number",
		})
	}

	// Get optimal hours from service
	result, err := h.powerService.GetOptimalHours(c.Request().Context(), maxHours, threshold)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get optimal hours: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

// HealthCheck handles GET /health
func (h *PowerHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
