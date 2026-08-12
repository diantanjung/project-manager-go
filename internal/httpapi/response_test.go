package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func TestRespondWrapsSingleResource(t *testing.T) {
	c, rec := testContext()

	respond(c, http.StatusOK, gin.H{"id": 1}, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got struct {
		Data map[string]int `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if got.Data["id"] != 1 {
		t.Fatalf("data.id = %d, want %d", got.Data["id"], 1)
	}
}

func TestRespondWrapsPaginatedResource(t *testing.T) {
	c, rec := testContext()
	payload := domain.Paginated[string]{
		Data: []string{"task"},
		Pagination: domain.Pagination{
			Page:       2,
			Limit:      10,
			TotalItems: 11,
			TotalPages: 2,
		},
	}

	respond(c, http.StatusOK, payload, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got struct {
		Data       []string          `json:"data"`
		Pagination domain.Pagination `json:"pagination"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0] != "task" {
		t.Fatalf("data = %#v, want [task]", got.Data)
	}
	if got.Pagination.TotalItems != 11 || got.Pagination.TotalPages != 2 {
		t.Fatalf("pagination = %+v, want totalItems 11 and totalPages 2", got.Pagination)
	}
}

func TestRespondErrorUsesSharedShape(t *testing.T) {
	c, rec := testContext()

	respondError(c, domain.NewError(domain.ErrValidation, "The given data was invalid."))

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	var got struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if got.Message != "The given data was invalid." {
		t.Fatalf("message = %q", got.Message)
	}

	var raw map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decoding raw response: %v", err)
	}
	if _, ok := raw["data"]; ok {
		t.Fatalf("unexpected data envelope in error response: %#v", raw)
	}
}

func TestBindJSONUsesGinValidationTags(t *testing.T) {
	c, rec := testContext()
	c.Request = httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"name":"Dian","email":"not-an-email","password":"secret123"}`),
	)
	c.Request.Header.Set("Content-Type", "application/json")

	var input service.CreateUserInput
	if bindJSON(c, &input) {
		t.Fatal("bindJSON() = true, want false")
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func testContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	return c, rec
}
