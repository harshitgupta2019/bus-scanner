package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Go Concept #1: Structs - Similar to classes in other languages
// Let's define our basic data structures
type Bus struct {
	ID            string    `json:"id"` // `json:"id"` is a struct tag for JSON serialization
	Operator      string    `json:"operator"`
	FromCity      string    `json:"from_city"`
	ToCity        string    `json:"to_city"`
	DepartureTime string    `json:"departure_time"`
	ArrivalTime   string    `json:"arrival_time"`
	Price         float64   `json:"price"`
	BusType       string    `json:"bus_type"`
	SeatsLeft     int       `json:"seats_left"`
	CreatedAt     time.Time `json:"created_at"`
}

// Go Concept #2: Structs for request/response
type SearchRequest struct {
	FromCity   string `json:"from_city" binding:"required"` // binding:"required" is Gin validation
	ToCity     string `json:"to_city" binding:"required"`
	Date       string `json:"date" binding:"required"`
	Passengers int    `json:"passengers" binding:"min=1"`
}

type SearchResponse struct {
	Buses     []Bus     `json:"buses"`
	Total     int       `json:"total"`
	Timestamp time.Time `json:"timestamp"`
	SearchID  string    `json:"search_id"`
}

// Go Concept #3: Methods on structs (like methods in a class)
// Let's add some methods to our Bus struct
func (b *Bus) IsExpensive() bool {
	return b.Price > 500.0
}

func (b *Bus) GetDuration() string {
	// Parse time strings and calculate duration
	dept, _ := time.Parse("15:04", b.DepartureTime)
	arr, _ := time.Parse("15:04", b.ArrivalTime)

	// Handle overnight journeys
	if arr.Before(dept) {
		arr = arr.Add(24 * time.Hour)
	}

	duration := arr.Sub(dept)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func (b *Bus) GetFormattedPrice() string {
	return fmt.Sprintf("₹%.2f", b.Price)
}

// Go Concept #4: Functions - Basic functions in Go
func generateSearchID() string {
	return fmt.Sprintf("search_%d", time.Now().UnixNano())
}

func createSampleBuses() []Bus {
	// Go Concept #5: Slices - Dynamic arrays in Go
	buses := []Bus{
		{
			ID:            "bus_001",
			Operator:      "RedBus Express",
			FromCity:      "Mumbai",
			ToCity:        "Pune",
			DepartureTime: "06:00",
			ArrivalTime:   "09:30",
			Price:         450.0,
			BusType:       "AC Sleeper",
			SeatsLeft:     25,
			CreatedAt:     time.Now(),
		},
		{
			ID:            "bus_002",
			Operator:      "VRL Travels",
			FromCity:      "Mumbai",
			ToCity:        "Pune",
			DepartureTime: "10:30",
			ArrivalTime:   "14:15",
			Price:         520.0,
			BusType:       "Volvo AC",
			SeatsLeft:     15,
			CreatedAt:     time.Now(),
		},
		{
			ID:            "bus_003",
			Operator:      "Orange Travels",
			FromCity:      "Mumbai",
			ToCity:        "Pune",
			DepartureTime: "18:45",
			ArrivalTime:   "22:15",
			Price:         395.0,
			BusType:       "Non-AC Sleeper",
			SeatsLeft:     30,
			CreatedAt:     time.Now(),
		},
	}

	return buses
}

// Go Concept #6: Function that returns multiple values (common in Go)
func filterBusesByPrice(buses []Bus, maxPrice float64) ([]Bus, int) {
	var filtered []Bus // Go Concept #7: var declaration with zero value

	// Go Concept #8: range - for loop over slices
	for _, bus := range buses { // _ discards the index, bus is the value
		if bus.Price <= maxPrice {
			filtered = append(filtered, bus) // append() adds to slice
		}
	}

	return filtered, len(filtered)
}

// Go Concept #9: Error handling - Go's approach to errors
func validateSearchRequest(req SearchRequest) error {
	if req.FromCity == req.ToCity {
		return fmt.Errorf("from_city and to_city cannot be the same")
	}

	// Parse and validate date
	_, err := time.Parse("2006-01-02", req.Date) // Go's reference time format
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	return nil // nil means no error in Go
}

// Go Concept #10: HTTP handlers with Gin framework
func searchBusesHandler(c *gin.Context) {
	var req SearchRequest

	// Go Concept #11: Error handling with if statement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate request
	if err := validateSearchRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create sample buses (in real app, this would query database/APIs)
	allBuses := createSampleBuses()

	// Filter buses by route (simple string matching for now)
	var matchingBuses []Bus
	for _, bus := range allBuses {
		if bus.FromCity == req.FromCity && bus.ToCity == req.ToCity {
			matchingBuses = append(matchingBuses, bus)
		}
	}

	// Create response
	response := SearchResponse{
		Buses:     matchingBuses,
		Total:     len(matchingBuses),
		Timestamp: time.Now(),
		SearchID:  generateSearchID(),
	}

	// Go Concept #12: JSON response
	c.JSON(http.StatusOK, response)
}

// Health check handler
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

// Go Concept #13: Maps - Key-value pairs
func busDetailsHandler(c *gin.Context) {
	busID := c.Param("id") // Get URL parameter

	// Sample bus details (in real app, query from database)
	busDetails := map[string]interface{}{
		"id":            busID,
		"operator":      "RedBus Express",
		"amenities":     []string{"WiFi", "Charging Port", "Water Bottle", "Blanket"},
		"pickup_points": []string{"Dadar", "Kurla", "Thane"},
		"drop_points":   []string{"Pune Station", "Katraj", "Hadapsar"},
		"cancellation_policy": map[string]interface{}{
			"free_cancellation_hours": 24,
			"cancellation_fee":        50.0,
		},
	}

	c.JSON(http.StatusOK, busDetails)
}

// Go Concept #14: main function - Entry point
func main() {
	// Go Concept #15: Package-level logging
	log.Println("🚌 Starting BusScanner API Server...")

	// Create Gin router
	r := gin.Default()

	// Go Concept #16: Middleware - Functions that run before handlers
	r.Use(func(c *gin.Context) {
		// Simple logging middleware
		start := time.Now()
		c.Next() // Continue to next handler
		latency := time.Since(start)

		log.Printf(
			"%s %s %d %v",
			c.Request.Method,
			c.Request.RequestURI,
			c.Writer.Status(),
			latency,
		)
	})

	// Define routes
	r.GET("/health", healthHandler)
	r.POST("/api/search", searchBusesHandler)
	r.GET("/api/bus/:id", busDetailsHandler)

	// Go Concept #17: Goroutine for background task (we'll explore more later)
	go func() {
		log.Println("Background task: Updating bus data every 5 minutes...")
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop() // defer runs when function exits

		for {
			select {
			case <-ticker.C:
				log.Println("Updating bus data...")
				// In real app: update prices, availability, etc.
			}
		}
	}()

	// Start server
	port := ":8080"
	log.Printf("✅ Server starting on http://localhost%s", port)
	log.Println("📍 Try these endpoints:")
	log.Println("   GET  http://localhost:8080/health")
	log.Println("   POST http://localhost:8080/api/search")
	log.Println("   GET  http://localhost:8080/api/bus/bus_001")

	// Go Concept #18: Error handling for server start
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
