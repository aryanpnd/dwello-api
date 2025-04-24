package utils

import (
	"context"
	"time"
)

func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// CustomTimeout creates a context with a custom timeout duration in seconds.
// It returns the context and a cancel function to release resources when done.
func CustomTimeout(seconds int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

// DatabaseContext creates a context with a timeout of 10 seconds for database operations.
// It returns the context and a cancel function to release resources when done.
func DatabaseContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
