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

type IProgramRepository interface {
	Create(ctx context.Context, program models.Program) uint
	GetById(ctx context.Context, id uint) *models.Program
	CountAll(ctx context.Context) int64
	GetAll(ctx context.Context, limit, offset int) []models.Program
	GetByName(ctx context.Context, name string) *models.Program
	CountNotAssignedToUser(ctx context.Context, userId int64) int64
	GetNotAssignedToUser(ctx context.Context, userId int64, limit, offset int) []models.Program
	UpdateById(ctx context.Context, id uint, program models.Program)
	DeleteById(ctx context.Context, id uint)
}

type programRepositoryDependencies struct {
	dig.In

	Database db.IDatabase   `name:"Database"`
	Config   config.IConfig `name:"Config"`
}

type programRepository struct {
	db *gorm.DB
}

func NewProgramRepository(deps programRepositoryDependencies) *programRepository {
	r := &programRepository{
		db: deps.Database.GetInstance(),
	}

	if deps.Config.RunMigrations() {
		err := r.db.AutoMigrate(&models.Program{})

		utils.PanicIfError(err)
	}

	return r
}

func (r *programRepository) CountNotAssignedToUser(ctx context.Context, userId int64) int64 {
	var count int64

	subQuery := r.db.WithContext(ctx).Model(&models.UserProgram{}).Select("program_id").Where("user_id = ?", userId)

	err := r.db.WithContext(ctx).Model(&models.Program{}).Where("id NOT IN (?)", subQuery).Count(&count).Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *programRepository) GetNotAssignedToUser(ctx context.Context, userId int64, limit, offset int) []models.Program {
	var programs []models.Program

	subQuery := r.db.WithContext(ctx).Model(&models.UserProgram{}).Select("program_id").Where("user_id = ?", userId)

	err := r.db.
		WithContext(ctx).
		Model(&models.Program{}).
		Where("id NOT IN (?)", subQuery).
		Limit(limit).
		Offset(offset).
		Find(&programs).
		Error

	utils.PanicIfError(err)

	return programs
}

func (r *programRepository) Create(ctx context.Context, program models.Program) uint {
	err := r.db.WithContext(ctx).Create(&program).Error

	utils.PanicIfNotContextError(err)

	return program.Id
}

func (r *programRepository) CountAll(ctx context.Context) int64 {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Program{}).Count(&count).Error

	utils.PanicIfNotContextError(err)

	return count
}

func (r *programRepository) GetById(ctx context.Context, id uint) *models.Program {
	var program models.Program
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Preload("Exercises").Where("id = ?", id).First(&program).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &program
}

func (r *programRepository) GetAll(ctx context.Context, limit, offset int) []models.Program {
	var programs []models.Program

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&programs).Error

	utils.PanicIfNotContextError(err)

	return programs
}

func (r *programRepository) GetByName(ctx context.Context, name string) *models.Program {
	var program models.Program
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&program).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &program
}

func (r *programRepository) UpdateById(ctx context.Context, id uint, program models.Program) {
	err := r.db.WithContext(ctx).Model(&models.Program{}).Where("id = ?", id).Updates(&program).Error

	utils.PanicIfNotContextError(err)
}

func (r *programRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Program{}).Error

	utils.PanicIfNotContextError(err)
}
