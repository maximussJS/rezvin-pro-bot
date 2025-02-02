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
	GetAllByUserIdAndMeasureId(ctx context.Context, userId int64, measureId uint) []models.UserMeasure
	DeleteById(ctx context.Context, id uint)
	DeleteByMeasureId(ctx context.Context, measureId uint)
	GetLastByUserIdAndMeasureId(ctx context.Context, userId int64, measureId uint) *models.UserMeasure
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

func (r *userMeasureRepository) GetAllByUserIdAndMeasureId(ctx context.Context, userId int64, measureId uint) []models.UserMeasure {
	var records []models.UserMeasure

	err := r.db.WithContext(ctx).
		Preload("Measure").
		Where("user_id = ? AND measure_id = ?", userId, measureId).
		Order("created_at asc").
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userMeasureRepository) GetLastByUserIdAndMeasureId(ctx context.Context, userId int64, measureId uint) *models.UserMeasure {
	var record models.UserMeasure

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND measure_id = ?", userId, measureId).
		Order("created_at desc").
		First(&record).
		Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	return &record
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

func (r *userMeasureRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.UserMeasure{}).
		Error

	utils.PanicIfNotContextError(err)
}

func (r *userMeasureRepository) DeleteByMeasureId(ctx context.Context, measureId uint) {
	err := r.db.WithContext(ctx).
		Where("measure_id = ?", measureId).
		Delete(&models.UserMeasure{}).
		Error

	utils.PanicIfNotContextError(err)
}
