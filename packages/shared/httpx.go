package shared

import (
	"encoding/json"
	"net/http"
)

// Problem is an RFC 7807 problem+json payload.
type Problem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
}

// WriteProblem writes a problem+json response with the given status.
func WriteProblem(w http.ResponseWriter, status int, problem Problem) error {
	problem.Status = status
	if problem.Title == "" {
		problem.Title = http.StatusText(status)
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(problem)
}

// WriteErrorProblem maps a shared.Error (or generic error) to problem+json.
func WriteErrorProblem(w http.ResponseWriter, r *http.Request, err error) error {
	status := http.StatusInternalServerError
	problem := Problem{
		Title:    http.StatusText(status),
		Status:   status,
		Detail:   "internal error",
		Instance: r.URL.Path,
		Code:     string(CodeInternal),
	}

	if e, ok := AsError(err); ok {
		problem.Code = string(e.Code)
		problem.Detail = e.Message
		switch e.Code {
		case CodeInvalid:
			status = http.StatusBadRequest
		case CodeNotFound:
			status = http.StatusNotFound
		case CodeConflict:
			status = http.StatusConflict
		case CodeForbidden:
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError
			problem.Detail = "internal error"
		}
		problem.Title = http.StatusText(status)
		problem.Status = status
	}

	return WriteProblem(w, status, problem)
}
