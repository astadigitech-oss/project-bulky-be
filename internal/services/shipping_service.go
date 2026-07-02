package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"log"

	"project-bulky-be/internal/models"
	"project-bulky-be/pkg/utils"

	"gorm.io/gorm"
)

// TrackingEvent satu entri riwayat status pengiriman
type TrackingEvent struct {
	Date   string `json:"date"`
	Time   string `json:"time"`
	Status string `json:"status"`
}

// TrackingResult hasil tracking pengiriman dari provider
type TrackingResult struct {
	Provider    string         `json:"provider"`
	BookingRef  string         `json:"booking_ref"`
	Status      string         `json:"status"`
	TrackingURL string         `json:"tracking_url,omitempty"`
	History     []TrackingEvent `json:"history"`
}

// DelivereeVehicleTypeInfo is the vehicle info returned by Deliveree.
type DelivereeVehicleTypeInfo struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	CargoLength    float64 `json:"cargo_length"`
	CargoHeight    float64 `json:"cargo_height"`
	CargoWidth     float64 `json:"cargo_width"`
	CargoWeight    float64 `json:"cargo_weight"`
	CargoCubicMeter float64 `json:"cargo_cubic_meter"`
}

// DelivereeDriver is the driver info returned by Deliveree.
type DelivereeDriver struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	Phone                  string  `json:"phone"`
	DriverImageURL         string  `json:"driver_image_url"`
	LastKnownPositionLat   float64 `json:"last_known_position_lat"`
	LastKnownPositionLng   float64 `json:"last_known_position_lng"`
}

// DelivereeDeliveryLocation is a location entry in the Deliveree delivery detail.
type DelivereeDeliveryLocation struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name"`
	DriverNote          string  `json:"driver_note"`
	Note                string  `json:"note"`
	RecipientName       string  `json:"recipient_name"`
	RecipientPhone      string  `json:"recipient_phone"`
	DeliveryStatus      string  `json:"delivery_status"`
	FailedDeliveryReason string `json:"failed_delivery_reason"`
	SignatureURL        string  `json:"signature_url"`
	ArrivedAt          string  `json:"arrived_at"`
	LeavedAt           string  `json:"leaved_at"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	ParkingFees        float64 `json:"parking_fees"`
	TollsFees          float64 `json:"tolls_fees"`
	WaitingTimeFees    float64 `json:"waiting_time_fees"`
	TrackingSharing    string  `json:"tracking_sharing"`
}

// DelivereeDeliveryDetail is the full delivery detail response from Deliveree API.
type DelivereeDeliveryDetail struct {
	ID                int                         `json:"id"`
	CustomerName      string                      `json:"customer_name"`
	DriverID          int                         `json:"driver_id"`
	VehicleTypeInfo   DelivereeVehicleTypeInfo     `json:"vehicle_type_info"`
	TimeType          string                      `json:"time_type"`
	Status            string                      `json:"status"`
	Note              string                      `json:"note"`
	TotalFees         float64                     `json:"total_fees"`
	Currency          string                      `json:"currency"`
	TrackingURL       string                      `json:"tracking_url"`
	JobOrderNumber    string                      `json:"job_order_number"`
	CreatedAt         string                      `json:"created_at"`
	PickupTime        string                      `json:"pickup_time"`
	CompletedAt       string                      `json:"completed_at"`
	Driver            *DelivereeDriver             `json:"driver"`
	Locations         []DelivereeDeliveryLocation `json:"locations"`
	RequireSignatures bool                        `json:"require_signatures"`
	DistanceFees      float64                     `json:"distance_fees"`
	CODPODFees        float64                     `json:"cod_pod_fees"`
	CODPOD            bool                        `json:"cod_pod"`
	SurchargesFees    float64                     `json:"surcharges_fees"`
	WayPointFees      float64                     `json:"way_point_fees"`
}

type ShippingService interface {
	// TriggerBookingAsync launches booking in a goroutine (non-blocking).
	TriggerBookingAsync(pesanan *models.Pesanan)
	// BookDelivery runs the booking synchronously and returns the result.
	BookDelivery(ctx context.Context, pesanan *models.Pesanan) (delivereeBookingID *string, forwarderTrackingNo *string, err error)
	// TrackDelivery retrieves live tracking info from the shipping provider.
	TrackDelivery(ctx context.Context, pesanan *models.Pesanan) (*TrackingResult, error)
	// GetDelivereeDetail retrieves full delivery detail from Deliveree API.
	GetDelivereeDetail(ctx context.Context, pesanan *models.Pesanan) (*DelivereeDeliveryDetail, error)
	// GetForwarderInvoice retrieves invoice list from Forwarder by booking number.
	GetForwarderInvoice(ctx context.Context, bookingNo string) ([]ForwarderInvoice, error)
}

type shippingService struct {
	db *gorm.DB
}

func NewShippingService(db *gorm.DB) ShippingService {
	return &shippingService{db: db}
}

var jawaBaliProvinces = map[string]bool{
	"dki jakarta":                true,
	"jakarta":                    true,
	"jawa barat":                 true,
	"jawa tengah":                true,
	"jawa timur":                 true,
	"d.i. yogyakarta":            true,
	"di yogyakarta":              true,
	"daerah istimewa yogyakarta": true,
	"yogyakarta":                 true,
	"banten":                     true,
	"bali":                       true,
}

func isLuarJawaBali(provinsi string) bool {
	return !jawaBaliProvinces[strings.ToLower(strings.TrimSpace(provinsi))]
}

func (s *shippingService) TriggerBookingAsync(pesanan *models.Pesanan) {
	go func(p *models.Pesanan) {
		// Skip jika booking sudah ada — cegah double-booking jika status READY dipanggil ulang
		if p.DelivereeBookingID != nil || p.ForwarderTrackingNo != nil {
			log.Printf("[shipping] skip booking async: pesanan=%s sudah punya booking_id/tracking_no", p.Kode)
			return
		}

		log.Printf("[shipping] trigger booking async: pesanan=%s delivery_type=%s", p.Kode, p.DeliveryType)
		ctx := context.Background()
		delivereeID, trackingNo, err := s.BookDelivery(ctx, p)
		updates := map[string]interface{}{}
		if err != nil {
			log.Printf("[shipping] booking gagal: pesanan=%s delivery_type=%s error=%v", p.Kode, p.DeliveryType, err)
			errMsg := err.Error()
			updates["booking_error"] = errMsg
		} else {
			log.Printf("[shipping] booking sukses: pesanan=%s delivery_type=%s deliveree_id=%v tracking_no=%v", p.Kode, p.DeliveryType, delivereeID, trackingNo)
			updates["booking_error"] = nil
			if delivereeID != nil {
				updates["deliveree_booking_id"] = *delivereeID
			}
			if trackingNo != nil {
				updates["forwarder_tracking_no"] = *trackingNo
			}
		}
		s.db.Model(&models.Pesanan{}).Where("id = ?", p.ID).UpdateColumns(updates)
	}(pesanan)
}

func (s *shippingService) BookDelivery(ctx context.Context, pesanan *models.Pesanan) (*string, *string, error) {
	switch pesanan.DeliveryType {
	case models.DeliveryTypeDeliveree:
		bookingID, err := s.bookDeliveree(ctx, pesanan)
		if err != nil {
			return nil, nil, err
		}
		return bookingID, nil, nil
	case models.DeliveryTypeForwarder, models.DeliveryTypeForwarderLCL:
		var trackingNo *string
		var err error
		if pesanan.AlamatBuyer != nil && isLuarJawaBali(pesanan.AlamatBuyer.Provinsi) {
			trackingNo, err = s.bookForwarderLCL(ctx, pesanan)
		} else {
			trackingNo, err = s.bookForwarder(ctx, pesanan)
		}
		if err != nil {
			return nil, nil, err
		}
		return nil, trackingNo, nil
	default:
		return nil, nil, fmt.Errorf("delivery type %s tidak memerlukan booking", pesanan.DeliveryType)
	}
}

// ─── Deliveree ────────────────────────────────────────────────────────────────

type delivereeLocation struct {
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	RecipientName  string  `json:"recipient_name"`
	RecipientPhone string  `json:"recipient_phone"`
	Note           string  `json:"note,omitempty"`
	IsPayer        bool    `json:"is_payer"`
	NeedCOD        bool    `json:"need_cod,omitempty"`
	CODNote        string  `json:"cod_note,omitempty"`
	CODInvoiceFees float64 `json:"cod_invoice_fees,omitempty"`
	NeedPOD        bool    `json:"need_pod,omitempty"`
	PODNote        string  `json:"pod_note,omitempty"`
}

type delivereeCreateRequest struct {
	VehicleTypeID        int                 `json:"vehicle_type_id"`
	BookingPaymentType   string              `json:"booking_payment_type"`
	Note                 string              `json:"note,omitempty"`
	TimeType             string              `json:"time_type"`
	PickupTime           string              `json:"pickup_time"`
	JobOrderNumber       string              `json:"job_order_number,omitempty"`
	AllowParkingFees     bool                `json:"allow_parking_fees"`
	AllowTollsFees       bool                `json:"allow_tolls_fees"`
	AllowWaitingTimeFees bool                `json:"allow_waiting_time_fees"`
	SendFirstToDriver    bool                `json:"send_first_to_driver"`
	MarkedAsFavorite     bool                `json:"marked_as_favorite"`
	Locations            []delivereeLocation `json:"locations"`
	RequireSignatures    bool                `json:"require_signatures"`
	ExtraServices        []interface{}       `json:"extra_services"`
	EstimateTransitTimes []interface{}       `json:"estimate_transit_times_attributes"`
}

type delivereeCreateResponse struct {
	BookingID int `json:"booking_id"`
}

func (s *shippingService) bookDeliveree(ctx context.Context, pesanan *models.Pesanan) (*string, error) {
	baseURL := os.Getenv("DELIVEREE_BASE_URL")
	apiKey := os.Getenv("DELIVEREE_API_KEY")
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("konfigurasi Deliveree tidak lengkap")
	}

	// Get first active warehouse
	warehouse, err := s.getActiveWarehouse()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data warehouse: %w", err)
	}

	// AlamatBuyer is required
	if pesanan.AlamatBuyer == nil {
		return nil, fmt.Errorf("pesanan tidak memiliki alamat pengiriman")
	}
	alamat := pesanan.AlamatBuyer

	// Vehicle type based on total qty
	totalQty := 0
	for _, item := range pesanan.Items {
		totalQty += item.Qty
	}
	vehicleTypeID := delivereeVehicleTypeID(totalQty, baseURL)

	warehouseAlamat := ""
	if warehouse.Alamat != nil {
		warehouseAlamat = *warehouse.Alamat
	}
	warehouseTelepon := ""
	if warehouse.Telepon != nil {
		warehouseTelepon = *warehouse.Telepon
	}
	warehouseLat := 0.0
	if warehouse.Latitude != nil {
		warehouseLat = *warehouse.Latitude
	}
	warehouseLng := 0.0
	if warehouse.Longitude != nil {
		warehouseLng = *warehouse.Longitude
	}

	buyerLat := 0.0
	if alamat.Latitude != nil {
		buyerLat = *alamat.Latitude
	}
	buyerLng := 0.0
	if alamat.Longitude != nil {
		buyerLng = *alamat.Longitude
	}

	req := delivereeCreateRequest{
		VehicleTypeID:        vehicleTypeID,
		BookingPaymentType:   "credit",
		TimeType:             "now",
		PickupTime:           "",
		JobOrderNumber:       pesanan.Kode,
		AllowParkingFees:     true,
		AllowTollsFees:       true,
		AllowWaitingTimeFees: true,
		SendFirstToDriver:    false,
		MarkedAsFavorite:     true,
		RequireSignatures:    true,
		ExtraServices:        []interface{}{},
		EstimateTransitTimes: []interface{}{},
		Locations: []delivereeLocation{
			{
				Address:        warehouseAlamat,
				Latitude:       warehouseLat,
				Longitude:      warehouseLng,
				RecipientName:  "Warehouse",
				RecipientPhone: warehouseTelepon,
				Note:           "Pickup at warehouse",
				IsPayer:        false,
			},
			{
				Address:        alamat.AlamatLengkap,
				Latitude:       buyerLat,
				Longitude:      buyerLng,
				RecipientName:  alamat.NamaPenerima,
				RecipientPhone: alamat.TeleponPenerima,
				Note:           "Drop at buyer location",
				IsPayer:        true,
				NeedCOD:        false,
				NeedPOD:        false,
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/deliveries", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", apiKey)

	log.Printf("[deliveree] --> POST /deliveries pesanan=%s vehicle_type_id=%d body=%s", pesanan.Kode, req.VehicleTypeID, string(body))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[deliveree] <-- POST /deliveries pesanan=%s error=%v", pesanan.Kode, err)
		return nil, fmt.Errorf("connection timeout atau gagal menghubungi Deliveree: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[deliveree] <-- POST /deliveries pesanan=%s status=%d body=%s", pesanan.Kode, resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Deliveree API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result delivereeCreateResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response Deliveree: %w", err)
	}
	if result.BookingID == 0 {
		return nil, fmt.Errorf("Deliveree tidak mengembalikan booking ID")
	}

	bookingID := strconv.Itoa(result.BookingID)
	return &bookingID, nil
}

func delivereeVehicleTypeID(totalQty int, baseURL string) int {
	isSandbox := strings.Contains(baseURL, "sandbox")
	switch {
	case totalQty <= 4:
		if isSandbox {
			return 14
		}
		return 2701
	case totalQty <= 8:
		if isSandbox {
			return 24
		}
		return 2703
	default:
		if isSandbox {
			return 36
		}
		return 2723
	}
}

// ─── Forwarder ────────────────────────────────────────────────────────────────

type forwarderTokenRequest struct {
	Scope string `json:"scope"`
}

type forwarderTokenResponse struct {
	AccessToken string `json:"access_token"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type forwarderLocation struct {
	Address   string `json:"address"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Type      string `json:"type"`
	Order     int    `json:"order"`
	PicName   string `json:"picname"`
	PicPhone  string `json:"picphone"`
	Detail    string `json:"detail,omitempty"`
}

type forwarderDataDetail struct {
	Packaging   string  `json:"packaging"`
	Commodity   string  `json:"commodity"`
	CargoDesc   string  `json:"cargodesc"`
	Qty         int     `json:"qty"`
	Length      float64 `json:"length"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Volume      float64 `json:"volume"`
	TotalVolume float64 `json:"totalvolume"`
	Weight      float64 `json:"weight"`
	TotalWeight float64 `json:"totalweight"`
}

type forwarderCreateBookingRequest struct {
	TransportID              string                `json:"transportid"`
	LoadID                   string                `json:"loadid"`
	ServiceID                string                `json:"serviceid"`
	OriginCityID             string                `json:"origincityid"`
	DestinationCityID        string                `json:"destinationcityid"`
	DestinationSubdistrictID string                `json:"destinationsubdistrictid"`
	PriceID                  int                   `json:"priceid"`
	PickupTimeType           string                `json:"pickuptimetype"`
	PickupTimeOn             string                `json:"pickuptimeon"`
	VehicleID                string                `json:"vehicleid"`
	VehicleQty               int                   `json:"vehicleqty"`
	ShipperName              string                `json:"shippername"`
	ShipperPhone             string                `json:"shipperphone"`
	ShipperAddress           string                `json:"shipperaddress"`
	ConsigneeName            string                `json:"consigneename"`
	ConsigneePhone           string                `json:"consigneephone"`
	ConsigneeAddress         string                `json:"consigneeaddress"`
	EstDistance              string                `json:"estdistance"`
	EstPrice                 string                `json:"estprice"`
	BasisPrice               string                `json:"basisprice"`
	Remark                   string                `json:"remark,omitempty"`
	WithInsurance            string                `json:"withinsurance"`
	CommodityAmount          string                `json:"commodityamount,omitempty"`
	InsuranceID              string                `json:"insuranceid,omitempty"`
	PremiAmount              string                `json:"premiamount,omitempty"`
	Locations                []forwarderLocation   `json:"locations"`
	DataDetail               []forwarderDataDetail `json:"datadetail"`
}

type forwarderBookingResponse struct {
	Msg       string `json:"msg"`
	Data      struct {
		BookingNo string `json:"booking_no"`
	} `json:"data"`
	IsSuccess string `json:"isSuccess"`
}

func (s *shippingService) bookForwarder(ctx context.Context, pesanan *models.Pesanan) (*string, error) {
	apiURL := os.Getenv("FORWARDER_API_URL")
	clientName := os.Getenv("FORWARDER_CLIENT_NAME")
	username := os.Getenv("FORWARDER_USERNAME")
	password := os.Getenv("FORWARDER_PASSWORD")
	if apiURL == "" || clientName == "" || username == "" || password == "" {
		return nil, fmt.Errorf("konfigurasi Forwarder tidak lengkap")
	}

	// Get access token
	token, err := s.getForwarderToken(ctx, apiURL, clientName, username, password, "CREATEBOOKINGLAND")
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan token Forwarder: %w", err)
	}

	// Get first active warehouse
	warehouse, err := s.getActiveWarehouse()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data warehouse: %w", err)
	}

	if pesanan.AlamatBuyer == nil {
		return nil, fmt.Errorf("pesanan tidak memiliki alamat pengiriman")
	}
	alamat := pesanan.AlamatBuyer

	// Lookup origin city (warehouse)
	if warehouse.Kota == nil || *warehouse.Kota == "" {
		return nil, fmt.Errorf("warehouse tidak memiliki data kota")
	}
	var originMapping models.ForwarderCityMapping
	if err := s.db.Where("kota_pattern = ?", utils.NormalizeKota(*warehouse.Kota)).First(&originMapping).Error; err != nil {
		return nil, fmt.Errorf("kota asal warehouse (%s) tidak ditemukan di Forwarder mapping", *warehouse.Kota)
	}

	// Lookup destination city (buyer)
	var destMapping models.ForwarderCityMapping
	if err := s.db.Where("kota_pattern = ?", utils.NormalizeKota(alamat.Kota)).First(&destMapping).Error; err != nil {
		return nil, fmt.Errorf("kota tujuan tidak ditemukan di Forwarder mapping. Silakan tambahkan mapping untuk kota: %s", alamat.Kota)
	}

	// Lookup destination subdistrict (optional for LTL)
	subdistrictID := "0"
	if alamat.Kecamatan != nil && *alamat.Kecamatan != "" {
		subdistrictID = s.resolveForwarderSubdistrictID(ctx, *alamat.Kecamatan, destMapping.ForwarderCityID, apiURL, clientName, username, password)
	}

	// Build datadetail from items
	dataDetail := make([]forwarderDataDetail, 0, len(pesanan.Items))
	for _, item := range pesanan.Items {
		p := item.Produk
		vol := (p.Panjang * p.Lebar * p.Tinggi) / 1_000_000 // cm³ → m³
		dataDetail = append(dataDetail, forwarderDataDetail{
			Packaging:   "5",
			Commodity:   "155",
			CargoDesc:   item.NamaProduk,
			Qty:         item.Qty,
			Length:      p.Panjang,
			Width:       p.Lebar,
			Height:      p.Tinggi,
			Volume:      vol,
			TotalVolume: vol * float64(item.Qty),
			Weight:      p.Berat,
			TotalWeight: p.Berat * float64(item.Qty),
		})
	}

	warehouseAlamat := ""
	if warehouse.Alamat != nil {
		warehouseAlamat = *warehouse.Alamat
	}
	warehouseTelepon := ""
	if warehouse.Telepon != nil {
		warehouseTelepon = *warehouse.Telepon
	}
	warehouseLat := "0"
	if warehouse.Latitude != nil {
		warehouseLat = strconv.FormatFloat(*warehouse.Latitude, 'f', 8, 64)
	}
	warehouseLng := "0"
	if warehouse.Longitude != nil {
		warehouseLng = strconv.FormatFloat(*warehouse.Longitude, 'f', 8, 64)
	}
	buyerLat := "0"
	if alamat.Latitude != nil {
		buyerLat = strconv.FormatFloat(*alamat.Latitude, 'f', 8, 64)
	}
	buyerLng := "0"
	if alamat.Longitude != nil {
		buyerLng = strconv.FormatFloat(*alamat.Longitude, 'f', 8, 64)
	}

	withInsurance := "0"
	commodityAmount := ""
	insuranceID := ""
	premiAmount := ""
	if pesanan.BiayaLainnya.IsPositive() {
		withInsurance = "1"
		commodityAmount = pesanan.BiayaProduk.StringFixed(0)
		insuranceID = "1"
		premiAmount = pesanan.BiayaLainnya.StringFixed(0)
	}

	bookingReq := forwarderCreateBookingRequest{
		TransportID:              "3",
		LoadID:                   "5",
		ServiceID:                "1",
		OriginCityID:             strconv.Itoa(originMapping.ForwarderCityID),
		DestinationCityID:        strconv.Itoa(destMapping.ForwarderCityID),
		DestinationSubdistrictID: subdistrictID,
		PriceID:                  1,
		PickupTimeType:           "SCHEDULE",
		PickupTimeOn:             "",
		VehicleID:                "0",
		VehicleQty:               1,
		ShipperName:              "Liquid8",
		ShipperPhone:             warehouseTelepon,
		ShipperAddress:           warehouseAlamat,
		ConsigneeName:            alamat.NamaPenerima,
		ConsigneePhone:           alamat.TeleponPenerima,
		ConsigneeAddress:         alamat.AlamatLengkap,
		EstDistance:              "0",
		EstPrice:                 "",
		BasisPrice:               "ECONOMY",
		Remark:                   pesanan.Kode,
		WithInsurance:            withInsurance,
		CommodityAmount:          commodityAmount,
		InsuranceID:              insuranceID,
		PremiAmount:              premiAmount,
		Locations: []forwarderLocation{
			{
				Address:   warehouseAlamat,
				Latitude:  warehouseLat,
				Longitude: warehouseLng,
				Type:      "PICKUP",
				Order:     1,
				PicName:   "Liquid8",
				PicPhone:  warehouseTelepon,
				Detail:    warehouseAlamat,
			},
			{
				Address:   alamat.AlamatLengkap,
				Latitude:  buyerLat,
				Longitude: buyerLng,
				Type:      "DELIVERY",
				Order:     2,
				PicName:   alamat.NamaPenerima,
				PicPhone:  alamat.TeleponPenerima,
				Detail:    alamat.AlamatLengkap,
			},
		},
		DataDetail: dataDetail,
	}

	body, err := json.Marshal(bookingReq)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/createbookingland", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Client_name", clientName)

	log.Printf("[forwarder] --> POST /createbookingland pesanan=%s transport=%s load=%s origin_city=%s dest_city=%s consignee=%s body=%s",
		pesanan.Kode, bookingReq.TransportID, bookingReq.LoadID,
		bookingReq.OriginCityID, bookingReq.DestinationCityID, bookingReq.ConsigneeName, string(body))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[forwarder] <-- POST /createbookingland pesanan=%s error=%v", pesanan.Kode, err)
		return nil, fmt.Errorf("connection timeout atau gagal menghubungi Forwarder: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[forwarder] <-- POST /createbookingland pesanan=%s status=%d body=%s", pesanan.Kode, resp.StatusCode, string(respBody))

	var result forwarderBookingResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response Forwarder: %w", err)
	}

	if result.IsSuccess != "ok" {
		return nil, fmt.Errorf("Forwarder API error: %s", result.Msg)
	}
	if result.Data.BookingNo == "" {
		return nil, fmt.Errorf("Forwarder tidak mengembalikan booking number")
	}

	return &result.Data.BookingNo, nil
}

// ─── Forwarder LCL (Sea Freight / Luar Jawa) ──────────────────────────────────

type forwarderBookingDetail struct {
	Qty             string  `json:"qty"`
	ContainerTypeID string  `json:"containertypeid"`
	PackageID       string  `json:"packageid"`
	Length          float64 `json:"length"`
	Width           float64 `json:"width"`
	Height          float64 `json:"height"`
	Volume          float64 `json:"volume"`
	Weight          float64 `json:"weight"`
	CargoID         string  `json:"cargoid"`
	CargoDesc       string  `json:"cargodesc"`
}

type forwarderCreateBookingLCLRequest struct {
	TransportID              string                   `json:"transportid"`
	MoveTypeID               string                   `json:"movetypeid"`
	LoadTypeID               string                   `json:"loadtypeid"`
	ServiceTypeID            string                   `json:"servicetypeid"`
	OriginCityID             string                   `json:"origincityid"`
	DestinationCityID        string                   `json:"destinationcityid"`
	DestinationSubdistrictID string                   `json:"destinationsubdistrictid"`
	LCLBasisID               string                   `json:"lclbasisid"`
	CargoReadyDate           string                   `json:"cargoreadydate"`
	Shipper                  string                   `json:"shipper"`
	ShipperAddress           string                   `json:"shipperaddress"`
	ShipperLat               string                   `json:"shipperlat"`
	ShipperLng               string                   `json:"shipperlng"`
	ShipperCountry           string                   `json:"shippercountry"`
	ShipperProvince          string                   `json:"shipperprovince"`
	ShipperCity              string                   `json:"shippercity"`
	ShipperPostalCode        string                   `json:"shipperpostalcode"`
	ShipperRemark            string                   `json:"shipperremark"`
	Consignee                string                   `json:"consignee"`
	ConsigneeAddress         string                   `json:"consigneeaddress"`
	ConsigneeLat             string                   `json:"consigneelat"`
	ConsigneeLng             string                   `json:"consigneelng"`
	ConsigneeCountry         string                   `json:"consigneecountry"`
	ConsigneeProvince        string                   `json:"consigneeprovince"`
	ConsigneeCity            string                   `json:"consigneecity"`
	ConsigneePostalCode      string                   `json:"consigneepostalcode"`
	ConsigneeRemark          string                   `json:"consigneeremark"`
	Pickup                   string                   `json:"pickup"`
	PickupAddress            string                   `json:"pickupaddress"`
	PickupLat                string                   `json:"pickuplat"`
	PickupLng                string                   `json:"pickuplng"`
	PickupCountry            string                   `json:"pickupcountry"`
	PickupProvince           string                   `json:"pickupprovince"`
	PickupCity               string                   `json:"pickupcity"`
	PickupPostalCode         string                   `json:"pickuppostalcode"`
	PickupPhone              string                   `json:"pickupphone"`
	PickupRemark             string                   `json:"pickupremark"`
	Delivery                 string                   `json:"delivery"`
	DeliveryAddress          string                   `json:"deliveryaddress"`
	DeliveryLat              string                   `json:"deliverylat"`
	DeliveryLng              string                   `json:"deliverylng"`
	DeliveryCountry          string                   `json:"deliverycountry"`
	DeliveryProvince         string                   `json:"deliveryprovince"`
	DeliveryCity             string                   `json:"deliverycity"`
	DeliveryPostalCode       string                   `json:"deliverypostalcode"`
	DeliveryPhone            string                   `json:"deliveryphone"`
	DeliveryRemark           string                   `json:"deliveryremark"`
	VoucherCode              string                   `json:"vouchercode"`
	CurrencyID               string                   `json:"currencyid"`
	Incoterm                 string                   `json:"incoterm"`
	WithInsurance            string                   `json:"withinsurance"`
	CommodityAmount          string                   `json:"commodityamount,omitempty"`
	InsuranceID              string                   `json:"insuranceid,omitempty"`
	PremiAmount              string                   `json:"premiamount,omitempty"`
	BookingDetail            []forwarderBookingDetail `json:"bookingdetail"`
}

func (s *shippingService) bookForwarderLCL(ctx context.Context, pesanan *models.Pesanan) (*string, error) {
	apiURL := os.Getenv("FORWARDER_API_URL")
	clientName := os.Getenv("FORWARDER_CLIENT_NAME")
	username := os.Getenv("FORWARDER_USERNAME")
	password := os.Getenv("FORWARDER_PASSWORD")
	if apiURL == "" || clientName == "" || username == "" || password == "" {
		return nil, fmt.Errorf("konfigurasi Forwarder tidak lengkap")
	}

	token, err := s.getForwarderToken(ctx, apiURL, clientName, username, password, "CREATEBOOKING")
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan token Forwarder: %w", err)
	}

	warehouse, err := s.getActiveWarehouse()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data warehouse: %w", err)
	}

	if pesanan.AlamatBuyer == nil {
		return nil, fmt.Errorf("pesanan tidak memiliki alamat pengiriman")
	}
	alamat := pesanan.AlamatBuyer

	if warehouse.Kota == nil || *warehouse.Kota == "" {
		return nil, fmt.Errorf("warehouse tidak memiliki data kota")
	}
	var originMapping models.ForwarderCityMapping
	if err := s.db.Where("kota_pattern = ?", utils.NormalizeKota(*warehouse.Kota)).First(&originMapping).Error; err != nil {
		return nil, fmt.Errorf("kota asal warehouse (%s) tidak ditemukan di Forwarder mapping", *warehouse.Kota)
	}

	var destMapping models.ForwarderCityMapping
	if err := s.db.Where("kota_pattern = ?", utils.NormalizeKota(alamat.Kota)).First(&destMapping).Error; err != nil {
		return nil, fmt.Errorf("kota tujuan tidak ditemukan di Forwarder mapping. Silakan tambahkan mapping untuk kota: %s", alamat.Kota)
	}

	subdistrictID := "0"
	if alamat.Kecamatan != nil && *alamat.Kecamatan != "" {
		subdistrictID = s.resolveForwarderSubdistrictID(ctx, *alamat.Kecamatan, destMapping.ForwarderCityID, apiURL, clientName, username, password)
	}

	bookingDetail := make([]forwarderBookingDetail, 0, len(pesanan.Items))
	for _, item := range pesanan.Items {
		p := item.Produk
		vol := (p.Panjang * p.Lebar * p.Tinggi) / 1_000_000
		bookingDetail = append(bookingDetail, forwarderBookingDetail{
			Qty:             strconv.Itoa(item.Qty),
			ContainerTypeID: "15",
			PackageID:       "7",
			Length:          p.Panjang,
			Width:           p.Lebar,
			Height:          p.Tinggi,
			Volume:          vol * float64(item.Qty),
			Weight:          p.Berat * float64(item.Qty),
			CargoID:         "78",
			CargoDesc:       item.NamaProduk,
		})
	}

	warehouseAlamat := ""
	if warehouse.Alamat != nil {
		warehouseAlamat = *warehouse.Alamat
	}
	warehouseTelepon := ""
	if warehouse.Telepon != nil {
		warehouseTelepon = *warehouse.Telepon
	}
	warehouseKodePos := ""
	if warehouse.KodePos != nil {
		warehouseKodePos = *warehouse.KodePos
	}
	warehouseLat := "0"
	if warehouse.Latitude != nil {
		warehouseLat = strconv.FormatFloat(*warehouse.Latitude, 'f', 8, 64)
	}
	warehouseLng := "0"
	if warehouse.Longitude != nil {
		warehouseLng = strconv.FormatFloat(*warehouse.Longitude, 'f', 8, 64)
	}
	buyerLat := "0"
	if alamat.Latitude != nil {
		buyerLat = strconv.FormatFloat(*alamat.Latitude, 'f', 8, 64)
	}
	buyerLng := "0"
	if alamat.Longitude != nil {
		buyerLng = strconv.FormatFloat(*alamat.Longitude, 'f', 8, 64)
	}
	buyerKodePos := ""
	if alamat.KodePos != nil {
		buyerKodePos = *alamat.KodePos
	}

	withInsurance := "0"
	commodityAmount := ""
	insuranceID := ""
	premiAmount := ""
	if pesanan.BiayaLainnya.IsPositive() {
		withInsurance = "1"
		commodityAmount = pesanan.BiayaProduk.StringFixed(0)
		insuranceID = "1"
		premiAmount = pesanan.BiayaLainnya.StringFixed(0)
	}

	bookingReq := forwarderCreateBookingLCLRequest{
		TransportID:              "1",
		MoveTypeID:               "1",
		LoadTypeID:               "2",
		ServiceTypeID:            "1",
		OriginCityID:             strconv.Itoa(originMapping.ForwarderCityID),
		DestinationCityID:        strconv.Itoa(destMapping.ForwarderCityID),
		DestinationSubdistrictID: subdistrictID,
		LCLBasisID:               "1",
		CargoReadyDate:           "",
		Shipper:                  "Liquid8",
		ShipperAddress:           warehouseAlamat,
		ShipperLat:               warehouseLat,
		ShipperLng:               warehouseLng,
		ShipperCountry:           "Indonesia",
		ShipperProvince:          "",
		ShipperCity:              *warehouse.Kota,
		ShipperPostalCode:        warehouseKodePos,
		ShipperRemark:            "",
		Consignee:                alamat.NamaPenerima,
		ConsigneeAddress:         alamat.AlamatLengkap,
		ConsigneeLat:             buyerLat,
		ConsigneeLng:             buyerLng,
		ConsigneeCountry:         "Indonesia",
		ConsigneeProvince:        alamat.Provinsi,
		ConsigneeCity:            alamat.Kota,
		ConsigneePostalCode:      buyerKodePos,
		ConsigneeRemark:          "",
		Pickup:                   "Liquid8",
		PickupAddress:            warehouseAlamat,
		PickupLat:                warehouseLat,
		PickupLng:                warehouseLng,
		PickupCountry:            "Indonesia",
		PickupProvince:           "",
		PickupCity:               *warehouse.Kota,
		PickupPostalCode:         warehouseKodePos,
		PickupPhone:              warehouseTelepon,
		PickupRemark:             "",
		Delivery:                 alamat.NamaPenerima,
		DeliveryAddress:          alamat.AlamatLengkap,
		DeliveryLat:              buyerLat,
		DeliveryLng:              buyerLng,
		DeliveryCountry:          "Indonesia",
		DeliveryProvince:         alamat.Provinsi,
		DeliveryCity:             alamat.Kota,
		DeliveryPostalCode:       buyerKodePos,
		DeliveryPhone:            alamat.TeleponPenerima,
		DeliveryRemark:           "",
		VoucherCode:              "",
		CurrencyID:               "1",
		Incoterm:                 "",
		WithInsurance:            withInsurance,
		CommodityAmount:          commodityAmount,
		InsuranceID:              insuranceID,
		PremiAmount:              premiAmount,
		BookingDetail:            bookingDetail,
	}

	body, err := json.Marshal(bookingReq)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/createbooking", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Client_name", clientName)

	log.Printf("[forwarder] --> POST /createbooking pesanan=%s transport=%s load=%s origin_city=%s dest_city=%s consignee=%s body=%s",
		pesanan.Kode, bookingReq.TransportID, bookingReq.LoadTypeID,
		bookingReq.OriginCityID, bookingReq.DestinationCityID, bookingReq.Consignee, string(body))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[forwarder] <-- POST /createbooking pesanan=%s error=%v", pesanan.Kode, err)
		return nil, fmt.Errorf("connection timeout atau gagal menghubungi Forwarder: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[forwarder] <-- POST /createbooking pesanan=%s status=%d body=%s", pesanan.Kode, resp.StatusCode, string(respBody))

	var result forwarderBookingResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response Forwarder: %w", err)
	}

	if result.IsSuccess != "ok" {
		return nil, fmt.Errorf("Forwarder API error: %s", result.Msg)
	}
	if result.Data.BookingNo == "" {
		return nil, fmt.Errorf("Forwarder tidak mengembalikan booking number")
	}

	return &result.Data.BookingNo, nil
}

// resolveForwarderSubdistrictID mencari subdistrict ID dengan 3 langkah fallback:
// 1. DB lookup by kecamatan_pattern + forwarder_city_id
// 2. DB lookup by kecamatan_pattern saja (any city)
// 3. API call ke /subdistrictlist
func (s *shippingService) resolveForwarderSubdistrictID(ctx context.Context, kecamatan string, forwarderCityID int, apiURL, clientName, username, password string) string {
	normalized := utils.NormalizeKecamatan(kecamatan)
	if normalized == "" {
		return "0"
	}

	// Step 1: match by kecamatan + city_id
	var m models.ForwarderSubdistrictMapping
	if err := s.db.Where("kecamatan_pattern = ? AND forwarder_city_id = ?", normalized, forwarderCityID).First(&m).Error; err == nil {
		return strconv.Itoa(m.ForwarderSubdistrictID)
	}

	// Step 2: match by kecamatan saja — hanya pakai kalau hasilnya tunggal
	var all []models.ForwarderSubdistrictMapping
	if err := s.db.Where("kecamatan_pattern = ?", normalized).Find(&all).Error; err == nil && len(all) == 1 {
		return strconv.Itoa(all[0].ForwarderSubdistrictID)
	}

	// Step 3: API call ke /subdistrictlist
	token, err := s.getForwarderToken(ctx, apiURL, clientName, username, password, "SUBDISTRICTLIST")
	if err != nil {
		log.Printf("[forwarder] resolveSubdistrict: gagal get token: %v", err)
		return "0"
	}

	reqBody, _ := json.Marshal(map[string]string{
		"subdistrict_name": strings.ToUpper(normalized),
		"city_id":          strconv.Itoa(forwarderCityID),
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/subdistrictlist", bytes.NewReader(reqBody))
	if err != nil {
		return "0"
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Client_name", clientName)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[forwarder] resolveSubdistrict: API error: %v", err)
		return "0"
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		IsSuccess string `json:"isSuccess"`
		Data      []struct {
			ItemID   int    `json:"item_id"`
			ItemName string `json:"item_name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil || result.IsSuccess != "ok" || len(result.Data) == 0 {
		log.Printf("[forwarder] resolveSubdistrict: tidak ditemukan via API untuk kecamatan=%s city_id=%d body=%s", normalized, forwarderCityID, string(respBody))
		return "0"
	}

	// Simpan ke DB agar lookup berikutnya tidak perlu API call
	s.db.Create(&models.ForwarderSubdistrictMapping{
		KecamatanPattern:         normalized,
		ForwarderCityID:          forwarderCityID,
		ForwarderSubdistrictID:   result.Data[0].ItemID,
		ForwarderSubdistrictName: result.Data[0].ItemName,
	})

	return strconv.Itoa(result.Data[0].ItemID)
}

func (s *shippingService) getForwarderToken(ctx context.Context, apiURL, clientName, username, password, scope string) (string, error) {
	tokenReq := forwarderTokenRequest{Scope: scope}
	body, _ := json.Marshal(tokenReq)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/accesstoken", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("client_name", clientName)
	httpReq.Header.Set("username", username)
	httpReq.Header.Set("password", password)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("connection timeout: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result forwarderTokenResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("gagal parse token response: %w", err)
	}
	if result.Status != "ok" || result.AccessToken == "" {
		return "", fmt.Errorf("gagal mendapatkan token: %s", result.Message)
	}
	return result.AccessToken, nil
}

// ─── Forwarder Invoice ────────────────────────────────────────────────────────

type ForwarderInvoiceDetail struct {
	FreightElementName string `json:"freight_element_name"`
	BasisName          string `json:"basis_name"`
	TotalIDR           string `json:"total_idr"`
	Amount             string `json:"amount"`
	Total              string `json:"total"`
	Subtotal           string `json:"subtotal"`
	InvoiceNo          string `json:"invoice_no"`
	Qty                string `json:"qty"`
	ContainerType      string `json:"container_type"`
	Currency           string `json:"currency"`
	Tax                string `json:"tax"`
	Remark             string `json:"remark"`
}

type ForwarderInvoice struct {
	BookingNo          string                   `json:"booking_no"`
	InvoiceNo          string                   `json:"invoice_no"`
	DueDate            string                   `json:"due_date"`
	InvoiceID          string                   `json:"invoice_id"`
	Currency           string                   `json:"currency"`
	Remark             string                   `json:"remark"`
	CreateDate         string                   `json:"create_date"`
	DownloadInvoiceURL string                   `json:"download_invoice_url"`
	DataDetail         []ForwarderInvoiceDetail `json:"data_detail"`
	InvoiceDate        string                   `json:"invoice_date"`
	QuotationNo        string                   `json:"quotation_no"`
	Status             string                   `json:"status"`
}

type forwarderInvoiceListRequest struct {
	UserName  string `json:"user_name"`
	BookingNo string `json:"booking_no"`
	InvoiceNo string `json:"invoice_no"`
}

type forwarderInvoiceListResponse struct {
	Msg       string             `json:"msg"`
	Data      []ForwarderInvoice `json:"data"`
	IsSuccess string             `json:"isSuccess"`
}

func (s *shippingService) GetForwarderInvoice(ctx context.Context, bookingNo string) ([]ForwarderInvoice, error) {
	apiURL := os.Getenv("FORWARDER_API_URL")
	clientName := os.Getenv("FORWARDER_CLIENT_NAME")
	username := os.Getenv("FORWARDER_USERNAME")
	password := os.Getenv("FORWARDER_PASSWORD")
	if apiURL == "" || clientName == "" || username == "" || password == "" {
		return nil, fmt.Errorf("konfigurasi Forwarder tidak lengkap")
	}

	token, err := s.getForwarderToken(ctx, apiURL, clientName, username, password, "INVOICELIST")
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan token Forwarder: %w", err)
	}

	reqBody := forwarderInvoiceListRequest{
		UserName:  clientName,
		BookingNo: bookingNo,
		InvoiceNo: "",
	}
	body, _ := json.Marshal(reqBody)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/invoicelist", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Client_name", clientName)

	log.Printf("[forwarder] --> POST /invoicelist booking_no=%s", bookingNo)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[forwarder] <-- POST /invoicelist booking_no=%s error=%v", bookingNo, err)
		return nil, fmt.Errorf("connection timeout atau gagal menghubungi Forwarder: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[forwarder] <-- POST /invoicelist booking_no=%s status=%d body=%s", bookingNo, resp.StatusCode, string(respBody))

	var result forwarderInvoiceListResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response Forwarder: %w", err)
	}
	if result.IsSuccess != "ok" {
		return nil, fmt.Errorf("Forwarder API error: %s", result.Msg)
	}

	return result.Data, nil
}

// ─── Tracking ─────────────────────────────────────────────────────────────────

func (s *shippingService) TrackDelivery(ctx context.Context, pesanan *models.Pesanan) (*TrackingResult, error) {
	switch pesanan.DeliveryType {
	case models.DeliveryTypeDeliveree:
		return s.trackDeliveree(ctx, pesanan)
	case models.DeliveryTypeForwarder, models.DeliveryTypeForwarderLCL:
		return s.trackForwarder(ctx, pesanan)
	default:
		return nil, fmt.Errorf("delivery type %s tidak mendukung tracking", pesanan.DeliveryType)
	}
}

func (s *shippingService) trackDeliveree(ctx context.Context, pesanan *models.Pesanan) (*TrackingResult, error) {
	baseURL := os.Getenv("DELIVEREE_BASE_URL")
	apiKey := os.Getenv("DELIVEREE_API_KEY")
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("konfigurasi Deliveree tidak lengkap")
	}
	if pesanan.DelivereeBookingID == nil {
		return nil, fmt.Errorf("pesanan belum memiliki Deliveree booking ID")
	}

	url := baseURL + "/deliveries/" + *pesanan.DelivereeBookingID
	log.Printf("[deliveree] --> GET /deliveries/%s pesanan=%s", *pesanan.DelivereeBookingID, pesanan.Kode)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[deliveree] <-- GET /deliveries/%s pesanan=%s error=%v", *pesanan.DelivereeBookingID, pesanan.Kode, err)
		return nil, fmt.Errorf("gagal menghubungi Deliveree: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[deliveree] <-- GET /deliveries/%s pesanan=%s status=%d body=%s", *pesanan.DelivereeBookingID, pesanan.Kode, resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Deliveree tracking error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var data struct {
		Status      string `json:"status"`
		TrackingURL string `json:"tracking_url"`
		Locations   []struct {
			Name           string `json:"name"`
			DeliveryStatus string `json:"delivery_status"`
		} `json:"locations"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, fmt.Errorf("gagal parse response Deliveree: %w", err)
	}

	history := make([]TrackingEvent, 0, len(data.Locations))
	for _, loc := range data.Locations {
		history = append(history, TrackingEvent{
			Status: loc.DeliveryStatus + " - " + loc.Name,
		})
	}

	return &TrackingResult{
		Provider:    "DELIVEREE",
		BookingRef:  *pesanan.DelivereeBookingID,
		Status:      data.Status,
		TrackingURL: data.TrackingURL,
		History:     history,
	}, nil
}

func (s *shippingService) GetDelivereeDetail(ctx context.Context, pesanan *models.Pesanan) (*DelivereeDeliveryDetail, error) {
	baseURL := os.Getenv("DELIVEREE_BASE_URL")
	apiKey := os.Getenv("DELIVEREE_API_KEY")
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf("konfigurasi Deliveree tidak lengkap")
	}
	if pesanan.DelivereeBookingID == nil {
		return nil, fmt.Errorf("pesanan belum memiliki Deliveree booking ID")
	}

	url := baseURL + "/deliveries/" + *pesanan.DelivereeBookingID
	log.Printf("[deliveree] --> GET /deliveries/%s (detail) pesanan=%s", *pesanan.DelivereeBookingID, pesanan.Kode)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[deliveree] <-- GET /deliveries/%s (detail) pesanan=%s error=%v", *pesanan.DelivereeBookingID, pesanan.Kode, err)
		return nil, fmt.Errorf("gagal menghubungi Deliveree: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[deliveree] <-- GET /deliveries/%s (detail) pesanan=%s status=%d", *pesanan.DelivereeBookingID, pesanan.Kode, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Deliveree API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var detail DelivereeDeliveryDetail
	if err := json.Unmarshal(respBody, &detail); err != nil {
		return nil, fmt.Errorf("gagal parse response Deliveree: %w", err)
	}

	return &detail, nil
}

func (s *shippingService) trackForwarder(ctx context.Context, pesanan *models.Pesanan) (*TrackingResult, error) {
	apiURL := os.Getenv("FORWARDER_API_URL")
	clientName := os.Getenv("FORWARDER_CLIENT_NAME")
	username := os.Getenv("FORWARDER_USERNAME")
	password := os.Getenv("FORWARDER_PASSWORD")
	if apiURL == "" || clientName == "" || username == "" || password == "" {
		return nil, fmt.Errorf("konfigurasi Forwarder tidak lengkap")
	}
	if pesanan.ForwarderTrackingNo == nil {
		return nil, fmt.Errorf("pesanan belum memiliki Forwarder tracking number")
	}

	token, err := s.getForwarderToken(ctx, apiURL, clientName, username, password, "TRACKANDTRACE")
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan token Forwarder: %w", err)
	}

	reqBody, _ := json.Marshal(map[string]string{
		"ref_cust_id": "",
		"booking_no":  *pesanan.ForwarderTrackingNo,
	})

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/trackandtrace", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Client_name", clientName)

	log.Printf("[forwarder] --> POST /trackandtrace pesanan=%s booking_no=%s", pesanan.Kode, *pesanan.ForwarderTrackingNo)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("[forwarder] <-- POST /trackandtrace pesanan=%s error=%v", pesanan.Kode, err)
		return nil, fmt.Errorf("gagal menghubungi Forwarder: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[forwarder] <-- POST /trackandtrace pesanan=%s status=%d body=%s", pesanan.Kode, resp.StatusCode, string(respBody))

	var result struct {
		Msg       string `json:"msg"`
		IsSuccess string `json:"isSuccess"`
		Data      []struct {
			BookingNumber string `json:"booking_number"`
			Status        []struct {
				StatusDate string `json:"status_date"`
				StatusName string `json:"status_name"`
				StatusTime string `json:"status_time"`
			} `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response Forwarder: %w", err)
	}
	if result.IsSuccess != "ok" {
		return nil, fmt.Errorf("Forwarder tracking error: %s", result.Msg)
	}

	history := []TrackingEvent{}
	currentStatus := ""
	if len(result.Data) > 0 {
		for _, s := range result.Data[0].Status {
			history = append(history, TrackingEvent{
				Date:   s.StatusDate,
				Time:   s.StatusTime,
				Status: s.StatusName,
			})
		}
		if len(history) > 0 {
			currentStatus = history[len(history)-1].Status
		}
	}

	return &TrackingResult{
		Provider:   "FORWARDER",
		BookingRef: *pesanan.ForwarderTrackingNo,
		Status:     currentStatus,
		History:    history,
	}, nil
}

func (s *shippingService) getActiveWarehouse() (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if err := s.db.Where("is_active = true").Order("created_at ASC").First(&warehouse).Error; err != nil {
		return nil, err
	}
	return &warehouse, nil
}
