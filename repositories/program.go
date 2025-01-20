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

type IProgramRepository interface {
	Create(ctx context.Context, program models.Program) uint
	GetById(ctx context.Context, id uint) *models.Program
	GetAll(ctx context.Context, limit, offset int) []models.Program
	GetByName(ctx context.Context, name string) *models.Program
	UpdateById(ctx context.Context, id uint, program models.Program)
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
		err := deps.DB.AutoMigrate(&models.Program{})

		utils.PanicIfError(err)
	}

	return &programRepository{
		db: deps.DB,
	}
}

func (r *programRepository) Create(ctx context.Context, program models.Program) uint {
	err := r.db.WithContext(ctx).Create(&program).Error

	utils.PanicIfNotContextError(err)

	return program.Id
}

func (r *programRepository) GetById(ctx context.Context, id uint) *models.Program {
	var program models.Program
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).First(&program).Error

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
