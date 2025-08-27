package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
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

var platformManager *PlatformManager

// CORS middleware to handle cross-origin requests
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// searchHandler handles bus search requests
func searchHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		sendJSON(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: "Only POST method is allowed",
		})
		return
	}

	// Parse request body
	var searchReq SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid JSON request body",
		})
		return
	}

	// Validate required fields
	if searchReq.FromCity == "" || searchReq.ToCity == "" {
		sendJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "from_city and to_city are required",
		})
		return
	}

	// Set default values
	if searchReq.Passengers == 0 {
		searchReq.Passengers = 1
	}
	if searchReq.Date.IsZero() {
		searchReq.Date = time.Now().AddDate(0, 0, 1) // Tomorrow
	}

	start := time.Now()

	// Search all platforms
	routes, err := platformManager.SearchAllPlatforms(searchReq)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: fmt.Sprintf("Search failed: %v", err),
		})
		return
	}

	// Sort routes by price (ascending)
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Price.Amount < routes[j].Price.Amount
	})

	searchTime := time.Since(start)
	searchID := fmt.Sprintf("search_%d", time.Now().Unix())

	response := SearchResponse{
		Status:     "success",
		Message:    "Routes found successfully",
		SearchID:   searchID,
		Routes:     routes,
		TotalFound: len(routes),
		SearchTime: fmt.Sprintf("%.2fs", searchTime.Seconds()),
	}

	sendJSON(w, http.StatusOK, response)
}

// routesHandler returns available routes (GET endpoint for testing)
func routesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// Get query parameters
	fromCity := r.URL.Query().Get("from")
	toCity := r.URL.Query().Get("to")
	dateStr := r.URL.Query().Get("date")
	passengersStr := r.URL.Query().Get("passengers")

	if fromCity == "" || toCity == "" {
		sendJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "from and to parameters are required",
		})
		return
	}

	// Parse date
	var searchDate time.Time
	if dateStr != "" {
		var err error
		searchDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			sendJSON(w, http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid date format. Use YYYY-MM-DD",
			})
			return
		}
	} else {
		searchDate = time.Now().AddDate(0, 0, 1) // Tomorrow
	}

	// Parse passengers
	passengers := 1
	if passengersStr != "" {
		var err error
		passengers, err = strconv.Atoi(passengersStr)
		if err != nil || passengers < 1 {
			passengers = 1
		}
	}

	searchReq := SearchRequest{
		FromCity:   fromCity,
		ToCity:     toCity,
		Date:       searchDate,
		Passengers: passengers,
	}

	start := time.Now()
	routes, err := platformManager.SearchAllPlatforms(searchReq)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: fmt.Sprintf("Search failed: %v", err),
		})
		return
	}

	// Sort by price
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Price.Amount < routes[j].Price.Amount
	})

	searchTime := time.Since(start)

	response := SearchResponse{
		Status:     "success",
		Message:    "Routes found successfully",
		SearchID:   fmt.Sprintf("search_%d", time.Now().Unix()),
		Routes:     routes,
		TotalFound: len(routes),
		SearchTime: fmt.Sprintf("%.2fs", searchTime.Seconds()),
	}

	sendJSON(w, http.StatusOK, response)
}

// citiesHandler returns available cities
func citiesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	cities := []map[string]string{}
	for _, loc := range GetSampleLocations() {
		cities = append(cities, map[string]string{
			"id":    loc.ID,
			"name":  loc.City,
			"state": loc.State,
		})
	}

	sendJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Cities retrieved successfully",
		Data:    cities,
	})
}

// platformsHandler returns available booking platforms
func platformsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	platforms := []map[string]interface{}{
		{"name": "RedBus", "active": true, "avg_response_time": "300ms"},
		{"name": "MakeMyTrip", "active": true, "avg_response_time": "450ms"},
		{"name": "Goibibo", "active": true, "avg_response_time": "325ms"},
	}

	sendJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Platforms retrieved successfully",
		Data:    platforms,
	})
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Initialize platform manager
	platformManager = NewPlatformManager()

	// Create HTTP multiplexer
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/routes", routesHandler)
	mux.HandleFunc("/cities", citiesHandler)
	mux.HandleFunc("/platforms", platformsHandler)

	// Server configuration
	port := ":8080"

	fmt.Printf("🚌 Bus Booking Aggregator API\n")
	fmt.Printf("📍 Server: http://localhost%s\n", port)
	fmt.Printf("📋 Endpoints:\n")
	fmt.Printf("   GET  /              - API info\n")
	fmt.Printf("   GET  /health        - Health check\n")
	fmt.Printf("   GET  /cities        - Available cities\n")
	fmt.Printf("   GET  /platforms     - Booking platforms\n")
	fmt.Printf("   GET  /routes        - Search routes (query params)\n")
	fmt.Printf("   POST /search        - Search routes (JSON body)\n")
	fmt.Printf("\n🔍 Example GET request:\n")
	fmt.Printf("   http://localhost%s/routes?from=Mumbai&to=Pune&date=2025-08-21&passengers=2\n", port)
	fmt.Printf("\n📝 Example POST request body:\n")
	fmt.Printf(`   {
     "from_city": "Mumbai",
     "to_city": "Pune", 
     "date": "2025-08-21T00:00:00Z",
     "passengers": 2
   }`)
	fmt.Printf("\n🚀 Starting server...\n")

	// Start server
	log.Fatal(http.ListenAndServe(port, mux))
}
