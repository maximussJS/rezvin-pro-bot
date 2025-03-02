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

type IExerciseRepository interface {
	Create(ctx context.Context, exercise models.Exercise) uint
	CountByProgramId(ctx context.Context, programId uint) int64
	GetByProgramId(ctx context.Context, programId uint, limit, offset int) []models.Exercise
	GetById(ctx context.Context, id uint) *models.Exercise
	GetByIdAndProgramId(ctx context.Context, id, programId uint) *models.Exercise
	GetAll(ctx context.Context, limit, offset int) []models.Exercise
	GetByNameAndProgramId(ctx context.Context, name string, programId uint) *models.Exercise
	UpdateById(ctx context.Context, id uint, exercise models.Exercise)
	DeleteById(ctx context.Context, id uint)
}

type exerciseRepositoryDependencies struct {
	dig.In

	Database db.IDatabase   `name:"Database"`
	Config   config.IConfig `name:"Config"`
}

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(deps exerciseRepositoryDependencies) *exerciseRepository {
	r := &exerciseRepository{
		db: deps.Database.GetInstance(),
	}

	if deps.Config.RunMigrations() {
		err := r.db.AutoMigrate(&models.Exercise{})

		utils.PanicIfError(err)
	}

	return r
}

func (r *exerciseRepository) Create(ctx context.Context, exercise models.Exercise) uint {
	err := r.db.WithContext(ctx).Create(&exercise).Error

	utils.PanicIfNotContextError(err)

	return exercise.Id
}

func (r *exerciseRepository) GetById(ctx context.Context, id uint) *models.Exercise {
	var exercise models.Exercise
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).First(&exercise).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &exercise
}

func (r *exerciseRepository) GetAllByProgramId(ctx context.Context, programId uint) []models.Exercise {
	var exercises []models.Exercise

	err := r.db.WithContext(ctx).Where("program_id = ?", programId).Find(&exercises).Error

	utils.PanicIfNotContextError(err)

	return exercises
}

func (r *exerciseRepository) GetByProgramId(ctx context.Context, programId uint, limit, offset int) []models.Exercise {
	var exercises []models.Exercise

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Where("program_id = ?", programId).Find(&exercises).Error

	utils.PanicIfNotContextError(err)

	return exercises
}

func (r *exerciseRepository) CountByProgramId(ctx context.Context, programId uint) int64 {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Exercise{}).Where("program_id = ?", programId).Count(&count).Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *exerciseRepository) GetByIdAndProgramId(ctx context.Context, id, programId uint) *models.Exercise {
	var exercise models.Exercise
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).Where("program_id = ?", programId).First(&exercise).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &exercise
}

func (r *exerciseRepository) GetAll(ctx context.Context, limit, offset int) []models.Exercise {
	var exercises []models.Exercise

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&exercises).Error

	utils.PanicIfNotContextError(err)

	return exercises
}

func (r *exerciseRepository) GetByNameAndProgramId(ctx context.Context, name string, programId uint) *models.Exercise {
	var exercise models.Exercise
	err := r.db.WithContext(ctx).Where("name = ?", name).Where("program_id = ?", programId).First(&exercise).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &exercise
}

func (r *exerciseRepository) UpdateById(ctx context.Context, id uint, exercise models.Exercise) {
	err := r.db.WithContext(ctx).Model(&models.Exercise{}).Where("id = ?", id).Updates(&exercise).Error

	utils.PanicIfNotContextError(err)
}

func (r *exerciseRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Exercise{}).Error

	utils.PanicIfNotContextError(err)
}
