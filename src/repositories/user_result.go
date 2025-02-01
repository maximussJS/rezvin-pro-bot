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

type userResultRepositoryDependencies struct {
	dig.In

	Database db.IDatabase   `name:"Database"`
	Config   config.IConfig `name:"Config"`
}

type IUserResultRepository interface {
	Create(ctx context.Context, record models.UserResult)
	CreateMany(ctx context.Context, records []models.UserResult)
	GetById(ctx context.Context, id uint) *models.UserResult
	CountAllByUserProgramId(ctx context.Context, userProgramId uint) int64
	GetAllByUserProgramIdAndExerciseId(ctx context.Context, userProgramId, exerciseId uint) []models.UserResult
	GetAllByUserProgramId(ctx context.Context, userProgramId uint) []models.UserResult
	GetAllByExerciseId(ctx context.Context, exerciseId uint) []models.UserResult
	GetByUserProgramId(ctx context.Context, userProgramId uint, limit, offset int) []models.UserResult
	UpdateById(ctx context.Context, id uint, record models.UserResult)
	UpdateByUserIdAndExerciseId(ctx context.Context, userId int64, exerciseId uint, record models.UserResult)
	DeleteByUserProgramId(ctx context.Context, userProgramId uint)
	DeleteByExerciseId(ctx context.Context, exerciseId uint)
}

type userResultRepository struct {
	db *gorm.DB
}

func NewUserResultRepository(deps userResultRepositoryDependencies) *userResultRepository {
	r := &userResultRepository{
		db: deps.Database.GetInstance(),
	}

	if deps.Config.RunMigrations() {
		err := r.db.AutoMigrate(&models.UserResult{})

		utils.PanicIfError(err)
	}

	return r
}

func (r *userResultRepository) Create(ctx context.Context, record models.UserResult) {
	err := r.db.WithContext(ctx).Create(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userResultRepository) CreateMany(ctx context.Context, records []models.UserResult) {
	if len(records) == 0 {
		return
	}

	err := r.db.WithContext(ctx).Create(&records).Error

	utils.PanicIfNotContextError(err)
}

func (r *userResultRepository) CountAllByUserProgramId(ctx context.Context, userProgramId uint) int64 {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserResult{}).
		Where("user_program_id = ?", userProgramId).
		Count(&count).
		Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *userResultRepository) GetById(ctx context.Context, id uint) *models.UserResult {
	var record models.UserResult

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("id = ?", id).
		First(&record).
		Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	return &record
}

func (r *userResultRepository) GetAllByUserProgramId(ctx context.Context, userProgramId uint) []models.UserResult {
	var records []models.UserResult

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("user_program_id = ?", userProgramId).
		Order("reps ASC").
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userResultRepository) GetAllByUserProgramIdAndExerciseId(
	ctx context.Context,
	userProgramId, exerciseId uint,
) []models.UserResult {
	var records []models.UserResult

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("user_program_id = ? AND exercise_id = ?", userProgramId, exerciseId).
		Order("reps ASC").
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userResultRepository) GetAllByExerciseId(ctx context.Context, exerciseId uint) []models.UserResult {
	var records []models.UserResult

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("exercise_id = ?", exerciseId).
		Order("reps ASC").
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userResultRepository) GetByUserProgramId(
	ctx context.Context,
	userProgramId uint,
	limit, offset int,
) []models.UserResult {
	var records []models.UserResult

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("user_program_id = ?", userProgramId).
		Order("reps ASC").
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userResultRepository) UpdateById(ctx context.Context, id uint, record models.UserResult) {
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userResultRepository) UpdateByUserIdAndExerciseId(
	ctx context.Context, userId int64,
	exerciseId uint,
	record models.UserResult,
) {
	err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Where("exercise_id = ?", exerciseId).
		Updates(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userResultRepository) DeleteByUserProgramId(ctx context.Context, userProgramId uint) {
	err := r.db.WithContext(ctx).
		Where("user_program_id = ?", userProgramId).
		Delete(&models.UserResult{}).
		Error

	utils.PanicIfNotContextError(err)
}

func (r *userResultRepository) DeleteByExerciseId(ctx context.Context, exerciseId uint) {
	err := r.db.WithContext(ctx).
		Where("exercise_id = ?", exerciseId).
		Delete(&models.UserResult{}).
		Error

	utils.PanicIfNotContextError(err)
}
