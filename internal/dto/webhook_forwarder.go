package dto

// ForwarderWebhookRequest adalah payload webhook dari Forwarder.
//
// Karena Forwarder & Deliveree berbagi platform on-demand yang sama, format
// payload identik dengan DelivereeWebhookRequest: status pengiriman, identifier
// booking/tracking, dan tracking URL opsional.
type ForwarderWebhookRequest struct {
	Status      string     `json:"status"`
	ID          FlexString `json:"id"`
	NoBooking   FlexString `json:"no_booking"`
	TrackingURL string     `json:"tracking_url"`
}

// ForwarderWebhookResponse adalah respons sukses webhook.
type ForwarderWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
