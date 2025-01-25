package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"rezvin-pro-bot/src/config"
	"rezvin-pro-bot/src/models"
	utils2 "rezvin-pro-bot/src/utils"
)

type userProgramRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type IUserProgramRepository interface {
	Create(ctx context.Context, userProgram models.UserProgram) uint
	GetById(ctx context.Context, id uint) *models.UserProgram
	GetAllByProgramId(ctx context.Context, programId uint) []models.UserProgram
	GetByUserIdAndProgramId(ctx context.Context, userId int64, programId uint) *models.UserProgram
	CountAllByUserId(ctx context.Context, userId int64) int64
	GetByUserId(ctx context.Context, userId int64, limit, offset int) []models.UserProgram
	DeleteById(ctx context.Context, id uint)
	DeleteByUserIdAndProgramId(ctx context.Context, userId int64, programId uint)
}

type userProgramRepository struct {
	db *gorm.DB
}

func NewUserProgramRepository(deps userProgramRepositoryDependencies) *userProgramRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models.UserProgram{})

		utils2.PanicIfError(err)
	}

	return &userProgramRepository{
		db: deps.DB,
	}
}

func (r *userProgramRepository) Create(ctx context.Context, userProgram models.UserProgram) uint {
	err := r.db.WithContext(ctx).Create(&userProgram).Error

	utils2.PanicIfNotContextError(err)

	return userProgram.Id
}

func (r *userProgramRepository) GetById(ctx context.Context, id uint) *models.UserProgram {
	var userProgram models.UserProgram

	err := r.db.WithContext(ctx).
		Preload("Program").
		Where("id = ?", id).
		First(&userProgram).
		Error

	if err != nil && utils2.IsRecordNotFoundError(err) {
		return nil
	}

	utils2.PanicIfNotRecordNotFound(err)

	return &userProgram
}

func (r *userProgramRepository) GetAllByProgramId(ctx context.Context, programId uint) []models.UserProgram {
	var userPrograms []models.UserProgram

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("program_id = ?", programId).
		Find(&userPrograms).
		Error

	utils2.PanicIfNotContextError(err)

	return userPrograms
}

func (r *userProgramRepository) CountAllByUserId(ctx context.Context, userId int64) int64 {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserProgram{}).
		Where("user_id = ?", userId).
		Count(&count).
		Error

	utils2.PanicIfNotContextError(err)

	return count
}

func (r *userProgramRepository) GetByUserIdAndProgramId(ctx context.Context, userId int64, programId uint) *models.UserProgram {
	var userProgram models.UserProgram

	err := r.db.WithContext(ctx).
		Preload("Program").
		Where("user_id = ?", userId).
		Where("program_id = ?", programId).
		First(&userProgram).
		Error

	if err != nil && utils2.IsRecordNotFoundError(err) {
		return nil
	}

	utils2.PanicIfNotRecordNotFound(err)

	return &userProgram
}

func (r *userProgramRepository) GetByUserId(ctx context.Context, userId int64, limit, offset int) []models.UserProgram {
	var userPrograms []models.UserProgram

	err := r.db.WithContext(ctx).
		Preload("Program").
		Where("user_id = ?", userId).
		Limit(limit).
		Offset(offset).
		Find(&userPrograms).
		Error

	utils2.PanicIfNotContextError(err)

	return userPrograms
}

func (r *userProgramRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.UserProgram{}).
		Error

	utils2.PanicIfNotContextError(err)
}

func (r *userProgramRepository) DeleteByUserIdAndProgramId(ctx context.Context, userId int64, programId uint) {
	err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Where("program_id = ?", programId).
		Delete(&models.UserProgram{}).
		Error

	utils2.PanicIfNotContextError(err)
}