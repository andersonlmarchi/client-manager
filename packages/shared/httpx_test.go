package shared_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/t-code/client-manager/packages/shared"
)

func TestWriteProblem(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	err := shared.WriteProblem(rec, http.StatusBadRequest, shared.Problem{
		Title:  "Bad Request",
		Detail: "field x is required",
		Code:   string(shared.CodeInvalid),
	})
	if err != nil {
		t.Fatal(err)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Fatalf("Content-Type = %q", ct)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
	var body shared.Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != http.StatusBadRequest || body.Detail != "field x is required" {
		t.Fatalf("body = %+v", body)
	}
}

func TestWriteErrorProblemMapsCodes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err    error
		status int
	}{
		{shared.NewError(shared.CodeInvalid, "bad"), http.StatusBadRequest},
		{shared.NewError(shared.CodeNotFound, "gone"), http.StatusNotFound},
		{shared.NewError(shared.CodeConflict, "dup"), http.StatusConflict},
		{shared.NewError(shared.CodeForbidden, "no"), http.StatusForbidden},
		{shared.NewError(shared.CodeInternal, "boom"), http.StatusInternalServerError},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "/v1/items", nil)
		rec := httptest.NewRecorder()
		if err := shared.WriteErrorProblem(rec, req, tc.err); err != nil {
			t.Fatal(err)
		}
		if rec.Code != tc.status {
			t.Fatalf("code %v: status = %d, want %d", tc.err, rec.Code, tc.status)
		}
	}
}
