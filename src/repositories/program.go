package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rezvin-pro-bot/src/config"
	models2 "rezvin-pro-bot/src/models"
	utils2 "rezvin-pro-bot/src/utils"
)

type IProgramRepository interface {
	Create(ctx context.Context, program models2.Program) uint
	GetById(ctx context.Context, id uint) *models2.Program
	CountAll(ctx context.Context) int64
	GetAll(ctx context.Context, limit, offset int) []models2.Program
	GetByName(ctx context.Context, name string) *models2.Program
	CountNotAssignedToUser(ctx context.Context, userId int64) int64
	GetNotAssignedToUser(ctx context.Context, userId int64, limit, offset int) []models2.Program
	UpdateById(ctx context.Context, id uint, program models2.Program)
	DeleteById(ctx context.Context, id uint)
}

type programRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type programRepository struct {
	db *gorm.DB
}

func NewProgramRepository(deps programRepositoryDependencies) *programRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models2.Program{})

		utils2.PanicIfError(err)
	}

	return &programRepository{
		db: deps.DB,
	}
}

func (r *programRepository) CountNotAssignedToUser(ctx context.Context, userId int64) int64 {
	var count int64

	subQuery := r.db.WithContext(ctx).Model(&models2.UserProgram{}).Select("program_id").Where("user_id = ?", userId)

	err := r.db.WithContext(ctx).Model(&models2.Program{}).Where("id NOT IN (?)", subQuery).Count(&count).Error

	utils2.PanicIfNotContextError(err)

	return count
}

func (r *programRepository) GetNotAssignedToUser(ctx context.Context, userId int64, limit, offset int) []models2.Program {
	var programs []models2.Program

	subQuery := r.db.WithContext(ctx).Model(&models2.UserProgram{}).Select("program_id").Where("user_id = ?", userId)

	err := r.db.
		WithContext(ctx).
		Model(&models2.Program{}).
		Where("id NOT IN (?)", subQuery).
		Limit(limit).
		Offset(offset).
		Find(&programs).
		Error

	utils2.PanicIfError(err)

	return programs
}

func (r *programRepository) Create(ctx context.Context, program models2.Program) uint {
	err := r.db.WithContext(ctx).Create(&program).Error

	utils2.PanicIfNotContextError(err)

	return program.Id
}

func (r *programRepository) CountAll(ctx context.Context) int64 {
	var count int64

	err := r.db.WithContext(ctx).Model(&models2.Program{}).Count(&count).Error

	utils2.PanicIfNotContextError(err)

	return count
}

func (r *programRepository) GetById(ctx context.Context, id uint) *models2.Program {
	var program models2.Program
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Preload("Exercises").Where("id = ?", id).First(&program).Error

	if err != nil && utils2.IsRecordNotFoundError(err) {
		return nil
	}

	utils2.PanicIfNotRecordNotFound(err)

	return &program
}

func (r *programRepository) GetAll(ctx context.Context, limit, offset int) []models2.Program {
	var programs []models2.Program

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&programs).Error

	utils2.PanicIfNotContextError(err)

	return programs
}

func (r *programRepository) GetByName(ctx context.Context, name string) *models2.Program {
	var program models2.Program
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&program).Error

	if err != nil && utils2.IsRecordNotFoundError(err) {
		return nil
	}

	utils2.PanicIfNotRecordNotFound(err)

	return &program
}

func (r *programRepository) UpdateById(ctx context.Context, id uint, program models2.Program) {
	err := r.db.WithContext(ctx).Model(&models2.Program{}).Where("id = ?", id).Updates(&program).Error

	utils2.PanicIfNotContextError(err)
}

func (r *programRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models2.Program{}).Error

	utils2.PanicIfNotContextError(err)
}
