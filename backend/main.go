package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	IsVisible bool   `json:"isVisible"`
}

type Service struct {
	ServiceID int    `json:"service_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Price     int    `json:"price"`
	Time      string `json:"time"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type ServiceArtistResponse struct {
	ServiceName string `json:"service_name"`
	UserName    string `json:"user_name"`
}

type ArtistsResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type BookingRequest struct {
	Artist    string `json:"artist"`
	Service   string `json:"service"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type BookingResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	BookingID   int    `json:"booking_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	UserName    string `json:"user_name,omitempty"`
}

type AvailableTimesResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Times   []string `json:"times,omitempty"`
}

type timeSlot struct {
	start time.Time
	end   time.Time
}

type AvailableTimesMonthlyResponse struct {
	Success bool              `json:"success"`
	Classes map[string]string `json:"classes"`
}

type Booking struct {
	BookingID    int    `json:"booking_id"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Date         string `json:"date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	ServiceID    int    `json:"service_id"`
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	ServiceName  string `json:"service_name"`
	ServiceColor string `json:"service_color"`
}

func main() {
	cors := corsMiddleware
	auth := AuthMiddleware

	http.Handle("/v0/api/login", cors(http.HandlerFunc(loginHandler)))
	http.Handle("/v0/api/artists", cors(http.HandlerFunc(artistsHandler)))
	http.Handle("/v0/api/add_booking", cors(http.HandlerFunc(createBookingHandler)))
	http.Handle("/v0/api/available_times", cors(http.HandlerFunc(availableTimesHandler)))
	http.Handle("/v0/api/availability", cors(http.HandlerFunc(availableTimesMonthlyHandler)))

	http.Handle("/v0/api/add_dash_booking", cors(auth(http.HandlerFunc(createDashBookingHandler))))
	http.Handle("/v0/api/get_bookings", cors(auth(http.HandlerFunc(getBookingsHandler))))
	http.Handle("/v0/api/bookings/", cors(auth(http.HandlerFunc(bookingHandler))))
	http.Handle("/v0/api/users", cors(auth(http.HandlerFunc(usersHandler))))
	http.Handle("/v0/api/user/", cors(auth(http.HandlerFunc(userHandler))))
	http.Handle("/v0/api/services", cors(auth(http.HandlerFunc(servicesHandler))))
	http.Handle("/v0/api/service/", cors(auth(http.HandlerFunc(serviceHandler))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func dbConn() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "123456"),
		getEnv("DB_NAME", "tuledb"),
	)
	return sql.Open("postgres", connStr)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

var jwtKey = []byte("your_secret_key")

func generateJWT(userID int) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"success":false,"message":"Missing token"}`, http.StatusUnauthorized)
			return
		}

		// Example: "Authorization <token>" — strip "Authorization " if you're using it
		tokenString := strings.TrimSpace(authHeader)

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"success":false,"message":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Store user ID if needed using claims.Subject (claims.Subject is userID in our case)

		next.ServeHTTP(w, r)
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"success":false,"message":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"message":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Retrieve the user's hashed password from the database
	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false,"message":"Database error"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user User
	var password string
	err = db.QueryRow("SELECT user_id, name, username, password FROM users WHERE username = $1", strings.TrimSpace(req.Username)).
		Scan(&user.UserID, &user.Name, &user.Username, &password)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Το όνομα χρήστη ή ο κωδικός πρόσβασης είναι λάθος"}`, http.StatusUnauthorized)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, password)
	if err != nil || !match {
		http.Error(w, `{"success":false,"message":"Το όνομα χρήστη ή ο κωδικός πρόσβασης είναι λάθος"}`, http.StatusUnauthorized)
		return
	}

	// If passwords match, generate JWT token
	token, err := generateJWT(user.UserID)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT user_id, name, username, is_visible FROM users ORDER BY user_id")
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if rows.Scan(&u.UserID, &u.Name, &u.Username, &u.IsVisible) == nil {
				users = append(users, u)
			}
		}
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var u struct {
			Name      string `json:"name"`
			Username  string `json:"username"`
			Password  string `json:"password"`
			IsVisible bool   `json:"isVisible"`
		}
		if json.NewDecoder(r.Body).Decode(&u) != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}

		hashedPassword, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
		if err != nil {
			http.Error(w, `{"success":false,"message":"Failed to hash password"}`, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (name, username, password, is_visible) VALUES ($1, $2, $3, $4)", u.Name, u.Username, hashedPassword, u.IsVisible)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"success":true}`))

	default:
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v0/api/user/")
	if id == "" {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodPut:
		var u struct {
			Name      string `json:"name"`
			Username  string `json:"username"`
			Password  string `json:"password"`
			IsVisible bool   `json:"isVisible"`
		}
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}

		var (
			query  string
			params []any
		)

		if u.Password != "" {
			hashedPassword, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
			if err != nil {
				http.Error(w, `{"success":false,"message":"Failed to hash password"}`, http.StatusInternalServerError)
				return
			}
			query = "UPDATE users SET name=$1, username=$2, password=$3, is_visible=$4 WHERE user_id=$5"
			params = []any{u.Name, u.Username, hashedPassword, u.IsVisible, id}
		} else {
			query = "UPDATE users SET name=$1, username=$2, is_visible=$3 WHERE user_id=$4"
			params = []any{u.Name, u.Username, u.IsVisible, id}
		}

		if _, err := db.Exec(query, params...); err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(`{"success":true}`))

	case http.MethodDelete:
		if _, err := db.Exec("DELETE FROM users WHERE user_id=$1", id); err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"success":true}`))

	default:
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
	}
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT service_id, name, color, price, time FROM services ORDER BY service_id")
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var services []Service
		for rows.Next() {
			var s Service
			if rows.Scan(&s.ServiceID, &s.Name, &s.Color, &s.Price, &s.Time) == nil {
				services = append(services, s)
			}
		}
		json.NewEncoder(w).Encode(services)

	case http.MethodPost:
		var s struct{ Name string; Color string; Price int; Time string }
		if json.NewDecoder(r.Body).Decode(&s) != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO services (name, color, price, time) VALUES ($1, $2, $3, $4)", s.Name, s.Color, s.Price, s.Time)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"success":true}`))

	default:
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
	}
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v0/api/service/")
	if id == "" {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodPut:
		var s Service
		if json.NewDecoder(r.Body).Decode(&s) != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE services SET name=$1, color=$2, price=$3, time=$4 WHERE service_id=$5", s.Name, s.Color, s.Price, s.Time, s.ServiceID)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"success":true}`))

	case http.MethodDelete:
		if _, err := db.Exec("DELETE FROM services WHERE service_id=$1", id); err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"success":true}`))

	default:
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
	}
}

func availableTimesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}

	date := r.URL.Query().Get("date")
	artistID, _ := strconv.Atoi(r.URL.Query().Get("artist_id"))

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT start_time, end_time FROM bookings WHERE user_id = $1 AND date = $2", artistID, date)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	bookedSlots := parseBookedSlots(rows)
	availableSlots := calculateAvailableSlots(bookedSlots)

	json.NewEncoder(w).Encode(AvailableTimesResponse{
		Success: true,
		Times:   availableSlots,
	})
}

func parseBookedSlots(rows *sql.Rows) []timeSlot {
	var slots []timeSlot
	for rows.Next() {
		var startStr, endStr string
		if rows.Scan(&startStr, &endStr) != nil {
			continue
		}

		start, _ := time.Parse("15:04", startStr)
		end, _ := time.Parse("15:04", endStr)
		slots = append(slots, timeSlot{start, end})
	}
	return slots
}

func calculateAvailableSlots(booked []timeSlot) []string {
	const slotDuration = 30 * time.Minute
	start, _ := time.Parse("15:04", "09:30")
	end, _ := time.Parse("15:04", "21:00")

	var available []string
	for current := start; current.Before(end); current = current.Add(slotDuration) {
		slotEnd := current.Add(slotDuration)
		if !isSlotBooked(current, slotEnd, booked) {
			available = append(available, fmt.Sprintf("%s - %s",
				current.Format("15:04"), slotEnd.Format("15:04")))
		}
	}
	return available
}

func isSlotBooked(start, end time.Time, booked []timeSlot) bool {
	for _, slot := range booked {
		if start.Before(slot.end) && end.After(slot.start) {
			return true
		}
	}
	return false
}

func availableTimesMonthlyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}

	artistID, _ := strconv.Atoi(r.URL.Query().Get("artist_id"))

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	startDate := time.Date(currentYear, currentMonth, now.Day(), 0, 0, 0, 0, currentLocation)
	endDate := time.Date(currentYear, currentMonth+2, 1, 0, 0, 0, 0, currentLocation).AddDate(0, 1, -1)

	rows, err := db.Query("SELECT date, start_time, end_time FROM bookings WHERE user_id = $1 AND date >= $2 AND date <= $3",
		artistID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	bookedSlotsByDate := parseBookedSlotsByDate(rows)

	result := make(map[string]string)
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		bookedSlots := bookedSlotsByDate[dateStr]
		availableSlots := calculateAvailableSlots(bookedSlots)
		slotCount := len(availableSlots)

		var availability string
		switch {
		case slotCount == 0:
			availability = "none"
		case slotCount <= 6:
			availability = "low"
		case slotCount <= 16:
			availability = "medium"
		default:
			availability = "high"
		}
		result[dateStr] = availability
	}

	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		bookedSlots := bookedSlotsByDate[dateStr]
		availableSlots := calculateAvailableSlots(bookedSlots)
		slotCount := len(availableSlots)

		var availability string
		switch {
		case slotCount == 0:
			availability = "none"
		case slotCount <= 6:
			availability = "low"
		case slotCount <= 16:
			availability = "medium"
		default:
			availability = "high"
		}
		result[dateStr] = availability
	}

	json.NewEncoder(w).Encode(AvailableTimesMonthlyResponse{
		Success: true,
		Classes: result,
	})
}

func parseBookedSlotsByDate(rows *sql.Rows) map[string][]timeSlot {
	slotsByDate := make(map[string][]timeSlot)

	startOfDay, _ := time.Parse("15:04", "09:00")
	endOfDay, _ := time.Parse("15:04", "21:00")

	for rows.Next() {
		var dateStr, startStr, endStr string
		if rows.Scan(&dateStr, &startStr, &endStr) != nil {
			continue
		}

		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			continue
		}

		start, _ := time.Parse("15:04", startStr)
		end, _ := time.Parse("15:04", endStr)

		if start.Before(startOfDay) || end.After(endOfDay) {
			continue
		}

		dateStr = parsedDate.Format("2006-01-02")

		slotsByDate[dateStr] = append(slotsByDate[dateStr], timeSlot{start, end})
	}
	return slotsByDate
}

func createBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	if req.Artist == "" || req.Service == "" || req.Date == "" || req.Time == "" ||
		req.FirstName == "" || req.LastName == "" || req.Phone == "" || req.Email == "" {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	timeParts := strings.Split(req.Time, " - ")
	if len(timeParts) != 2 {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(req.Artist)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	serviceID, err := strconv.Atoi(req.Service)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var bookingID int
	err = db.QueryRow(`
        INSERT INTO bookings (firstname, lastname, email, phone, date, 
            start_time, end_time, user_id, service_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING booking_id`,
		req.FirstName, req.LastName, req.Email, req.Phone,
		req.Date, timeParts[0], timeParts[1], artistID, serviceID,
	).Scan(&bookingID)

	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}

	response := BookingResponse{
		Success:   true,
		Message:   "Booking created successfully",
		BookingID: bookingID,
	}

	json.NewEncoder(w).Encode(response)
}

func createDashBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	if req.Artist == "" || req.Service == "" || req.Date == "" || req.Time == "" {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	timeParts := strings.Split(req.Time, " - ")
	if len(timeParts) != 2 {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(req.Artist)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	serviceID, err := strconv.Atoi(req.Service)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var bookingID int
	err = db.QueryRow(`
        INSERT INTO bookings (firstname, lastname, email, phone, date, 
            start_time, end_time, user_id, service_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING booking_id`,
		req.FirstName, req.LastName, req.Email, req.Phone,
		req.Date, timeParts[0], timeParts[1], artistID, serviceID,
	).Scan(&bookingID)

	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}

	response := BookingResponse{
		Success:   true,
		Message:   "Booking created successfully",
		BookingID: bookingID,
	}

	json.NewEncoder(w).Encode(response)
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT user_id, name FROM users WHERE is_visible = true")
	if err != nil {
		http.Error(w, `{"success":false,"error":"database query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var artists []ArtistsResponse

	for rows.Next() {
		var artist ArtistsResponse
		err := rows.Scan(&artist.UserID, &artist.Name)
		if err != nil {
			http.Error(w, `{"success":false,"error":"data scan failed"}`, http.StatusInternalServerError)
			return
		}
		artists = append(artists, artist)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, `{"success":false,"error":"rows iteration failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

func getBookingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT 
			b.booking_id, b.firstname, b.lastname, b.email, b.phone, 
			b.date, b.start_time, b.end_time, b.service_id, 
			u.user_id, u.name, s.name, s.color
		FROM bookings b
		JOIN users u ON u.user_id = b.user_id
		JOIN services s ON s.service_id = b.service_id
		ORDER BY b.date, b.start_time ASC
	`)

	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		err := rows.Scan(
			&b.BookingID, &b.FirstName, &b.LastName, &b.Email, &b.Phone,
			&b.Date, &b.StartTime, &b.EndTime, &b.ServiceID,
			&b.UserID, &b.UserName, &b.ServiceName, &b.ServiceColor,
		)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v0/api/bookings/")
	if id == "" {
		http.Error(w, `{"success":false, "error": "Booking ID missing"}`, http.StatusBadRequest)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false, "error": "Failed to connect to database"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodPut:
		var b BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, `{"success":false, "error": "Failed to parse request body"}`, http.StatusBadRequest)
			return
		}

		timeParts := strings.Split(b.Time, " - ")
		if len(timeParts) != 2 {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}

		query := `UPDATE bookings SET user_id=$1, service_id=$2, date=$3, start_time=$4,
					end_time=$5, firstname=$6, lastname=$7, phone=$8, email=$9 
		          WHERE booking_id=$10`
		_, err = db.Exec(query, b.Artist, b.Service, b.Date, timeParts[0], timeParts[1], b.FirstName,
			b.LastName, b.Phone, b.Email, id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success":false, "error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"success":true}`))

	case http.MethodDelete:
		if _, err := db.Exec("DELETE FROM bookings WHERE booking_id=$1", id); err != nil {
			http.Error(w, `{"success":false, "error": "Failed to delete booking"}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"success":true}`))

	default:
		http.Error(w, `{"success":false, "error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
