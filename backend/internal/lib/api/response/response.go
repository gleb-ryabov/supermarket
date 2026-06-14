package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response consists of required fields for all responses.
type Response struct {
	Status string `json:"status"`          // Response status
	Error  string `json:"error,omitempty"` // Error if response is error
}

const (
	// StatusOK - Status without error.
	StatusOK = "OK"
	// StatusError - Status with error.
	StatusError = "Error"
)

// OK returns a response with OK status.
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// Error returns a response with the error status and error message.
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// Respond writes HTTP response in JSON format with status code.
func Respond(c *fiber.Ctx, code int, resp any) error {
	return c.Status(code).JSON(resp)
}
