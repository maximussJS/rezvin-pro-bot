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

type userExerciseRecordRepositoryDependencies struct {
	dig.In

	Database db.IDatabase   `name:"Database"`
	Config   config.IConfig `name:"Config"`
}

type IUserExerciseRecordRepository interface {
	Create(ctx context.Context, record models.UserExerciseRecord)
	CreateMany(ctx context.Context, records []models.UserExerciseRecord)
	GetById(ctx context.Context, id uint) *models.UserExerciseRecord
	CountAllByUserProgramId(ctx context.Context, userProgramId uint) int64
	GetAllByUserProgramIdAndExerciseId(ctx context.Context, userProgramId, exerciseId uint) []models.UserExerciseRecord
	GetAllByUserProgramId(ctx context.Context, userProgramId uint) []models.UserExerciseRecord
	GetAllByExerciseId(ctx context.Context, exerciseId uint) []models.UserExerciseRecord
	GetByUserProgramId(ctx context.Context, userProgramId uint, limit, offset int) []models.UserExerciseRecord
	UpdateById(ctx context.Context, id uint, record models.UserExerciseRecord)
	UpdateByUserIdAndExerciseId(ctx context.Context, userId int64, exerciseId uint, record models.UserExerciseRecord)
	DeleteByUserProgramId(ctx context.Context, userProgramId uint)
	DeleteByExerciseId(ctx context.Context, exerciseId uint)
}

type userExerciseRecordRepository struct {
	db *gorm.DB
}

func NewUserExerciseRecordRepository(deps userExerciseRecordRepositoryDependencies) *userExerciseRecordRepository {
	r := &userExerciseRecordRepository{
		db: deps.Database.GetInstance(),
	}

	if deps.Config.RunMigrations() {
		err := r.db.AutoMigrate(&models.UserExerciseRecord{})

		utils.PanicIfError(err)
	}

	return r
}

func (r *userExerciseRecordRepository) Create(ctx context.Context, record models.UserExerciseRecord) {
	err := r.db.WithContext(ctx).Create(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userExerciseRecordRepository) CreateMany(ctx context.Context, records []models.UserExerciseRecord) {
	if len(records) == 0 {
		return
	}

	err := r.db.WithContext(ctx).Create(&records).Error

	utils.PanicIfNotContextError(err)
}

func (r *userExerciseRecordRepository) CountAllByUserProgramId(ctx context.Context, userProgramId uint) int64 {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserExerciseRecord{}).
		Where("user_program_id = ?", userProgramId).
		Count(&count).
		Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *userExerciseRecordRepository) GetById(ctx context.Context, id uint) *models.UserExerciseRecord {
	var record models.UserExerciseRecord

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

func (r *userExerciseRecordRepository) GetAllByUserProgramId(ctx context.Context, userProgramId uint) []models.UserExerciseRecord {
	var records []models.UserExerciseRecord

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("user_program_id = ?", userProgramId).
		Order("reps ASC").
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userExerciseRecordRepository) GetAllByUserProgramIdAndExerciseId(
	ctx context.Context,
	userProgramId, exerciseId uint,
) []models.UserExerciseRecord {
	var records []models.UserExerciseRecord

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("user_program_id = ? AND exercise_id = ?", userProgramId, exerciseId).
		Order("reps ASC").
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userExerciseRecordRepository) GetAllByExerciseId(ctx context.Context, exerciseId uint) []models.UserExerciseRecord {
	var records []models.UserExerciseRecord

	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Where("exercise_id = ?", exerciseId).
		Order("reps ASC").
		Find(&records).
		Error

	utils.PanicIfNotContextError(err)

	return records
}

func (r *userExerciseRecordRepository) GetByUserProgramId(
	ctx context.Context,
	userProgramId uint,
	limit, offset int,
) []models.UserExerciseRecord {
	var records []models.UserExerciseRecord

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

func (r *userExerciseRecordRepository) UpdateById(ctx context.Context, id uint, record models.UserExerciseRecord) {
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userExerciseRecordRepository) UpdateByUserIdAndExerciseId(
	ctx context.Context, userId int64,
	exerciseId uint,
	record models.UserExerciseRecord,
) {
	err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Where("exercise_id = ?", exerciseId).
		Updates(&record).Error

	utils.PanicIfNotContextError(err)
}

func (r *userExerciseRecordRepository) DeleteByUserProgramId(ctx context.Context, userProgramId uint) {
	err := r.db.WithContext(ctx).
		Where("user_program_id = ?", userProgramId).
		Delete(&models.UserExerciseRecord{}).
		Error

	utils.PanicIfNotContextError(err)
}

func (r *userExerciseRecordRepository) DeleteByExerciseId(ctx context.Context, exerciseId uint) {
	err := r.db.WithContext(ctx).
		Where("exercise_id = ?", exerciseId).
		Delete(&models.UserExerciseRecord{}).
		Error

	utils.PanicIfNotContextError(err)
}
