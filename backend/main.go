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
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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
	Email     string `json:"email"`
	IsVisible bool   `json:"isVisible"`
}

type Service struct {
	ServiceID int    `json:"service_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
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
	Bath      string `json:"bath"`
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

type timeRange struct {
	startStr string
	endStr   string
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
	Bath         bool   `json:"bath"`
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	ServiceName  string `json:"service_name"`
	ServiceColor string `json:"service_color"`
}

type DisableDayRequest struct {
	Status bool   `json:"status"`
	Date   string `json:"date"`
	UserID int    `json:"userId"`
}

type UserWithWorkingHours struct {
	UserID       int                            `json:"userId"`
	Name         string                         `json:"name"`
	Username     string                         `json:"username"`
	Email        string                         `json:"email"`
	IsVisible    bool                           `json:"isVisible"`
	WorkingHours map[string][]map[string]string `json:"workingHours"`
}

type Client struct {
	ClientID  int    `json:"client_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func main() {
	cors := corsMiddleware
	auth := AuthMiddleware

	http.Handle("/v0/api/login", cors(http.HandlerFunc(loginHandler)))
	http.Handle("/v0/api/artists", cors(http.HandlerFunc(artistsHandler)))
	http.Handle("/v0/api/add_booking", cors(http.HandlerFunc(createBookingHandler)))
	http.Handle("/v0/api/available_times", cors(http.HandlerFunc(availableTimesHandler)))
	http.Handle("/v0/api/availability", cors(http.HandlerFunc(availableTimesMonthlyHandler)))
	http.Handle("/v0/api/forgot_password", cors(http.HandlerFunc(forgotPasswordHandler)))

	http.Handle("/v0/api/add_dash_booking", cors(auth(http.HandlerFunc(createDashBookingHandler))))
	http.Handle("/v0/api/disable_day", cors(auth(http.HandlerFunc(disableDayHandler))))
	http.Handle("/v0/api/get_client_data", cors(auth(http.HandlerFunc(clientDataHandler))))
	http.Handle("/v0/api/dash_artists", cors(auth(http.HandlerFunc(dashArtistsHandler))))
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
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST"),
		getEnv("DB_PORT"),
		getEnv("DB_USER"),
		getEnv("DB_PASSWORD"),
		getEnv("DB_NAME"),
	)
	return sql.Open("postgres", connStr)
}

func getEnv(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var jwtKey = []byte(getEnv("JWT_SECRET"))

func generateJWT(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * 30 * 24 * time.Hour)),
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

		claims := &jwt.RegisteredClaims{}
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
	err = db.QueryRow("SELECT user_id, name, username, email, password FROM users WHERE username = $1 OR email = $1", strings.TrimSpace(req.Username)).
		Scan(&user.UserID, &user.Name, &user.Username, &user.Email, &password)
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
		rows, err := db.Query(`
	SELECT u.user_id, u.name, u.username, u.email, u.is_visible, 
	       w.day_of_week, w.start_time, w.end_time
	FROM users u
	LEFT JOIN working_hours w ON w.user_id = u.user_id
	ORDER BY u.user_id DESC`)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		usersMap := make(map[int]*UserWithWorkingHours)

		for rows.Next() {
			var (
				id              int
				name            string
				username        string
				email           string
				isVisible       bool
				day, start, end sql.NullString
			)

			err := rows.Scan(&id, &name, &username, &email, &isVisible, &day, &start, &end)
			if err != nil {
				http.Error(w, `{"success":false}`, http.StatusInternalServerError)
				return
			}

			if _, exists := usersMap[id]; !exists {
				usersMap[id] = &UserWithWorkingHours{
					UserID:       id,
					Name:         name,
					Username:     username,
					Email:        email,
					IsVisible:    isVisible,
					WorkingHours: make(map[string][]map[string]string),
				}
			}

			if day.Valid && start.Valid && end.Valid {
				usersMap[id].WorkingHours[day.String] = append(usersMap[id].WorkingHours[day.String], map[string]string{
					"start": start.String,
					"end":   end.String,
				})
			}
		}

		var users []UserWithWorkingHours
		for _, u := range usersMap {
			users = append(users, *u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var u struct {
			Name         string                         `json:"name"`
			Username     string                         `json:"username"`
			Email        string                         `json:"email"`
			Password     string                         `json:"password"`
			IsVisible    bool                           `json:"isVisible"`
			WorkingHours map[string][]map[string]string `json:"workingHours"`
		}

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}

		hashedPassword, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
		if err != nil {
			http.Error(w, `{"success":false,"message":"Failed to hash password"}`, http.StatusInternalServerError)
			return
		}

		// Insert the user and get user_id
		var userID int
		err = db.QueryRow(
			"INSERT INTO users (name, username, email, password, is_visible) VALUES ($1, $2, $3, $4, $5) RETURNING user_id",
			u.Name, u.Username, u.Email, hashedPassword, u.IsVisible,
		).Scan(&userID)

		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}

		// Insert working hours
		stmt, err := db.Prepare("INSERT INTO working_hours (user_id, day_of_week, start_time, end_time) VALUES ($1, $2, $3, $4)")
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		for day, slots := range u.WorkingHours {
			for _, slot := range slots {
				start, end := slot["start"], slot["end"]
				if start != "" && end != "" {
					_, err := stmt.Exec(userID, day, start, end)
					if err != nil {
						http.Error(w, `{"success":false}`, http.StatusInternalServerError)
						return
					}
				}
			}
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
			Name         string                         `json:"name"`
			Username     string                         `json:"username"`
			Email        string                         `json:"email"`
			Password     string                         `json:"password"`
			IsVisible    bool                           `json:"isVisible"`
			WorkingHours map[string][]map[string]string `json:"workingHours"`
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
			query = "UPDATE users SET name=$1, username=$2, email=$3, password=$4, is_visible=$5 WHERE user_id=$6"
			params = []any{u.Name, u.Username, u.Email, hashedPassword, u.IsVisible, id}
		} else {
			query = "UPDATE users SET name=$1, username=$2, email=$3, is_visible=$4 WHERE user_id=$5"
			params = []any{u.Name, u.Username, u.Email, u.IsVisible, id}
		}

		if _, err := db.Exec(query, params...); err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}

		// Remove existing working hours
		if _, err := db.Exec("DELETE FROM working_hours WHERE user_id = $1", id); err != nil {
			http.Error(w, `{"success":false,"message":"Failed to clear working hours"}`, http.StatusInternalServerError)
			return
		}

		// Insert updated working hours
		stmt, err := db.Prepare("INSERT INTO working_hours (user_id, day_of_week, start_time, end_time) VALUES ($1, $2, $3, $4)")
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		for day, slots := range u.WorkingHours {
			for _, slot := range slots {
				start, end := slot["start"], slot["end"]
				if start != "" && end != "" {
					if _, err := stmt.Exec(id, day, start, end); err != nil {
						http.Error(w, `{"success":false}`, http.StatusInternalServerError)
						return
					}
				}
			}
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
		rows, err := db.Query("SELECT service_id, name, color, time FROM services ORDER BY service_id")
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var services []Service
		for rows.Next() {
			var s Service
			if rows.Scan(&s.ServiceID, &s.Name, &s.Color, &s.Time) == nil {
				services = append(services, s)
			}
		}
		json.NewEncoder(w).Encode(services)

	case http.MethodPost:
		var s struct {
			Name  string
			Color string
			Time  string
		}
		if json.NewDecoder(r.Body).Decode(&s) != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO services (name, color, time) VALUES ($1, $2, $3)", s.Name, s.Color, s.Time)
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

		_, err = db.Exec("UPDATE services SET name=$1, color=$2, time=$3 WHERE service_id=$4", s.Name, s.Color, s.Time, s.ServiceID)
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

func getWorkingHours(db *sql.DB, artistID int, weekday string) ([]timeRange, error) {
	rows, err := db.Query(`
		SELECT start_time, end_time 
		FROM working_hours 
		WHERE user_id = $1 AND LOWER(day_of_week) = LOWER($2)
	`, artistID, weekday)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hours []timeRange
	for rows.Next() {
		var start, end string
		if err := rows.Scan(&start, &end); err != nil {
			continue
		}
		hours = append(hours, timeRange{startStr: start, endStr: end})
	}
	return hours, nil
}

func availableTimesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")
	artistID, _ := strconv.Atoi(r.URL.Query().Get("artist_id"))

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT date, start_time, end_time FROM bookings WHERE user_id = $1 AND date = $2", artistID, dateStr)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	location, _ := time.LoadLocation("Europe/Athens")
	today := time.Now().In(location).Format("2006-01-02")
	isToday := dateStr == today

	bookedSlots := parseBookedSlots(rows)

	date, _ := time.Parse("2006-01-02", dateStr)

	dayOfWeek := date.Weekday().String()
	workingHours, err := getWorkingHours(db, artistID, dayOfWeek)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}

	availableSlots := calculateAvailableSlots(bookedSlots, isToday, workingHours, date)

	json.NewEncoder(w).Encode(AvailableTimesResponse{
		Success: true,
		Times:   availableSlots,
	})
}

func parseBookedSlots(rows *sql.Rows) []timeSlot {
	var slots []timeSlot
	for rows.Next() {
		var date time.Time
		var startStr, endStr string
		if err := rows.Scan(&date, &startStr, &endStr); err != nil {
			continue
		}

		loc, _ := time.LoadLocation("Europe/Athens")
		// Combine actual date object with start/end time strings
		startTime, errStart := time.ParseInLocation("15:04", startStr, loc)
		endTime, errEnd := time.ParseInLocation("15:04", endStr, loc)

		if errStart != nil || errEnd != nil {
			continue
		}

		start := time.Date(date.Year(), date.Month(), date.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0, loc)
		end := time.Date(date.Year(), date.Month(), date.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0, loc)

		slots = append(slots, timeSlot{start, end})
	}
	return slots
}

func calculateAvailableSlots(booked []timeSlot, isToday bool, workingHours []timeRange, date time.Time) []string {
	const slotDuration = 30 * time.Minute
	loc, _ := time.LoadLocation("Europe/Athens")

	var available []string

	for _, tr := range workingHours {
		startTime, err1 := time.ParseInLocation("15:04", tr.startStr, loc)
		endTime, err2 := time.ParseInLocation("15:04", tr.endStr, loc)
		if err1 != nil || err2 != nil {
			continue
		}

		start := time.Date(date.Year(), date.Month(), date.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0, loc)
		end := time.Date(date.Year(), date.Month(), date.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0, loc)

		current := start
		for current.Add(slotDuration).Before(end) || current.Add(slotDuration).Equal(end) {
			slotEnd := current.Add(slotDuration)

			if !isToday || (isToday && current.After(time.Now().In(loc))) {
				if !isSlotBooked(current, slotEnd, booked) {
					available = append(available, fmt.Sprintf("%s - %s",
						current.Format("15:04"), slotEnd.Format("15:04")))
				}
			}

			current = current.Add(slotDuration)
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

	location, _ := time.LoadLocation("Europe/Athens")
	today := time.Now().In(location).Format("2006-01-02")

	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		bookedSlots := bookedSlotsByDate[dateStr]

		isToday := dateStr == today

		workingHours, err := getWorkingHours(db, artistID, d.Weekday().String())
		
		if err != nil {
			continue
		}
		availableSlots := calculateAvailableSlots(bookedSlots, isToday, workingHours, d)

		slotCount := len(availableSlots)

		var availability string
		switch {
		case slotCount == 0:
			availability = "none"
		case slotCount <= 6:
			availability = "low"
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
    loc, _ := time.LoadLocation("Europe/Athens")

    for rows.Next() {
        var date time.Time
        var startStr, endStr string
        if err := rows.Scan(&date, &startStr, &endStr); err != nil {
            continue
        }

        // Parse just the clock times in the correct location
        tStart, err1 := time.ParseInLocation("15:04", startStr, loc)
        tEnd,   err2 := time.ParseInLocation("15:04", endStr,   loc)
        if err1 != nil || err2 != nil {
            continue
        }

        // Combine into full‐date time.Time values
        start := time.Date(date.Year(), date.Month(), date.Day(),
            tStart.Hour(), tStart.Minute(), 0, 0, loc)
        end := time.Date(date.Year(), date.Month(), date.Day(),
            tEnd.Hour(),   tEnd.Minute(),   0, 0, loc)

        key := date.Format("2006-01-02")
        slotsByDate[key] = append(slotsByDate[key], timeSlot{start, end})
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
            start_time, end_time, user_id, service_id, bath)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING booking_id`,
		req.FirstName, req.LastName, req.Email, req.Phone,
		req.Date, timeParts[0], timeParts[1], artistID, serviceID, req.Bath,
	).Scan(&bookingID)

	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}

	var clientID int
	err = db.QueryRow(`
		SELECT client_id FROM clients WHERE phone = $1`,
		req.Phone,
	).Scan(&clientID)

	if err != nil {
		if err == sql.ErrNoRows {
			// No existing client, proceed to insert
			_, err = db.Exec(`
			INSERT INTO clients (firstname, lastname, email, phone)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (phone) DO NOTHING`,
				req.FirstName, req.LastName, req.Email, req.Phone,
			)
			if err != nil {
				http.Error(w, `{"success":false}`, http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
			return
		}
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
        start_time, end_time, user_id, service_id, bath)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING booking_id`,
		req.FirstName, req.LastName, req.Email, req.Phone,
		req.Date, timeParts[0], timeParts[1], artistID, serviceID, req.Bath,
	).Scan(&bookingID)

	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`
		INSERT INTO clients (firstname, lastname, email, phone)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (phone) DO UPDATE
		SET firstname = EXCLUDED.firstname,
		lastname = EXCLUDED.lastname,
		email = EXCLUDED.email;`,
		req.FirstName, req.LastName, req.Email, req.Phone,
	)
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

	rows, err := db.Query("SELECT user_id, name FROM users WHERE is_visible = true ORDER BY user_id")
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

func dashArtistsHandler(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT DISTINCT user_id, name FROM users WHERE is_visible = true ORDER BY user_id")
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
		http.Error(w, `{"success":}`, http.StatusMethodNotAllowed)
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
			b.date, b.start_time, b.end_time, b.service_id, b.bath,
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
			&b.Date, &b.StartTime, &b.EndTime, &b.ServiceID, &b.Bath,
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
					end_time=$5, firstname=$6, lastname=$7, phone=$8, email=$9, bath=$10
		          WHERE booking_id=$11`
		_, err = db.Exec(query, b.Artist, b.Service, b.Date, timeParts[0], timeParts[1], b.FirstName,
			b.LastName, b.Phone, b.Email, b.Bath, id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success":false, "error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO clients (firstname, lastname, email, phone)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (phone) DO UPDATE
			SET firstname = EXCLUDED.firstname,
			lastname = EXCLUDED.lastname,
			email = EXCLUDED.email;`,
			b.FirstName, b.LastName, b.Email, b.Phone,
		)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusInternalServerError)
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

func disableDayHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var req DisableDayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse the date
	bookingDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if req.Status {
		// Insert booking
		_, err := db.Exec(`
			INSERT INTO bookings (user_id, date, service_id, start_time, end_time, firstname, lastname, email, phone)
			VALUES ($1, $2, '19', '09:00', '21:00', '', '', '', '')`,
			req.UserID, bookingDate,
		)
		if err != nil {
			log.Printf("Insert error: %v", err)
			http.Error(w, "Failed to insert booking", http.StatusInternalServerError)
			return
		}
	} else {
		// Delete booking
		_, err := db.Exec(`
			DELETE FROM bookings
			WHERE user_id = $1 AND date = $2 AND start_time = '09:00' AND end_time = '21:00'`,
			req.UserID, bookingDate,
		)
		if err != nil {
			log.Printf("Delete error: %v", err)
			http.Error(w, "Failed to delete booking", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func clientDataHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	phone := "%" + r.URL.Query().Get("phone") + "%"
	rows, err := db.Query(`
			SELECT client_id, firstname, lastname, email, phone 
			FROM clients 
			WHERE phone LIKE $1 
			ORDER BY client_id DESC`, phone)
	if err != nil {
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}

	var clients []Client
	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ClientID, &c.Firstname, &c.Lastname, &c.Email, &c.Phone); err == nil {
			clients = append(clients, c)
		}
	}

	json.NewEncoder(w).Encode(clients)
}

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

}