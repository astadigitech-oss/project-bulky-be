package repositories

import (
	"context"
	"time"

	"project-bulky-be/internal/models"

	"gorm.io/gorm"
)

type BuyerRepository interface {
	FindByID(ctx context.Context, id string) (*models.Buyer, error)
	FindByIDWithAlamat(ctx context.Context, id string) (*models.Buyer, error)
	FindAll(ctx context.Context, params *models.BuyerFilterRequest) ([]models.Buyer, int64, error)
	Update(ctx context.Context, buyer *models.Buyer) error
	Delete(ctx context.Context, id string) error
	ExistsByUsername(ctx context.Context, username string, excludeID *string) (bool, error)
	ExistsByEmail(ctx context.Context, email string, excludeID *string) (bool, error)
	CountAlamat(ctx context.Context, buyerID string) (int64, error)
	GetStatistik(ctx context.Context) (*models.BuyerStatistikResponse, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	GetRegistrationByMonth(ctx context.Context, startDate, endDate time.Time) ([]models.ChartData, int64, error)
	GetRegistrationByDay(ctx context.Context, startDate, endDate time.Time) ([]models.ChartData, int64, error)
}

type buyerRepository struct {
	db *gorm.DB
}

func NewBuyerRepository(db *gorm.DB) BuyerRepository {
	return &buyerRepository{db: db}
}

func (r *buyerRepository) FindByID(ctx context.Context, id string) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&buyer).Error
	return &buyer, err
}

func (r *buyerRepository) FindByIDWithAlamat(ctx context.Context, id string) (*models.Buyer, error) {
	var buyer models.Buyer
	err := r.db.WithContext(ctx).
		Preload("Alamat", "deleted_at IS NULL").
		Where("id = ?", id).First(&buyer).Error
	return &buyer, err
}

func (r *buyerRepository) FindAll(ctx context.Context, params *models.BuyerFilterRequest) ([]models.Buyer, int64, error) {
	var buyers []models.Buyer
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Buyer{})

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("nama ILIKE ? OR username ILIKE ? OR email ILIKE ? OR telepon ILIKE ?", search, search, search, search)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// if params.IsVerified != nil {
	// 	query = query.Where("is_verified = ?", *params.IsVerified)
	// }

	query.Count(&total)

	orderClause := params.SortBy + " " + params.Order
	err := query.Order(orderClause).
		Offset(params.GetOffset()).
		Limit(params.PerPage).
		Find(&buyers).Error

	return buyers, total, err
}

func (r *buyerRepository) Update(ctx context.Context, buyer *models.Buyer) error {
	return r.db.WithContext(ctx).Save(buyer).Error
}

func (r *buyerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Buyer{}, "id = ?", id).Error
}

func (r *buyerRepository) ExistsByUsername(ctx context.Context, username string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Buyer{}).Where("username = ?", username)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *buyerRepository) ExistsByEmail(ctx context.Context, email string, excludeID *string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Buyer{}).Where("email = ?", email)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *buyerRepository) CountAlamat(ctx context.Context, buyerID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AlamatBuyer{}).Where("buyer_id = ?", buyerID).Count(&count).Error
	return count, err
}

func (r *buyerRepository) GetStatistik(ctx context.Context) (*models.BuyerStatistikResponse, error) {
	var stats models.BuyerStatistikResponse
	now := time.Now()

	// Total buyer
	r.db.WithContext(ctx).Model(&models.Buyer{}).Count(&stats.TotalBuyer)

	// Buyer verified
	r.db.WithContext(ctx).Model(&models.Buyer{}).Where("is_verified = ?", true).Count(&stats.BuyerVerified)

	// Bulan ini vs bulan lalu
	startThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startLastMonth := startThisMonth.AddDate(0, -1, 0)

	var countThisMonth, countLastMonth int64
	r.db.WithContext(ctx).Model(&models.Buyer{}).
		Where("created_at >= ? AND created_at < ?", startThisMonth, now).
		Count(&countThisMonth)
	r.db.WithContext(ctx).Model(&models.Buyer{}).
		Where("created_at >= ? AND created_at < ?", startLastMonth, startThisMonth).
		Count(&countLastMonth)

	stats.PersentaseBulanIni = calculatePersentase(countThisMonth, countLastMonth)
	stats.RegistrasiBulanIni = countThisMonth

	// Tahun ini vs tahun lalu
	startThisYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	startLastYear := startThisYear.AddDate(-1, 0, 0)

	var countThisYear, countLastYear int64
	r.db.WithContext(ctx).Model(&models.Buyer{}).
		Where("created_at >= ? AND created_at < ?", startThisYear, now).
		Count(&countThisYear)
	r.db.WithContext(ctx).Model(&models.Buyer{}).
		Where("created_at >= ? AND created_at < ?", startLastYear, startThisYear).
		Count(&countLastYear)

	stats.PersentaseTahunIni = calculatePersentase(countThisYear, countLastYear)
	stats.RegistrasiTahunIni = countThisYear

	return &stats, nil
}

func calculatePersentase(current, previous int64) models.PersentaseData {
	var value float64
	var trend string

	if previous == 0 {
		if current > 0 {
			value = 100.0
			trend = "up"
		} else {
			value = 0
			trend = "stable"
		}
	} else {
		value = float64(current-previous) / float64(previous) * 100
		if value > 0 {
			trend = "up"
		} else if value < 0 {
			trend = "down"
			value = -value // Make positive for display
		} else {
			trend = "stable"
		}
	}

	// Round to 1 decimal place
	value = float64(int(value*10)) / 10

	return models.PersentaseData{
		Value:    value,
		Trend:    trend,
		Current:  current,
		Previous: previous,
	}
}

func (r *buyerRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Buyer{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

func (r *buyerRepository) GetRegistrationByMonth(ctx context.Context, startDate, endDate time.Time) ([]models.ChartData, int64, error) {
	var results []struct {
		Month time.Time // ← Ganti dari string ke time.Time
		Count int64
	}

	err := r.db.WithContext(ctx).
		Model(&models.Buyer{}).
		Select("DATE_TRUNC('month', created_at) as month, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("DATE_TRUNC('month', created_at)").
		Order("month").
		Scan(&results).Error

	if err != nil {
		return nil, 0, err
	}

	// Create chart data with all months
	chartData := []models.ChartData{}
	var total int64

	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		count := int64(0)

		// Find count for this month - compare Year & Month
		for _, r := range results {
			if r.Month.Year() == current.Year() && r.Month.Month() == current.Month() {
				count = r.Count
				break
			}
		}

		chartData = append(chartData, models.ChartData{
			Date: current,
			User: int(count),
		})
		total += count

		// Move to next month
		current = current.AddDate(0, 1, 0)
	}

	return chartData, total, nil
}

func (r *buyerRepository) GetRegistrationByDay(ctx context.Context, startDate, endDate time.Time) ([]models.ChartData, int64, error) {
	var results []struct {
		Day   time.Time
		Count int64
	}

	err := r.db.WithContext(ctx).
		Model(&models.Buyer{}).
		Select("DATE_TRUNC('day', created_at) as day, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("DATE_TRUNC('day', created_at)").
		Order("day").
		Scan(&results).Error

	if err != nil {
		return nil, 0, err
	}

	// Create chart data with all days
	chartData := []models.ChartData{}
	var total int64

	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		count := int64(0)

		// Find count for this day - compare Year, Month & Day
		for _, r := range results {
			if r.Day.Year() == current.Year() && r.Day.Month() == current.Month() && r.Day.Day() == current.Day() {
				count = r.Count
				break
			}
		}

		chartData = append(chartData, models.ChartData{
			Date: current,
			User: int(count),
		})
		total += count

		// Move to next day
		current = current.AddDate(0, 0, 1)
	}

	return chartData, total, nil
}
