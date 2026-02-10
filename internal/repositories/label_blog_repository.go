package repositories

import (
	"context"

	"project-bulky-be/internal/dto"
	"project-bulky-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabelBlogRepository interface {
	Create(ctx context.Context, label *models.LabelBlog) error
	Update(ctx context.Context, label *models.LabelBlog) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.LabelBlog, error)
	FindBySlug(ctx context.Context, slug string) (*models.LabelBlog, error)
	FindAll(ctx context.Context, params *dto.LabelBlogFilterRequest) ([]models.LabelBlog, models.PaginationMeta, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.LabelBlog, error)
	CountBlogByLabel(ctx context.Context, LabelIDs uuid.UUID) (int64, error)
	UpdateUrutan(ctx context.Context, id uuid.UUID, urutan int) error
	FindAllPublicWithCount(ctx context.Context) ([]models.LabelBlog, error)
}

type labelBlogRepository struct {
	db *gorm.DB
}

func NewLabelBlogRepository(db *gorm.DB) LabelBlogRepository {
	return &labelBlogRepository{db: db}
}

func (r *labelBlogRepository) Create(ctx context.Context, label *models.LabelBlog) error {
	return r.db.WithContext(ctx).Create(label).Error
}

func (r *labelBlogRepository) Update(ctx context.Context, label *models.LabelBlog) error {
	return r.db.WithContext(ctx).Save(label).Error
}

func (r *labelBlogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.LabelBlog{}, id).Error
}

func (r *labelBlogRepository) CountBlogByLabel(ctx context.Context, LabelIDs uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BlogLabel{}).
		Where("label_id = ?", LabelIDs).
		Count(&count).Error
	return count, err
}

func (r *labelBlogRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.LabelBlog, error) {
	var label models.LabelBlog
	err := r.db.WithContext(ctx).First(&label, id).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *labelBlogRepository) FindBySlug(ctx context.Context, slug string) (*models.LabelBlog, error) {
	var label models.LabelBlog
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *labelBlogRepository) FindAll(ctx context.Context, params *dto.LabelBlogFilterRequest) ([]models.LabelBlog, models.PaginationMeta, error) {
	var labels []models.LabelBlog
	err := r.db.WithContext(ctx).Order("urutan ASC, nama_id ASC").Find(&labels).Error
	return labels, models.PaginationMeta{}, err
}

func (r *labelBlogRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.LabelBlog, error) {
	var labels []models.LabelBlog
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&labels).Error
	return labels, err
}

func (r *labelBlogRepository) UpdateUrutan(ctx context.Context, id uuid.UUID, urutan int) error {
	return r.db.WithContext(ctx).Model(&models.LabelBlog{}).
		Where("id = ?", id).
		Update("urutan", urutan).Error
}

func (r *labelBlogRepository) FindAllPublicWithCount(ctx context.Context) ([]models.LabelBlog, error) {
	var labels []models.LabelBlog
	err := r.db.WithContext(ctx).
		Select("label_blog.*, COUNT(blog_label.blog_id) as blog_count").
		Joins("LEFT JOIN blog_label ON blog_label.label_id = label_blog.id").
		Joins("LEFT JOIN blog ON blog.id = blog_label.blog_id AND blog.deleted_at IS NULL AND blog.is_active = true").
		Where("label_blog.deleted_at IS NULL").
		Group("label_blog.id").
		Having("COUNT(blog_label.blog_id) > 0").
		Order("label_blog.urutan ASC").
		Find(&labels).Error
	return labels, err
}
