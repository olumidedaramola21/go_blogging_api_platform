package middleware

import (
	"encoding/json"		// for encoding responses to JSON
	"log" 				// for logging request info
	"net/http"			// for handling HTTP request and responses 
	"time"				// for tracking request duration
)

// responseWriter wraps http.ResponseWriter to capture status - useful for logging as status code is not easily captured by http.ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	statusCode int		//  captures the status code returned by the handler
}

// LoggingMiddleware logs HTTP requests details like method, url path, and time taken to process request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()	// records the current time request was given
		// Create a custom ResponseWriter to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)	// passes control to the next middleware or handler	
		duration := time.Since(start)	// how long request took to be processed 
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// AuthenticationMiddleware processes the request before the main API handler and checks if the user is authenticated 
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for Authorization header or session cookie (depending on your authentication method)
		token := r.Header.Get("Authorization")
		if token == "" {
			// if no token is provided, respond with an Unauthorized error
			ErrorResponse(w, http.StatusUnauthorized, "Authorizaation token is missing")
			return 
		}
		// Authentication Code here
	})
}

// JSONMiddleware ensures that every HTTP response has  content type set to JSON
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles Cross-Origin Resource sharing - allows specification of which domain and HTTP methods are allowed to access resource on the server.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			// OPTIONS is used by browser pre-flight CORS
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// ErrorHandlingMiddleware logs any error that occurs in the Main handler
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// catch any panic that occurs in the handler
		defer func()  {
			if err := recover(); err != nil {
				// log the error (you could also add logging here)
				log.Printf("Error: %v", err)

				// Send a generic error response
				ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		// if no errors, continue to the next handler
		next.ServeHTTP(w, r)
	})
}


// Utilities

// ErrorResponse sends a JSON response indicating error with a status code and message
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"success": false,
		"error": message,
	}
	json.NewEncoder(w).Encode(response)	// serialises map into JSON and writes it to the response body
}

// SuccessResponse sends a JSON response indicating success with a status code, message and data
func SuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) { 
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"success": true,
		"message": message,
		"data": data,
	}
	json.NewEncoder(w).Encode(response)	// serialises map into JSON and writes it to the response body
}


