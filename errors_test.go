package ujeebu

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAPIError_Unmarshal_ErrorCodeNumber(t *testing.T) {
	var e APIError
	if err := json.Unmarshal([]byte(`{"message":"bad","error_code":404}`), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	e.StatusCode = 404

	if got := e.errorCodeString(); got != "404" {
		t.Fatalf("expected error code '404', got %q", got)
	}
	if got := e.Error(); !strings.Contains(got, "code: 404") {
		t.Fatalf("expected Error() to include code: 404, got %q", got)
	}
}

func TestAPIError_Unmarshal_ErrorCodeString(t *testing.T) {
	var e APIError
	if err := json.Unmarshal([]byte(`{"message":"bad","error_code":"AUTH"}`), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	e.StatusCode = 401

	if got := e.errorCodeString(); got != "AUTH" {
		t.Fatalf("expected error code 'AUTH', got %q", got)
	}
	if got := e.Error(); !strings.Contains(got, "code: AUTH") {
		t.Fatalf("expected Error() to include code: AUTH, got %q", got)
	}
}
