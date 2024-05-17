package repository

import (
	"context"
	"github.com/Data-service/internal/model"
)

type Repository interface {
	SaveItem(ctx context.Context, file model.Item) error
	GetItem(ctx context.Context, itemId int64) (model.Item, error)
	GetItemsByUserId(ctx context.Context, userId int64) ([]model.Item, error)
	GetItemsByIds(ctx context.Context, ids []int64) ([]model.Item, error)
	UpdateItem(ctx context.Context, file model.Item) error

	CreateCapsule(ctx context.Context, capsule model.Capsule) error
	GetCapsule(ctx context.Context, capsuleId int64) (model.Capsule, error)
	GetCapsulesByUser(ctx context.Context, userId int64) ([]model.Capsule, error)
	DeleteCapsule(ctx context.Context, capsuleId int64) error
	UpdateCapsuleName(ctx context.Context, capsule model.Capsule) error
	DeleteItem(ctx context.Context, itemId int64) error

	AddItemToCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error
	DeleteItemFromCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error
	GetCapsuleItems(ctx context.Context, capsuleId int64) ([]int64, error)

	GetLook(ctx context.Context, lookId int64) (model.Look, error)
	GetItemsArray(ctx context.Context, itemIds []int64) ([]model.Item, error)
	GetLookItems(ctx context.Context, lookId int64) ([]model.LooksItems, error)
	DeleteLookItems(ctx context.Context, lookId int64) error
	DeleteLook(ctx context.Context, lookId int64) error
	DeleteItemFromLook(ctx context.Context, lookId, itemId int64) error
	AddItemToLook(ctx context.Context, lookId, itemId int64) error
	CreateLook(ctx context.Context, look model.Look) error
}
