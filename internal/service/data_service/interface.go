package data_service

import (
	"context"
	"github.com/Data-service/internal/model"
	"io"
)

type Service interface {
	CreateFile(ctx context.Context, item model.Item, fileName string, bucketName string, fileSize int64, reader io.Reader) error
	GetItem(ctx context.Context, itemId int64) (model.Item, error)
	GetItemsByUser(ctx context.Context, userId int64) ([]model.Item, error)
	DeleteItem(ctx context.Context, itemId int64) error
	UpdateItem(ctx context.Context, item model.Item) error

	ClearFile(fileName string, reader io.Reader) (io.Reader, int64, error)

	CreateCapsule(ctx context.Context, capsule model.Capsule) error
	UpdateCapsule(ctx context.Context, capsule model.Capsule) error
	DeleteCapsule(ctx context.Context, capsuleId int64) error
	GetCapsule(ctx context.Context, capsuleId int64) (model.Capsule, error)
	GetCapsulesByUser(ctx context.Context, userId int64) ([]model.Capsule, error)
	AddItemToCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error
	DeleteItemFromCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error
	GetItemsInCapsule(ctx context.Context, capsuleId int64) ([]model.Item, error)

	GetFilesByUser(ctx context.Context, userId string) ([]*model.File, error)
	DeleteFile(ctx context.Context, bucketName, fileName string) error
	GetFilesArray(ctx context.Context, bucketName string, fileNames map[string]struct{}) ([]*model.File, error)

	GetLookData(ctx context.Context, lookId int64) (model.LookData, error)
	AddToLook(ctx context.Context, lookId, itemId int64) error
	DeleteLook(ctx context.Context, lookId int64) error
	CreateLook(ctx context.Context, look model.Look) error
}
