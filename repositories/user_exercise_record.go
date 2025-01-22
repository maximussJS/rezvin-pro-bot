package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"rezvin-pro-bot/config"
	"rezvin-pro-bot/models"
	"rezvin-pro-bot/utils"
)

type userExerciseRecordRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type IUserExerciseRecordRepository interface {
	Create(ctx context.Context, record models.UserExerciseRecord)
	CreateMany(ctx context.Context, records []models.UserExerciseRecord)
	UpdateByUserIdAndExerciseId(ctx context.Context, userId int64, exerciseId uint, record models.UserExerciseRecord)
}

type userExerciseRecordRepository struct {
	db *gorm.DB
}

func NewUserExerciseRecordRepository(deps userExerciseRecordRepositoryDependencies) *userExerciseRecordRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models.UserExerciseRecord{})

		utils.PanicIfError(err)
	}

	return &userExerciseRecordRepository{
		db: deps.DB,
	}
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
