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

type userRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type IUserRepository interface {
	Create(ctx context.Context, user models.User) int64
	GetById(ctx context.Context, id int64) *models.User
	UpdateById(ctx context.Context, id int64, user models.User)
	DeleteById(ctx context.Context, id int64)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(deps userRepositoryDependencies) *userRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models.User{})

		utils.PanicIfError(err)
	}

	return &userRepository{
		db: deps.DB,
	}
}

func (r *userRepository) Create(ctx context.Context, user models.User) int64 {
	err := r.db.WithContext(ctx).Create(&user).Error

	utils.PanicIfNotContextError(err)

	return user.Id
}

func (r *userRepository) GetById(ctx context.Context, id int64) *models.User {
	var user models.User
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).First(&user).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &user
}

func (r *userRepository) UpdateById(ctx context.Context, id int64, user models.User) {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(&user).Error

	utils.PanicIfNotContextError(err)
}

func (r *userRepository) DeleteById(ctx context.Context, id int64) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error

	utils.PanicIfNotContextError(err)
}
