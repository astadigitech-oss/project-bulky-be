package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"project-bulky-be/internal/config"
)

type PushTestRecipient struct {
	BuyerID     string `json:"buyer_id"`
	BuyerName   string `json:"buyer_name"`
	Email       string `json:"email"`
	DeviceCount int64  `json:"device_count"`
}

type PushTestSummary struct {
	SuccessCount int `json:"success_count"`
	FailureCount int `json:"failure_count"`
}

type pushTestEnvelope[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type PushTestService struct {
	storefrontBaseURL string
	internalAPIKey    string
	client            *http.Client
}

func NewPushTestService(cfg *config.Config) *PushTestService {
	return &PushTestService{
		storefrontBaseURL: strings.TrimRight(cfg.StorefrontBaseURL, "/"),
		internalAPIKey:    cfg.InternalAPIKey,
		client: &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (s *PushTestService) FindRecipients(ctx context.Context, search string) ([]PushTestRecipient, error) {
	var result []PushTestRecipient
	path := "/internal/notifications/push-test/recipients"
	if search != "" {
		path += "?search=" + url.QueryEscape(search)
	}
	if err := s.do(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	if result == nil {
		result = []PushTestRecipient{}
	}
	return result, nil
}

func (s *PushTestService) Send(ctx context.Context, buyerIDs []string) (*PushTestSummary, error) {
	var result PushTestSummary
	body := map[string]any{"buyer_ids": buyerIDs}
	if err := s.do(ctx, http.MethodPost, "/internal/notifications/push-test/send", body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *PushTestService) NotifyNewProduct(ctx context.Context, productID, slug, name string) error {
	var result PushTestSummary
	return s.do(ctx, http.MethodPost, "/internal/notifications/events/new-product", map[string]string{
		"product_id": productID,
		"slug":       slug,
		"name":       name,
	}, &result)
}

func (s *PushTestService) NotifyOrderStatusChanged(ctx context.Context, buyerID, orderID, orderCode, previousStatus, orderStatus string) error {
	var result PushTestSummary
	return s.do(ctx, http.MethodPost, "/internal/notifications/events/order-status", map[string]string{
		"buyer_id":        buyerID,
		"order_id":        orderID,
		"order_code":      orderCode,
		"previous_status": previousStatus,
		"order_status":    orderStatus,
	}, &result)
}

func (s *PushTestService) do(ctx context.Context, method, path string, body any, output any) error {
	if s.storefrontBaseURL == "" || s.internalAPIKey == "" {
		return fmt.Errorf("integrasi push test dengan storefront belum dikonfigurasi")
	}
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = strings.NewReader(string(encoded))
	}
	req, err := http.NewRequestWithContext(ctx, method, s.storefrontBaseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("X-Internal-Key", s.internalAPIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	response, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal menghubungi storefront: %w", err)
	}
	defer response.Body.Close()
	var envelope pushTestEnvelope[json.RawMessage]
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&envelope); err != nil {
		return fmt.Errorf("respons storefront tidak valid (HTTP %d): %w", response.StatusCode, err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 || !envelope.Success {
		if envelope.Message == "" {
			envelope.Message = http.StatusText(response.StatusCode)
		}
		return fmt.Errorf("storefront menolak permintaan: %s", envelope.Message)
	}
	if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(envelope.Data, output); err != nil {
		return fmt.Errorf("data respons storefront tidak valid: %w", err)
	}
	return nil
}
