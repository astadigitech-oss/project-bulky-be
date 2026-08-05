package dto

import (
	"encoding/json"
	"fmt"
)

// FlexString menerima nilai JSON string ATAU number dan menyimpannya sebagai string.
// Dipakai untuk field id/no_booking webhook Deliveree — provider mengirim booking ID
// sebagai integer (mis. 123456), sementara handler lama (Laravel) menerimanya
// fleksibel. Struct Go dengan tipe string akan gagal unmarshal saat menerima number.
type FlexString string

// UnmarshalJSON menerima JSON string, number, atau null.
func (f *FlexString) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		*f = ""
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return fmt.Errorf("cannot unmarshal %s into FlexString", string(b))
	}
	*f = FlexString(n.String())
	return nil
}

// DelivereeWebhookRequest adalah payload webhook dari Deliveree.
// Mengikuti format yang sama seperti handler webhook di BE lama (panel-bulky).
type DelivereeWebhookRequest struct {
	Status      string     `json:"status"`
	ID          FlexString `json:"id"`
	NoBooking   FlexString `json:"no_booking"`
	TrackingURL string     `json:"tracking_url"`
}

// DelivereeWebhookResponse adalah respons sukses webhook.
type DelivereeWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
