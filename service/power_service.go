package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"power-price-monitor/config"
	"power-price-monitor/model"
)

const (
	SpainGeoID = 8741
)

type PowerPriceService struct {
	config     *config.Config
	cachedData *model.CachedPowerData
}

func NewPowerPriceService(cfg *config.Config) *PowerPriceService {
	return &PowerPriceService{
		config: cfg,
	}
}

// GetOptimalHours returns the optimal hours based on max_hours and threshold
func (s *PowerPriceService) GetOptimalHours(ctx context.Context, maxHours int, thresholdKWh float64) (*model.OptimalHoursResponse, error) {
	// Ensure we have current data
	if err := s.ensureCurrentData(ctx); err != nil {
		return nil, fmt.Errorf("failed to get current data: %w", err)
	}

	// Convert threshold from €/kWh to €/MWh
	thresholdMWh := thresholdKWh * 1000

	// Get timezone
	loc := s.getTimezone()

	// Get values for today in local timezone
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")
	var todayValues []model.Value

	for _, value := range s.cachedData.Values {
		if value.GeoID == SpainGeoID {
			valueTime, err := parseValueTime(value.DateTime)
			if err != nil {
				continue
			}
			valueTimeLocal := valueTime.In(loc)
			if valueTimeLocal.Format("2006-01-02") == today {
				todayValues = append(todayValues, value)
			}
		}
	}

	// Sort by price to find optimal hours
	sort.Slice(todayValues, func(i, j int) bool {
		return todayValues[i].Value < todayValues[j].Value
	})

	// Filter by threshold and max hours
	var optimalValues []model.Value

	// If maxHours is 0, return empty slice (disables scheduling)
	if maxHours > 0 {
		for _, value := range todayValues {
			if len(optimalValues) >= maxHours {
				break
			}
			if value.Value <= thresholdMWh {
				optimalValues = append(optimalValues, value)
			}
		}
	}

	// Sort optimal values by time for response
	sort.Slice(optimalValues, func(i, j int) bool {
		timeI, _ := parseValueTime(optimalValues[i].DateTime)
		timeJ, _ := parseValueTime(optimalValues[j].DateTime)
		return timeI.Before(timeJ)
	})

	// Find next start time (next optimal hour from now) in local timezone
	var nextStart *time.Time
	for _, value := range optimalValues {
		valueTime, err := parseValueTime(value.DateTime)
		if err != nil {
			continue
		}
		valueTimeLocal := valueTime.In(loc)
		if valueTimeLocal.After(now) {
			nextStart = &valueTimeLocal
			break
		}
	}

	// Get current price
	currentPrice := s.getCurrentPrice(todayValues, loc)

	// Build response with local timezone
	optimalHours := make([]string, 0)
	for _, value := range optimalValues {
		valueTime, err := parseValueTime(value.DateTime)
		if err != nil {
			continue
		}
		valueTimeLocal := valueTime.In(loc)
		optimalHours = append(optimalHours, valueTimeLocal.Format("15:04"))
	}

	response := &model.OptimalHoursResponse{
		OptimalHours:       optimalHours,
		TotalHoursSelected: len(optimalValues),
		CurrentPrice:       currentPrice / 1000, // Convert to €/kWh
		ThresholdUsed:      thresholdKWh,
		MaxHoursUsed:       maxHours,
	}

	if nextStart != nil {
		response.NextStart = nextStart.Format(time.RFC3339)
	}

	return response, nil
}

// ensureCurrentData fetches new data if cache is stale
func (s *PowerPriceService) ensureCurrentData(ctx context.Context) error {
	if s.cachedData != nil && s.isCacheValid() {
		return nil
	}

	data, err := s.fetchPowerPrices(ctx)
	if err != nil {
		return err
	}

	s.cachedData = &model.CachedPowerData{
		Values:     data.Indicator.Values,
		CachedAt:   time.Now(),
		ValidUntil: s.getNextDayStart(),
	}

	return nil
}

// isCacheValid checks if cached data is still valid (invalidate when day changes)
func (s *PowerPriceService) isCacheValid() bool {
	if s.cachedData == nil {
		return false
	}

	// Get timezone
	loc := s.getTimezone()

	now := time.Now().In(loc)
	cacheTime := s.cachedData.CachedAt.In(loc)

	// Cache is invalid if the day has changed
	return now.Format("2006-01-02") == cacheTime.Format("2006-01-02")
}

// getNextDayStart returns the start of the next day in local timezone
func (s *PowerPriceService) getNextDayStart() time.Time {
	loc := s.getTimezone()

	now := time.Now().In(loc)
	tomorrow := now.AddDate(0, 0, 1)
	return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, loc)
}

// fetchPowerPrices fetches data from the API
func (s *PowerPriceService) fetchPowerPrices(ctx context.Context) (*model.IndicatorResponse, error) {
	url := "https://api.esios.ree.es/indicators/1001"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json; application/vnd.esios-api-v2+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "api.esios.ree.es")
	req.Header.Set("x-api-key", s.config.Token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var response model.IndicatorResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &response, nil
}

// getCurrentPrice returns the current hour's price
func (s *PowerPriceService) getCurrentPrice(values []model.Value, loc *time.Location) float64 {
	now := time.Now().In(loc)
	currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)

	for _, value := range values {
		valueTime, err := parseValueTime(value.DateTime)
		if err != nil {
			continue
		}
		valueTimeLocal := valueTime.In(loc)
		valueHour := time.Date(valueTimeLocal.Year(), valueTimeLocal.Month(), valueTimeLocal.Day(), valueTimeLocal.Hour(), 0, 0, 0, loc)

		if valueHour.Equal(currentHour) {
			return value.Value
		}
	}

	return 0.0 // Default if current hour not found
}

func (s *PowerPriceService) getTimezone() *time.Location {
	loc, err := time.LoadLocation(s.config.Timezone)
	if err != nil {
		return time.UTC // fallback to UTC if timezone is invalid
	}
	return loc
}

func parseValueTime(dateTime string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000-07:00", dateTime)
}
