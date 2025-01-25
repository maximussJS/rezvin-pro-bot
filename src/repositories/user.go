package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rezvin-pro-bot/src/config"
	"rezvin-pro-bot/src/models"
	utils2 "rezvin-pro-bot/src/utils"
)

type userRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type IUserRepository interface {
	Create(ctx context.Context, user models.User) int64
	GetById(ctx context.Context, id int64) *models.User
	GetAdminUsers(ctx context.Context) []models.User
	CountClients(ctx context.Context) int64
	GetClients(ctx context.Context, limit, offset int) []models.User
	CountPendingUsers(ctx context.Context) int64
	GetPendingUsers(ctx context.Context, limit, offset int) []models.User
	UpdateById(ctx context.Context, id int64, user models.User)
	DeleteById(ctx context.Context, id int64)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(deps userRepositoryDependencies) *userRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models.User{})

		utils2.PanicIfError(err)
	}

	return &userRepository{
		db: deps.DB,
	}
}

func (r *userRepository) CountPendingUsers(ctx context.Context) int64 {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("is_approved = ?", false).
		Where("is_admin = ?", false).
		Where("is_declined = ?", false).
		Count(&count).
		Error

	utils2.PanicIfNotContextError(err)

	return count
}

func (r *userRepository) GetPendingUsers(ctx context.Context, limit, offset int) []models.User {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("is_approved = ?", false).
		Where("is_admin = ?", false).
		Where("is_declined = ?", false).
		Limit(limit).
		Offset(offset).
		Find(&users).
		Error

	utils2.PanicIfNotContextError(err)

	return users
}

func (r *userRepository) CountClients(ctx context.Context) int64 {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("is_approved = ?", true).
		Where("is_admin = ?", false).
		Where("is_declined = ?", false).
		Count(&count).
		Error

	utils2.PanicIfNotContextError(err)

	return count
}

func (r *userRepository) GetClients(ctx context.Context, limit, offset int) []models.User {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("is_approved = ?", true).
		Where("is_admin = ?", false).
		Where("is_declined = ?", false).
		Limit(limit).
		Offset(offset).
		Find(&users).
		Error

	utils2.PanicIfNotContextError(err)

	return users
}

func (r *userRepository) GetAdminUsers(ctx context.Context) []models.User {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("is_admin = ?", true).
		Find(&users).
		Error

	utils2.PanicIfNotContextError(err)

	return users
}

func (r *userRepository) Create(ctx context.Context, user models.User) int64 {
	err := r.db.WithContext(ctx).Create(&user).Error

	utils2.PanicIfNotContextError(err)

	return user.Id
}

func (r *userRepository) GetById(ctx context.Context, id int64) *models.User {
	var user models.User
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).First(&user).Error

	if err != nil && utils2.IsRecordNotFoundError(err) {
		return nil
	}

	utils2.PanicIfNotRecordNotFound(err)

	return &user
}

func (r *userRepository) UpdateById(ctx context.Context, id int64, user models.User) {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(&user).Error

	utils2.PanicIfNotContextError(err)
}

func (r *userRepository) DeleteById(ctx context.Context, id int64) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error

	utils2.PanicIfNotContextError(err)
}