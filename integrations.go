package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// APIConfig holds configuration for different APIs
type APIConfig struct {
	BaseURL   string
	APIKey    string
	SecretKey string
	UserAgent string
	Timeout   time.Duration
	RateLimit time.Duration
}

// HTTPClient wraps http.Client with additional functionality
type HTTPClient struct {
	client   *http.Client
	config   APIConfig
	lastCall time.Time
}

func NewHTTPClient(config APIConfig) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

// RateLimit ensures we don't exceed API rate limits
func (h *HTTPClient) RateLimit() {
	if h.config.RateLimit > 0 {
		elapsed := time.Since(h.lastCall)
		if elapsed < h.config.RateLimit {
			time.Sleep(h.config.RateLimit - elapsed)
		}
	}
	h.lastCall = time.Now()
}

// MakeRequest performs HTTP request with proper headers and error handling
func (h *HTTPClient) MakeRequest(method, endpoint string, headers map[string]string, body interface{}) ([]byte, error) {
	h.RateLimit()

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, h.config.BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set default headers
	req.Header.Set("User-Agent", h.config.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set API key in header (common pattern)
	if h.config.APIKey != "" {
		req.Header.Set("X-API-Key", h.config.APIKey)
		req.Header.Set("Authorization", "Bearer "+h.config.APIKey)
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return responseBody, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// RealRedBusService integrates with actual RedBus API
type RealRedBusService struct {
	Name   string
	client *HTTPClient
}

func NewRealRedBusService(apiKey string) *RealRedBusService {
	config := APIConfig{
		BaseURL:   "https://api.redbus.com/v1", // Note: This is a placeholder - actual endpoint may differ
		APIKey:    apiKey,
		UserAgent: "BusAggregator/1.0",
		Timeout:   30 * time.Second,
		RateLimit: 1 * time.Second, // 1 request per second
	}

	return &RealRedBusService{
		Name:   "RedBus",
		client: NewHTTPClient(config),
	}
}

func (r *RealRedBusService) GetPlatformName() string {
	return r.Name
}

// RedBusSearchRequest represents the API request format
type RedBusSearchRequest struct {
	FromCityID    string `json:"fromCityId"`
	ToCityID      string `json:"toCityId"`
	DepartureDate string `json:"departureDate"`
	Passengers    int    `json:"passengers"`
}

// RedBusRoute represents the API response format
type RedBusRoute struct {
	ID             string   `json:"id"`
	OperatorName   string   `json:"operatorName"`
	BusType        string   `json:"busType"`
	DepartureTime  string   `json:"departureTime"`
	ArrivalTime    string   `json:"arrivalTime"`
	Duration       string   `json:"duration"`
	Fare           float64  `json:"fare"`
	AvailableSeats int      `json:"availableSeats"`
	Amenities      []string `json:"amenities"`
	BoardingPoints []string `json:"boardingPoints"`
	DroppingPoints []string `json:"droppingPoints"`
}

func (r *RealRedBusService) SearchRoutes(req SearchRequest) ([]Route, error) {
	// Convert our internal request format to RedBus API format
	redBusReq := RedBusSearchRequest{
		FromCityID:    r.getCityID(req.FromCity),
		ToCityID:      r.getCityID(req.ToCity),
		DepartureDate: req.Date.Format("2006-01-02"),
		Passengers:    req.Passengers,
	}

	endpoint := "/routes/search"
	responseBody, err := r.client.MakeRequest("POST", endpoint, nil, redBusReq)
	if err != nil {
		return nil, fmt.Errorf("RedBus API error: %v", err)
	}

	var apiResponse struct {
		Status string        `json:"status"`
		Data   []RedBusRoute `json:"data"`
	}

	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse RedBus response: %v", err)
	}

	// Convert RedBus routes to our internal format
	var routes []Route
	for _, rbRoute := range apiResponse.Data {
		route, err := r.convertRedBusRoute(rbRoute, req)
		if err != nil {
			continue // Skip invalid routes
		}
		routes = append(routes, route)
	}

	return routes, nil
}

func (r *RealRedBusService) getCityID(cityName string) string {
	// In real implementation, you'd have a mapping or call a cities API
	cityMap := map[string]string{
		"Mumbai":    "MUMBAI001",
		"Pune":      "PUNE001",
		"Bangalore": "BANGALORE001",
		"Delhi":     "DELHI001",
	}
	return cityMap[cityName]
}

func (r *RealRedBusService) convertRedBusRoute(rbRoute RedBusRoute, req SearchRequest) (Route, error) {
	locations := GetSampleLocations()
	var fromLoc, toLoc Location

	for _, loc := range locations {
		if loc.City == req.FromCity {
			fromLoc = loc
		}
		if loc.City == req.ToCity {
			toLoc = loc
		}
	}

	// Parse time strings
	departureTime, err := time.Parse("15:04", rbRoute.DepartureTime)
	if err != nil {
		return Route{}, fmt.Errorf("invalid departure time: %v", err)
	}

	arrivalTime, err := time.Parse("15:04", rbRoute.ArrivalTime)
	if err != nil {
		return Route{}, fmt.Errorf("invalid arrival time: %v", err)
	}

	// Combine with search date
	departureDateTime := time.Date(req.Date.Year(), req.Date.Month(), req.Date.Day(),
		departureTime.Hour(), departureTime.Minute(), 0, 0, req.Date.Location())
	arrivalDateTime := time.Date(req.Date.Year(), req.Date.Month(), req.Date.Day(),
		arrivalTime.Hour(), arrivalTime.Minute(), 0, 0, req.Date.Location())

	// Handle next day arrival
	if arrivalTime.Before(departureTime) {
		arrivalDateTime = arrivalDateTime.Add(24 * time.Hour)
	}

	return Route{
		ID:   rbRoute.ID,
		From: fromLoc,
		To:   toLoc,
		Operator: BusOperator{
			ID:       strings.ToLower(strings.ReplaceAll(rbRoute.OperatorName, " ", "_")),
			Name:     rbRoute.OperatorName,
			Platform: "RedBus",
			Rating:   4.0, // Default rating
		},
		BusType: BusType{
			ID:          strings.ToLower(strings.ReplaceAll(rbRoute.BusType, " ", "_")),
			Name:        rbRoute.BusType,
			Amenities:   rbRoute.Amenities,
			Description: rbRoute.BusType,
		},
		DepartureTime: departureDateTime,
		ArrivalTime:   arrivalDateTime,
		Duration:      rbRoute.Duration,
		Price: Price{
			Amount:   rbRoute.Fare,
			Currency: "INR",
			Platform: "RedBus",
		},
		AvailableSeats: rbRoute.AvailableSeats,
		BookingURL:     fmt.Sprintf("https://redbus.com/bus-tickets/%s", rbRoute.ID),
	}, nil
}

// RapidAPIBusService integrates with transportation APIs from RapidAPI
type RapidAPIBusService struct {
	Name   string
	client *HTTPClient
}

func NewRapidAPIBusService(apiKey string) *RapidAPIBusService {
	config := APIConfig{
		BaseURL:   "https://transport-api.p.rapidapi.com",
		APIKey:    apiKey,
		UserAgent: "BusAggregator/1.0",
		Timeout:   30 * time.Second,
		RateLimit: 2 * time.Second, // RapidAPI rate limit
	}

	return &RapidAPIBusService{
		Name:   "Transport API",
		client: NewHTTPClient(config),
	}
}

func (r *RapidAPIBusService) GetPlatformName() string {
	return r.Name
}

func (r *RapidAPIBusService) SearchRoutes(req SearchRequest) ([]Route, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("from", req.FromCity)
	params.Set("to", req.ToCity)
	params.Set("date", req.Date.Format("2006-01-02"))
	params.Set("passengers", strconv.Itoa(req.Passengers))

	endpoint := "/bus/search?" + params.Encode()

	headers := map[string]string{
		"X-RapidAPI-Host": "transport-api.p.rapidapi.com",
		"X-RapidAPI-Key":  r.client.config.APIKey,
	}

	responseBody, err := r.client.MakeRequest("GET", endpoint, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("RapidAPI error: %v", err)
	}

	// Parse response (format depends on the specific API)
	var apiResponse struct {
		Routes []struct {
			ID        string   `json:"id"`
			Operator  string   `json:"operator"`
			Departure string   `json:"departure"`
			Arrival   string   `json:"arrival"`
			Price     float64  `json:"price"`
			Duration  string   `json:"duration"`
			BusType   string   `json:"busType"`
			Seats     int      `json:"availableSeats"`
			Amenities []string `json:"amenities"`
		} `json:"routes"`
	}

	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	// Convert to our internal format
	var routes []Route
	locations := GetSampleLocations()
	var fromLoc, toLoc Location

	for _, loc := range locations {
		if loc.City == req.FromCity {
			fromLoc = loc
		}
		if loc.City == req.ToCity {
			toLoc = loc
		}
	}

	for _, apiRoute := range apiResponse.Routes {
		// Parse times
		departureTime, _ := time.Parse(time.RFC3339, apiRoute.Departure)
		arrivalTime, _ := time.Parse(time.RFC3339, apiRoute.Arrival)

		route := Route{
			ID:   apiRoute.ID,
			From: fromLoc,
			To:   toLoc,
			Operator: BusOperator{
				ID:       strings.ToLower(strings.ReplaceAll(apiRoute.Operator, " ", "_")),
				Name:     apiRoute.Operator,
				Platform: "Transport API",
				Rating:   3.8,
			},
			BusType: BusType{
				ID:        strings.ToLower(strings.ReplaceAll(apiRoute.BusType, " ", "_")),
				Name:      apiRoute.BusType,
				Amenities: apiRoute.Amenities,
			},
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			Duration:      apiRoute.Duration,
			Price: Price{
				Amount:   apiRoute.Price,
				Currency: "INR",
				Platform: "Transport API",
			},
			AvailableSeats: apiRoute.Seats,
			BookingURL:     fmt.Sprintf("https://example-booking.com/book/%s", apiRoute.ID),
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// Enhanced Platform Manager with real APIs
type RealPlatformManager struct {
	platforms []PlatformService
}

func NewRealPlatformManager(redBusAPIKey, rapidAPIKey string) *RealPlatformManager {
	platforms := []PlatformService{}

	// Add mock services for testing
	platforms = append(platforms, &RedBusService{Name: "RedBus Mock"})

	// Add real APIs if keys are provided
	if redBusAPIKey != "" {
		platforms = append(platforms, NewRealRedBusService(redBusAPIKey))
	}

	if rapidAPIKey != "" {
		platforms = append(platforms, NewRapidAPIBusService(rapidAPIKey))
	}

	return &RealPlatformManager{
		platforms: platforms,
	}
}

func (pm *RealPlatformManager) SearchAllPlatforms(req SearchRequest) ([]Route, error) {
	var allRoutes []Route

	// Channel to collect routes from all platforms
	routesChan := make(chan []Route, len(pm.platforms))
	errorsChan := make(chan error, len(pm.platforms))

	// Search all platforms concurrently
	for _, platform := range pm.platforms {
		go func(p PlatformService) {
			routes, err := p.SearchRoutes(req)
			if err != nil {
				fmt.Printf("Warning: %s error: %v\n", p.GetPlatformName(), err)
				errorsChan <- err
				routesChan <- nil
				return
			}
			routesChan <- routes
			errorsChan <- nil
		}(platform)
	}

	// Collect results
	for i := 0; i < len(pm.platforms); i++ {
		routes := <-routesChan
		err := <-errorsChan

		if err != nil {
			continue // Already logged above
		}

		if routes != nil {
			allRoutes = append(allRoutes, routes...)
		}
	}

	return allRoutes, nil
}

// Configuration structure for API keys
type Config struct {
	RedBusAPIKey string `json:"redbus_api_key"`
	RapidAPIKey  string `json:"rapidapi_key"`
	ServerPort   string `json:"server_port"`
}

// LoadConfig loads configuration from environment variables or config file
func LoadConfig() Config {
	// In production, you'd load from environment variables or config file
	return Config{
		RedBusAPIKey: "", // Set your RedBus API key here
		RapidAPIKey:  "", // Set your RapidAPI key here
		ServerPort:   ":8080",
	}
}
