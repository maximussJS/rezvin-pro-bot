package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rezvin-pro-bot/config"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/utils"
)

type IExerciseRepository interface {
	Create(ctx context.Context, exercise models.Exercise) uint
	GetById(ctx context.Context, id uint) *models.Exercise
	GetByIdAndProgramId(ctx context.Context, id, programId uint) *models.Exercise
	GetAll(ctx context.Context, limit, offset int) []models.Exercise
	GetByNameAndProgramId(ctx context.Context, name string, programId uint) *models.Exercise
	UpdateById(ctx context.Context, id uint, exercise models.Exercise)
	DeleteById(ctx context.Context, id uint)
}

type exerciseRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(deps exerciseRepositoryDependencies) *exerciseRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models.Exercise{})

		utils.PanicIfError(err)
	}

	return &exerciseRepository{
		db: deps.DB,
	}
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