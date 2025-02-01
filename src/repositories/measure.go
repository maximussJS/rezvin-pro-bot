package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rezvin-pro-bot/src/config"
	"rezvin-pro-bot/src/internal/db"
	"rezvin-pro-bot/src/models"
	"rezvin-pro-bot/src/utils"
)

type IMeasureRepository interface {
	Create(ctx context.Context, measure models.Measure) uint
	GetById(ctx context.Context, id uint) *models.Measure
	CountAll(ctx context.Context) int64
	GetAll(ctx context.Context, limit, offset int) []models.Measure
	GetByName(ctx context.Context, name string) *models.Measure
	CountNotAssignedToUser(ctx context.Context, userId int64) int64
	GetNotAssignedToUser(ctx context.Context, userId int64, limit, offset int) []models.Measure
	UpdateById(ctx context.Context, id uint, measure models.Measure)
	DeleteById(ctx context.Context, id uint)
}

type measureRepositoryDependencies struct {
	dig.In

	Database db.IDatabase   `name:"Database"`
	Config   config.IConfig `name:"Config"`
}

type measureRepository struct {
	db *gorm.DB
}

func NewMeasureRepository(deps measureRepositoryDependencies) *measureRepository {
	r := &measureRepository{
		db: deps.Database.GetInstance(),
	}

	if deps.Config.RunMigrations() {
		err := r.db.AutoMigrate(&models.Measure{})

		utils.PanicIfError(err)
	}

	return r
}

func (r *measureRepository) CountNotAssignedToUser(ctx context.Context, userId int64) int64 {
	var count int64

	subQuery := r.db.WithContext(ctx).Model(&models.UserMeasure{}).Select("measure_id").Where("user_id = ?", userId)

	err := r.db.WithContext(ctx).Model(&models.Measure{}).Where("id NOT IN (?)", subQuery).Count(&count).Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *measureRepository) GetNotAssignedToUser(ctx context.Context, userId int64, limit, offset int) []models.Measure {
	var measures []models.Measure

	subQuery := r.db.WithContext(ctx).Model(&models.UserMeasure{}).Select("measure_id").Where("user_id = ?", userId)

	err := r.db.
		WithContext(ctx).
		Model(&models.Measure{}).
		Where("id NOT IN (?)", subQuery).
		Limit(limit).
		Offset(offset).
		Find(&measures).
		Error

	utils.PanicIfError(err)

	return measures
}

func (r *measureRepository) Create(ctx context.Context, measure models.Measure) uint {
	err := r.db.WithContext(ctx).Create(&measure).Error

	utils.PanicIfNotContextError(err)

	return measure.Id
}

func (r *measureRepository) CountAll(ctx context.Context) int64 {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Measure{}).Count(&count).Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *measureRepository) GetById(ctx context.Context, id uint) *models.Measure {
	var measure models.Measure
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).First(&measure).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &measure
}

func (r *measureRepository) GetAll(ctx context.Context, limit, offset int) []models.Measure {
	var measures []models.Measure

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&measures).Error

	utils.PanicIfNotContextError(err)

	return measures
}

func (r *measureRepository) GetByName(ctx context.Context, name string) *models.Measure {
	var measure models.Measure
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&measure).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &measure
}

func (r *measureRepository) UpdateById(ctx context.Context, id uint, measure models.Measure) {
	err := r.db.WithContext(ctx).Model(&models.Measure{}).Where("id = ?", id).Updates(&measure).Error

	utils.PanicIfNotContextError(err)
}

func (r *measureRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Measure{}).Error

	utils.PanicIfNotContextError(err)
}
