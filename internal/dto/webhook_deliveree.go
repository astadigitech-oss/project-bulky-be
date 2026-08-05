package dto

// DelivereeWebhookRequest adalah payload webhook dari Deliveree.
// Mengikuti format yang sama seperti handler webhook di BE lama (panel-bulky).
type DelivereeWebhookRequest struct {
	Status      string `json:"status"`
	ID          string `json:"id"`
	NoBooking   string `json:"no_booking"`
	TrackingURL string `json:"tracking_url"`
}

// DelivereeWebhookResponse adalah respons sukses webhook.
type DelivereeWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
