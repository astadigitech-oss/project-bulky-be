package dto

import (
	"encoding/json"
	"testing"
)

func TestFlexStringUnmarshal(t *testing.T) {
	// id dikirim sebagai angka (integer) oleh Deliveree
	payload := `{"status":"delivery_in_progress","id":123456,"tracking_url":"https://track.example.com/abc"}`
	var req DelivereeWebhookRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		t.Fatalf("unmarshal number id gagal: %v", err)
	}
	if string(req.ID) != "123456" {
		t.Fatalf("id number: expected 123456, got %q", string(req.ID))
	}

	// id dikirim sebagai string
	payload2 := `{"status":"delivery_completed","id":"123456"}`
	var req2 DelivereeWebhookRequest
	if err := json.Unmarshal([]byte(payload2), &req2); err != nil {
		t.Fatalf("unmarshal string id gagal: %v", err)
	}
	if string(req2.ID) != "123456" {
		t.Fatalf("id string: expected 123456, got %q", string(req2.ID))
	}

	// no_booking sebagai angka, id kosong
	payload3 := `{"status":"canceled","no_booking":9999}`
	var req3 DelivereeWebhookRequest
	if err := json.Unmarshal([]byte(payload3), &req3); err != nil {
		t.Fatalf("unmarshal no_booking number gagal: %v", err)
	}
	if string(req3.NoBooking) != "9999" {
		t.Fatalf("no_booking: expected 9999, got %q", string(req3.NoBooking))
	}

	// id null
	payload4 := `{"status":"locating_driver","id":null}`
	var req4 DelivereeWebhookRequest
	if err := json.Unmarshal([]byte(payload4), &req4); err != nil {
		t.Fatalf("unmarshal id null gagal: %v", err)
	}
	if string(req4.ID) != "" {
		t.Fatalf("id null: expected empty, got %q", string(req4.ID))
	}
}
