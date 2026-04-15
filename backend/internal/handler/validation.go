package handler

import (
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate     *validator.Validate
	validateOnce sync.Once
)

func getValidator() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New()
	})
	return validate
}

// validateRequest validates a struct and writes a 400 response with the first error if invalid.
// Returns true if valid, false if an error was written.
func validateRequest(w http.ResponseWriter, req interface{}) bool {
	if err := getValidator().Struct(req); err != nil {
		var msg string
		if ve, ok := err.(validator.ValidationErrors); ok && len(ve) > 0 {
			msg = ve[0].Translate(nil)
			if msg == "" {
				msg = ve[0].Error()
			}
		} else {
			msg = err.Error()
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
		return false
	}
	return true
}
