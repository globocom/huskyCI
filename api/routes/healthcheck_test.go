package routes

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	err := HealthCheck(ctx)
	if err != nil {
		t.Errorf("Error to return the health check.")
	}
	if rec.Code != 200 {
		t.Errorf("Different status code.")
	}
	got := rec.Body.String()
	if got != "WORKING\n" {
		t.Errorf("HealthCheck is not working.")
	}
}