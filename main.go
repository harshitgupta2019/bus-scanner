package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Location represents a city or bus station
type Location struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	City    string  `json:"city"`
	State   string  `json:"state"`
	Country string  `json:"country"`
	Lat     float64 `json:"latitude"`
	Lng     float64 `json:"longitude"`
}

// BusOperator represents a bus company
type BusOperator struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Logo     string  `json:"logo"`
	Rating   float64 `json:"rating"`
	Platform string  `json:"platform"` // which booking platform
}

// BusType represents different types of buses
type BusType struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Seats       int      `json:"seats"`
	Amenities   []string `json:"amenities"`
	Description string   `json:"description"`
}

// Route represents a bus route with pricing
type Route struct {
	ID             string      `json:"id"`
	From           Location    `json:"from"`
	To             Location    `json:"to"`
	Operator       BusOperator `json:"operator"`
	BusType        BusType     `json:"bus_type"`
	DepartureTime  time.Time   `json:"departure_time"`
	ArrivalTime    time.Time   `json:"arrival_time"`
	Duration       string      `json:"duration"`
	Price          Price       `json:"price"`
	AvailableSeats int         `json:"available_seats"`
	BookingURL     string      `json:"booking_url"`
}

// Price represents pricing information
type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Platform string  `json:"platform"`
}

// SearchRequest represents a search query
type SearchRequest struct {
	FromCity   string    `json:"from_city"`
	ToCity     string    `json:"to_city"`
	Date       time.Time `json:"date"`
	Passengers int       `json:"passengers"`
}

// SearchResponse represents the aggregated search results
type SearchResponse struct {
	Status     string  `json:"status"`
	Message    string  `json:"message"`
	SearchID   string  `json:"search_id"`
	Routes     []Route `json:"routes"`
	TotalFound int     `json:"total_found"`
	SearchTime string  `json:"search_time"`
}

// BookingPlatform represents external booking platforms
type BookingPlatform struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key,omitempty"`
	IsActive bool   `json:"is_active"`
}

// Some sample data for testing
func GetSampleLocations() []Location {
	return []Location{
		{
			ID: "mumbai", Name: "Mumbai Central", City: "Mumbai",
			State: "Maharashtra", Country: "India", Lat: 19.0760, Lng: 72.8777,
		},
		{
			ID: "pune", Name: "Pune Station", City: "Pune",
			State: "Maharashtra", Country: "India", Lat: 18.5204, Lng: 73.8567,
		},
		{
			ID: "bangalore", Name: "Bangalore Majestic", City: "Bangalore",
			State: "Karnataka", Country: "India", Lat: 12.9716, Lng: 77.5946,
		},
		{
			ID: "delhi", Name: "Delhi ISBT", City: "Delhi",
			State: "Delhi", Country: "India", Lat: 28.7041, Lng: 77.1025,
		},
	}
}

func GetSampleOperators() []BusOperator {
	return []BusOperator{
		{ID: "redbus", Name: "RedBus", Logo: "redbus.png", Rating: 4.2, Platform: "redbus"},
		{ID: "makemytrip", Name: "MakeMyTrip", Logo: "mmt.png", Rating: 4.0, Platform: "makemytrip"},
		{ID: "goibibo", Name: "Goibibo", Logo: "goibibo.png", Rating: 3.9, Platform: "goibibo"},
		{ID: "abhibus", Name: "AbhiBus", Logo: "abhibus.png", Rating: 4.1, Platform: "abhibus"},
	}
}

func GetSampleBusTypes() []BusType {
	return []BusType{
		{
			ID: "ac_sleeper", Name: "AC Sleeper", Seats: 40,
			Amenities:   []string{"AC", "Sleeper", "Blanket", "Pillow"},
			Description: "Air conditioned sleeper bus with comfortable berths",
		},
		{
			ID: "non_ac_seater", Name: "Non-AC Seater", Seats: 50,
			Amenities:   []string{"Pushback Seats", "Charging Point"},
			Description: "Comfortable seater bus for day travel",
		},
		{
			ID: "volvo_ac", Name: "Volvo AC", Seats: 45,
			Amenities:   []string{"AC", "WiFi", "Entertainment", "USB Charging"},
			Description: "Premium Volvo bus with luxury amenities",
		},
	}
}

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// homeHandler handles the root endpoint
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Welcome to Bus Booking Aggregator API",
		Data:    map[string]string{"version": "1.0.0"},
	}
	sendJSON(w, http.StatusOK, response)
}

// healthHandler handles health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Service is healthy",
	}
	sendJSON(w, http.StatusOK, response)
}

func main() {
	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)

	// Server configuration
	port := ":8080"

	fmt.Printf("🚌 Bus Booking Aggregator starting on port %s\n", port)
	fmt.Printf("📍 http://localhost%s\n", port)

	// Start the server
	log.Fatal(http.ListenAndServe(port, mux))
}

type PlatformService interface {
	SearchRoutes(req SearchRequest) ([]Route, error)
	GetPlatformName() string
}

// RedBusService simulates RedBus API
type RedBusService struct {
	Name string
}

func (r *RedBusService) GetPlatformName() string {
	return "RedBus"
}

func (r *RedBusService) SearchRoutes(req SearchRequest) ([]Route, error) {
	// Simulate API delay
	time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)

	locations := GetSampleLocations()
	operators := GetSampleOperators()
	busTypes := GetSampleBusTypes()

	// Find from and to locations
	var fromLoc, toLoc Location
	for _, loc := range locations {
		if loc.City == req.FromCity {
			fromLoc = loc
		}
		if loc.City == req.ToCity {
			toLoc = loc
		}
	}

	// Generate mock routes
	routes := []Route{}
	basePrice := 500 + rand.Float64()*1000 // Random base price between 500-1500

	for i := 0; i < rand.Intn(3)+2; i++ { // 2-4 routes
		route := Route{
			ID:            fmt.Sprintf("redbus_%d", i+1),
			From:          fromLoc,
			To:            toLoc,
			Operator:      operators[0], // RedBus operator
			BusType:       busTypes[rand.Intn(len(busTypes))],
			DepartureTime: req.Date.Add(time.Hour * time.Duration(6+i*4)),   // 6AM, 10AM, 2PM, 6PM
			ArrivalTime:   req.Date.Add(time.Hour * time.Duration(6+i*4+8)), // +8 hours journey
			Duration:      "8h 0m",
			Price: Price{
				Amount:   basePrice + float64(i*100),
				Currency: "INR",
				Platform: "RedBus",
			},
			AvailableSeats: rand.Intn(20) + 5, // 5-25 seats
			BookingURL:     "https://redbus.in/book/route123",
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// MakeMyTripService simulates MakeMyTrip API
type MakeMyTripService struct {
	Name string
}

func (m *MakeMyTripService) GetPlatformName() string {
	return "MakeMyTrip"
}

func (m *MakeMyTripService) SearchRoutes(req SearchRequest) ([]Route, error) {
	time.Sleep(time.Duration(rand.Intn(600)+300) * time.Millisecond)

	locations := GetSampleLocations()
	operators := GetSampleOperators()
	busTypes := GetSampleBusTypes()

	var fromLoc, toLoc Location
	for _, loc := range locations {
		if loc.City == req.FromCity {
			fromLoc = loc
		}
		if loc.City == req.ToCity {
			toLoc = loc
		}
	}

	routes := []Route{}
	basePrice := 450 + rand.Float64()*1200

	for i := 0; i < rand.Intn(4)+1; i++ { // 1-4 routes
		route := Route{
			ID:            fmt.Sprintf("mmt_%d", i+1),
			From:          fromLoc,
			To:            toLoc,
			Operator:      operators[1], // MakeMyTrip operator
			BusType:       busTypes[rand.Intn(len(busTypes))],
			DepartureTime: req.Date.Add(time.Hour * time.Duration(7+i*3)),
			ArrivalTime:   req.Date.Add(time.Hour * time.Duration(7+i*3+9)),
			Duration:      "9h 0m",
			Price: Price{
				Amount:   basePrice + float64(i*150),
				Currency: "INR",
				Platform: "MakeMyTrip",
			},
			AvailableSeats: rand.Intn(15) + 3,
			BookingURL:     "https://makemytrip.com/bus/book/xyz",
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// GoibiboService simulates Goibibo API
type GoibiboService struct {
	Name string
}

func (g *GoibiboService) GetPlatformName() string {
	return "Goibibo"
}

func (g *GoibiboService) SearchRoutes(req SearchRequest) ([]Route, error) {
	time.Sleep(time.Duration(rand.Intn(400)+250) * time.Millisecond)

	locations := GetSampleLocations()
	operators := GetSampleOperators()
	busTypes := GetSampleBusTypes()

	var fromLoc, toLoc Location
	for _, loc := range locations {
		if loc.City == req.FromCity {
			fromLoc = loc
		}
		if loc.City == req.ToCity {
			toLoc = loc
		}
	}

	routes := []Route{}
	basePrice := 600 + rand.Float64()*900

	for i := 0; i < rand.Intn(3)+2; i++ { // 2-4 routes
		route := Route{
			ID:            fmt.Sprintf("goibibo_%d", i+1),
			From:          fromLoc,
			To:            toLoc,
			Operator:      operators[2], // Goibibo operator
			BusType:       busTypes[rand.Intn(len(busTypes))],
			DepartureTime: req.Date.Add(time.Hour * time.Duration(8+i*4)),
			ArrivalTime:   req.Date.Add(time.Hour * time.Duration(8+i*4+7)),
			Duration:      "7h 30m",
			Price: Price{
				Amount:   basePrice + float64(i*80),
				Currency: "INR",
				Platform: "Goibibo",
			},
			AvailableSeats: rand.Intn(25) + 8,
			BookingURL:     "https://goibibo.com/bus/booking/abc",
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// PlatformManager manages all booking platforms
type PlatformManager struct {
	platforms []PlatformService
}

func NewPlatformManager() *PlatformManager {
	return &PlatformManager{
		platforms: []PlatformService{
			&RedBusService{Name: "RedBus"},
			&MakeMyTripService{Name: "MakeMyTrip"},
			&GoibiboService{Name: "Goibibo"},
		},
	}
}

func (pm *PlatformManager) SearchAllPlatforms(req SearchRequest) ([]Route, error) {
	var allRoutes []Route

	// Channel to collect routes from all platforms
	routesChan := make(chan []Route, len(pm.platforms))
	errorsChan := make(chan error, len(pm.platforms))

	// Search all platforms concurrently
	for _, platform := range pm.platforms {
		go func(p PlatformService) {
			routes, err := p.SearchRoutes(req)
			if err != nil {
				errorsChan <- fmt.Errorf("%s error: %v", p.GetPlatformName(), err)
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
			fmt.Printf("Warning: %v\n", err)
			continue
		}

		if routes != nil {
			allRoutes = append(allRoutes, routes...)
		}
	}

	return allRoutes, nil
}
