package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	Create(ctx context.Context, req *entity.Settings) error
	FindByName(ctx context.Context, name string) (*entity.Settings, error)
	Update(ctx context.Context, settings *entity.Settings) error
	Delete(ctx context.Context, settings *entity.Settings) error
	FindAll(ctx context.Context, q dto.SettingsListQuery) (items []*entity.Settings, meta pagination.Meta, err error)
	FindByID(ctx context.Context, id int) (*entity.Settings, error)
}

type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{
		db: db,
	}
}

// Create implements SettingsRepository.
func (s *settingsRepository) Create(ctx context.Context, req *entity.Settings) error {
	return s.db.WithContext(ctx).Create(req).Error
}

// Delete implements SettingsRepository.
func (s *settingsRepository) Delete(ctx context.Context, settings *entity.Settings) error {
	return s.db.WithContext(ctx).Delete(settings).Error
}

// FindAll implements SettingsRepository.
func (s *settingsRepository) FindAll(ctx context.Context, q dto.SettingsListQuery) (items []*entity.Settings, meta pagination.Meta, err error) {
	q.Normalize() // Terapkan DefaultPage dan DefaultLimit

	base := s.db.WithContext(ctx).Model(&entity.Settings{})

	// Tentukan kolom yang boleh di-sort
	allowedSort := map[string]struct{}{"created_at": {}, "name": {}, "id": {}}

	// Terapkan filter dan search
	filtered := base.
		Scopes(
			SettingsFilters(q), // Filter by name
			// Search by name dan value
			ILike([]string{"settings.name", "settings.value"}, q.Q),
		)

	// Hitung total data
	var total int64
	if err = filtered.Count(&total).Error; err != nil {
		return
	}

	// Ambil data untuk halaman ini
	if err = filtered.
		Scopes(
			func(db *gorm.DB) *gorm.DB { return pagination.ScopeSort(db, q.Sort, allowedSort) },
			func(db *gorm.DB) *gorm.DB { return pagination.ScopePaginate(db, q.Page, q.Limit) },
		).
		Find(&items).Error; err != nil {
		return
	}

	// Hitung metadata paginasi
	meta = pagination.CalcMeta(int(total), q.Page, q.Limit)
	return
}

// FindByID implements SettingsRepository.
func (s *settingsRepository) FindByID(ctx context.Context, id int) (*entity.Settings, error) {
	var settings entity.Settings
	if err := s.db.WithContext(ctx).First(&settings, id).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

// FindByName implements SettingsRepository.
func (s *settingsRepository) FindByName(ctx context.Context, name string) (*entity.Settings, error) {
	var settings entity.Settings
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&settings).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

// Update implements SettingsRepository.
func (s *settingsRepository) Update(ctx context.Context, settings *entity.Settings) error {
	return s.db.WithContext(ctx).Save(settings).Error
}

func SettingsFilters(q dto.SettingsListQuery) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.Name != nil {
			// Gunakan ILIKE agar lebih fleksibel
			db = db.Where("name ILIKE ?", "%"+*q.Name+"%")
		}
		return db
	}
}
