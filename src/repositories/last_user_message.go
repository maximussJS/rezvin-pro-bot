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

type ILastUserMessageRepository interface {
	Create(ctx context.Context, exercise models.LastUserMessage) int64
	GetByChatId(ctx context.Context, id int64) *models.LastUserMessage
	UpdateByChatId(ctx context.Context, id int64, exercise models.LastUserMessage)
	DeleteByChatId(ctx context.Context, id int64)
}

type lastUserMessageRepositoryDependencies struct {
	dig.In

	Database db.IDatabase   `name:"Database"`
	Config   config.IConfig `name:"Config"`
}

type lastUserMessageRepository struct {
	db *gorm.DB
}

func NewLastUserMessageRepository(deps lastUserMessageRepositoryDependencies) *lastUserMessageRepository {
	r := &lastUserMessageRepository{
		db: deps.Database.GetInstance(),
	}

	if deps.Config.RunMigrations() {
		err := r.db.AutoMigrate(&models.LastUserMessage{})

		utils.PanicIfError(err)
	}

	return r
}

func (r *lastUserMessageRepository) Create(ctx context.Context, msg models.LastUserMessage) int64 {
	err := r.db.WithContext(ctx).Create(&msg).Error

	utils.PanicIfNotContextError(err)

	return msg.ChatId
}

func (r *lastUserMessageRepository) GetByChatId(ctx context.Context, id int64) *models.LastUserMessage {
	var msg models.LastUserMessage
	err := r.db.WithContext(ctx).Where("chat_id = ?", id).First(&msg).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotContextError(err)

	return &msg
}

func (r *lastUserMessageRepository) UpdateByChatId(ctx context.Context, id int64, msg models.LastUserMessage) {
	err := r.db.WithContext(ctx).Where("chat_id = ?", id).Updates(&msg).Error

	utils.PanicIfNotContextError(err)
}

func (r *lastUserMessageRepository) DeleteByChatId(ctx context.Context, id int64) {
	err := r.db.WithContext(ctx).Where("chat_id = ?", id).Delete(&models.LastUserMessage{}).Error

	utils.PanicIfNotContextError(err)
}
