package repository

import (
	"context"
	"fmt"
	"github.com/Data-service/internal/model"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

const (
	capsulesTable      = "capsules"
	capsulesItemsTable = "capsules_items"
	itemsTable         = "items"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) SaveItem(ctx context.Context, file model.Item) error {
	fmt.Println("saving item")
	query := `INSERT INTO items (url, name, category, size_text, size_number, description, color, user_id, file_name)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	if _, err := r.db.ExecContext(ctx, query, file.Url, file.Name, file.Category, file.SizeText, file.SizeNumber,
		file.Description, file.Color, file.UserId, file.File.Name); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (r *repository) UpdateItem(ctx context.Context, file model.Item) error {
	query, _, err := goqu.Update(itemsTable).Set(file).Where(goqu.Ex{"id": file.Id}).ToSQL()
	if err != nil {
		return fmt.Errorf("cannot configure query: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("cannot update item: %w", err)
	}
	return nil
}

func (r *repository) GetItem(ctx context.Context, itemId int64) (model.Item, error) {
	var item model.Item
	query, _, err := goqu.From(itemsTable).Select().Where(goqu.Ex{"id": itemId}).ToSQL()
	if err != nil {
		return item, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.GetContext(ctx, &item, query)
	if err != nil {
		return item, fmt.Errorf("cannot get item: %w", err)
	}
	return item, nil
}

func (r *repository) GetItemsByUserId(ctx context.Context, userId int64) ([]model.Item, error) {
	var items []model.Item
	query, _, err := goqu.From(itemsTable).Select().Where(goqu.Ex{"user_id": userId}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("cannot get item: %w", err)
	}
	return items, nil
}

func (r *repository) GetItemsByIds(ctx context.Context, ids []int64) ([]model.Item, error) {
	var items []model.Item
	query, _, err := goqu.From(itemsTable).Select().Where(goqu.Ex{"id": ids}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("cannot get item: %w", err)
	}
	return items, nil
}

func (r *repository) CreateCapsule(ctx context.Context, capsule model.Capsule) error {
	query := `INSERT INTO capsules (name, user_id) values ($1, $2)`
	if _, err := r.db.ExecContext(ctx, query, capsule.Name, capsule.UserId); err != nil {
		return fmt.Errorf("cannot create capsule: %w", err)
	}
	return nil
}

func (r *repository) GetCapsule(ctx context.Context, capsuleId int64) (model.Capsule, error) {
	var capsule model.Capsule
	query, _, err := goqu.From(capsulesTable).Select().Where(goqu.Ex{"id": capsuleId}).ToSQL()
	fmt.Println(query)
	if err != nil {
		return model.Capsule{}, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.GetContext(ctx, &capsule, query)
	if err != nil {
		return model.Capsule{}, fmt.Errorf("cannot get item: %w", err)
	}
	return capsule, nil
}

func (r *repository) GetCapsulesByUser(ctx context.Context, userId int64) ([]model.Capsule, error) {
	var capsules []model.Capsule
	query, _, err := goqu.From(capsulesTable).Select().Where(goqu.Ex{"user_id": userId}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.SelectContext(ctx, &capsules, query)
	if err != nil {
		return nil, fmt.Errorf("cannot get item: %w", err)
	}
	return capsules, nil
}

func (r *repository) AddItemToCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error {
	query := `INSERT INTO capsules_items (item_id, capsule_id) values ($1, $2)`
	if _, err := r.db.ExecContext(ctx, query, itemCapsule.ItemId, itemCapsule.CapsuleId); err != nil {
		return fmt.Errorf("cannot add item to capsule: %w", err)
	}
	return nil
}

func (r *repository) DeleteItemFromCapsule(ctx context.Context, itemCapsule model.ItemCapsule) error {
	query := `DELETE FROM capsules_items WHERE item_id = $1 and capsule_id = $2`
	if _, err := r.db.ExecContext(ctx, query, itemCapsule.ItemId, itemCapsule.CapsuleId); err != nil {
		return fmt.Errorf("cannot delete item from capsule: %w", err)
	}
	return nil
}

func (r *repository) GetCapsuleItems(ctx context.Context, capsuleId int64) ([]int64, error) {
	var ids []int64
	query, _, err := goqu.From(capsulesItemsTable).Select("item_id").Where(goqu.Ex{"capsule_id": capsuleId}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.SelectContext(ctx, &ids, query)
	if err != nil {
		return nil, fmt.Errorf("cannot get items: %w", err)
	}
	return ids, nil
}

func (r *repository) DeleteCapsule(ctx context.Context, capsuleId int64) error {
	query := `DELETE FROM capsules WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, query, capsuleId); err != nil {
		return fmt.Errorf("cannot delete capsule:%w", err)
	}
	return nil
}

func (r *repository) UpdateCapsuleName(ctx context.Context, capsule model.Capsule) error {
	query := `UPDATE capsules SET name = $1 WHERE capsule_id = $2`
	if _, err := r.db.ExecContext(ctx, query, capsule.Name, capsule.Id); err != nil {
		return fmt.Errorf("cannot update capsule name:%w", err)
	}
	return nil
}

func (r *repository) DeleteItem(ctx context.Context, itemId int64) error {
	query := `DELETE FROM items WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, query, itemId); err != nil {
		return fmt.Errorf("cannot delete item:%w", err)
	}
	return nil
}

func (r *repository) CreateLook(ctx context.Context, look model.Look) error {
	query := `INSERT INTO looks (user_id, look_name, stylist_id) values ($1, $2, $3)`
	if _, err := r.db.ExecContext(ctx, query, look.UserId, look.LookName, look.StylistId); err != nil {
		return fmt.Errorf("cannot create look:%w", err)
	}
	return nil
}

func (r *repository) GetLook(ctx context.Context, lookId int64) (model.Look, error) {
	var look model.Look
	query, _, err := goqu.From("looks").Select().Where(goqu.Ex{"id": lookId}).ToSQL()
	if err != nil {
		return model.Look{}, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.GetContext(ctx, &look, query)
	if err != nil {
		return model.Look{}, fmt.Errorf("cannot get items: %w", err)
	}
	return look, nil
}

func (r *repository) AddItemToLook(ctx context.Context, lookId, itemId int64) error {
	query := `INSERT INTO looks_items (look_id, item_id) values ($1, $2)`
	if _, err := r.db.ExecContext(ctx, query, lookId, itemId); err != nil {
		return fmt.Errorf("cannot add to look:%w", err)
	}
	return nil
}

func (r *repository) DeleteItemFromLook(ctx context.Context, lookId, itemId int64) error {
	query := `DELETE FROM looks_items WHERE look_id = $1 AND item_id = $2`
	if _, err := r.db.ExecContext(ctx, query, lookId, itemId); err != nil {
		return fmt.Errorf("cannot add to look:%w", err)
	}
	return nil
}

func (r *repository) DeleteLook(ctx context.Context, lookId int64) error {
	query := `DELETE FROM looks WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, query, lookId); err != nil {
		return fmt.Errorf("cannot delete look:%w", err)
	}
	return nil
}

func (r *repository) DeleteLookItems(ctx context.Context, lookId int64) error {
	query := `DELETE FROM looks_items WHERE look_id = $1`
	if _, err := r.db.ExecContext(ctx, query, lookId); err != nil {
		return fmt.Errorf("cannot delete look items:%w", err)
	}
	return nil
}

func (r *repository) GetLookItems(ctx context.Context, lookId int64) ([]model.LooksItems, error) {
	var pairs []model.LooksItems
	query, _, err := goqu.From("looks_items").Select().Where(goqu.Ex{"look_id": lookId}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.SelectContext(ctx, &pairs, query)
	if err != nil {
		return nil, fmt.Errorf("cannot get items: %w", err)
	}
	return pairs, nil
}

func (r *repository) GetItemsArray(ctx context.Context, itemIds []int64) ([]model.Item, error) {
	var items []model.Item
	query, _, err := goqu.From(itemsTable).Select().Where(goqu.Ex{"id": itemIds}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("cannot configure query: %w", err)
	}
	err = r.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("cannot get items: %w", err)
	}
	return items, nil
}
