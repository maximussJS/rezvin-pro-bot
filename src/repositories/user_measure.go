package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"rezvin-pro-bot/src/config"
	"rezvin-pro-bot/src/internal/db"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/utils"
)

type userMeasureRepositoryDependencies struct {
	dig.In

	Database db.IDatabase   `name:"Database"`
	Config   config.IConfig `name:"Config"`
}

type IUserMeasureRepository interface {
	Create(ctx context.Context, record models.UserMeasure)
	GetById(ctx context.Context, id uint) *models.UserMeasure
	CountAllByMeasureId(ctx context.Context, measureId uint) int64
	GetAllByMeasureId(ctx context.Context, measureId uint) []models.UserMeasure
	GetByMeasureId(ctx context.Context, measureId uint, limit, offset int) []models.UserMeasure
	CountAllByUserId(ctx context.Context, userId int64) int64
	GetByUserId(ctx context.Context, userId int64, limit, offset int) []models.UserMeasure
	UpdateById(ctx context.Context, id uint, record models.UserMeasure)
	DeleteByMeasureId(ctx context.Context, measureId uint)
}

type userMeasureRepository struct {
	db *gorm.DB
}

func NewUserMeasureRepository(deps userMeasureRepositoryDependencies) *userMeasureRepository {
	r := &userMeasureRepository{
		db: deps.Database.GetInstance(),
	}

	if deps.Config.RunMigrations() {
		err := r.db.AutoMigrate(&models.UserMeasure{})

		utils.PanicIfError(err)
	}

	return r
}

func (r *userMeasureRepository) Create(ctx context.Context, record models.UserMeasure) {
	err := r.db.WithContext(ctx).Create(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userMeasureRepository) CountAllByMeasureId(ctx context.Context, measureId uint) int64 {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserMeasure{}).
		Where("measure_id = ?", measureId).
		Count(&count).
		Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *userMeasureRepository) CountAllByUserId(ctx context.Context, userId int64) int64 {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserMeasure{}).
		Where("user_id = ?", userId).
		Count(&count).
		Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *userMeasureRepository) GetById(ctx context.Context, id uint) *models.UserMeasure {
	var record models.UserMeasure

	err := r.db.WithContext(ctx).
		Preload("Measure").
		Where("id = ?", id).
		First(&record).
		Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	return &record
}

func (r *userMeasureRepository) GetAllByMeasureId(ctx context.Context, measureId uint) []models.UserMeasure {
	var records []models.UserMeasure

	err := r.db.WithContext(ctx).
		Where("measure_id = ?", measureId).
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userMeasureRepository) GetByMeasureId(
	ctx context.Context,
	measureId, limit, offset int,
) []models.UserMeasure {
	var records []models.UserMeasure

	err := r.db.WithContext(ctx).
		Where("measure_id = ?", measureId).
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userMeasureRepository) GetByUserId(
	ctx context.Context,
	userId int64,
	limit, offset int,
) []models.UserMeasure {
	var records []models.UserMeasure

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userMeasureRepository) UpdateById(ctx context.Context, id uint, record models.UserMeasure) {
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userMeasureRepository) DeleteByMeasureId(ctx context.Context, measureId uint) {
	err := r.db.WithContext(ctx).
		Where("measure_id = ?", measureId).
		Delete(&models.UserMeasure{}).
		Error

	utils.PanicIfNotContextError(err)
}
