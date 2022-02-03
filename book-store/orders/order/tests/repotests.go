package tests

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	"github.com/Vesninovich/go-tasks/book-store/orders/order"
)

var orders = []order.CreateDTO{
	{Description: "a", BookID: uuid.New()},
	{Description: "b", BookID: uuid.New()},
}
var ctx = context.Background()

// Constructor is function that constructs repository to test
type Constructor func(*testing.T) order.Repository

// RepoGetAll tests getting all data
func RepoGetAll(t *testing.T, c Constructor) {
	repo := setup(t, c)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Error while getting all stored items: %s", err)
	}
	if len(stored) != len(orders) {
		t.Errorf("Expected to have %d items stored, got %d", len(orders), len(stored))
	}
}

// RepoGet tests getting item by id
func RepoGet(t *testing.T, c Constructor) {
	repo := setup(t, c)
	stored, _ := repo.GetAll(ctx)
	for _, item := range stored {
		found, err := repo.Get(ctx, item.ID)
		if err != nil {
			t.Fatalf("Error while getting item: %s", err)
		}
		if found.Description != item.Description || found.ID != item.ID {
			t.Error("Got wrong item")
		}
	}
}

// RepoGetNonExisting tests getting non-existing item by id
func RepoGetNonExisting(t *testing.T, c Constructor) {
	repo := setup(t, c)
	_, err := repo.Get(ctx, uuid.New())
	checkNotFound(t, err)
}

// RepoCreate tests creating items
func RepoCreate(t *testing.T, c Constructor) {
	repo := setup(t, c)
	repo.GetAll(ctx)
}

// RepoUpdate tests updating item
func RepoUpdate(t *testing.T, c Constructor) {
	repo, stored := setupMutation(t, c)
	id := stored[0].ID
	desc := "d"
	replaced, err := repo.Update(ctx, order.DTO{
		ID: id,
		CreateDTO: order.CreateDTO{
			Description: desc,
			BookID:      uuid.New(),
		},
	})
	if err != nil {
		t.Fatalf("Error while updating item: %s", err)
	}
	if replaced.Description != desc {
		t.Errorf("Expected to update item with data %s, got %s", desc, replaced.Description)
	}
}

// RepoUpdateNonExisting tests updating non-existing item
func RepoUpdateNonExisting(t *testing.T, c Constructor) {
	repo := setup(t, c)
	_, err := repo.Update(ctx, order.DTO{
		ID: uuid.New(),
		CreateDTO: order.CreateDTO{
			Description: "",
			BookID:      uuid.New(),
		},
	})
	checkNotFound(t, err)
}

// RepoUpdateDeleted tests updating deleted item
func RepoUpdateDeleted(t *testing.T, c Constructor) {
	repo, id, _ := setupAlreadyDeleted(t, c)
	_, err := repo.Update(ctx, order.DTO{
		ID: id,
		CreateDTO: order.CreateDTO{
			Description: "",
			BookID:      uuid.New(),
		},
	})
	checkNotFound(t, err)
}

// RepoUpdateWithSomeDeleted tests updating item if some are deleted
func RepoUpdateWithSomeDeleted(t *testing.T, c Constructor) {
	repo, id, stored := setupAlreadyDeleted(t, c)
	var item order.DTO
	for _, item = range stored {
		if item.ID != id {
			break
		}
	}
	desc := "c"
	replaced, err := repo.Update(ctx, order.DTO{
		ID: item.ID,
		CreateDTO: order.CreateDTO{
			Description: desc,
			BookID:      uuid.New(),
		},
	})
	if err != nil {
		t.Fatalf("Error while updating item with some deleted: %s", err)
	}
	if replaced.Description != desc {
		t.Errorf("Expected to update item with data %s, got %s", desc, replaced.Description)
	}
}

// RepoDelete tests deleting item
func RepoDelete(t *testing.T, c Constructor) {
	repo, stored := setupMutation(t, c)
	for _, item := range stored {
		_, err := repo.Delete(ctx, item.ID)
		if err != nil {
			t.Fatalf("Error while deleting item: %s", err)
		}
		_, err = repo.Get(ctx, item.ID)
		checkNotFound(t, err)
	}
}

// RepoDeleteTwice tests deleting item twice
func RepoDeleteTwice(t *testing.T, c Constructor) {
	repo, id, _ := setupAlreadyDeleted(t, c)
	_, err := repo.Delete(ctx, id)
	checkNotFound(t, err)
}

// RepoDeleteNonExisting tests deleting non-existing item
func RepoDeleteNonExisting(t *testing.T, c Constructor) {
	repo := setup(t, c)
	_, err := repo.Delete(ctx, uuid.New())
	checkNotFound(t, err)
}

func findByDescription(name string, data []order.DTO, t *testing.T) order.DTO {
	for _, item := range data {
		if item.Description == name {
			return item
		}
	}
	t.Errorf("Item with name %s not found", name)
	return order.DTO{}
}

func setup(t *testing.T, c Constructor) order.Repository {
	repo := c(t)
	for _, o := range orders {
		_, err := repo.Create(ctx, o)
		if err != nil {
			t.Fatalf("Error while creating order %s: %s", o, err)
		}
	}
	return repo
}

func setupMutation(t *testing.T, c Constructor) (order.Repository, []order.DTO) {
	repo := setup(t, c)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Error while fetching all data: %s", err)
	}
	return repo, stored
}

func setupAlreadyDeleted(t *testing.T, c Constructor) (order.Repository, uuid.UUID, []order.DTO) {
	repo, stored := setupMutation(t, c)
	id := stored[0].ID
	_, err := repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Error deleting item: %s", err)
	}
	return repo, id, stored
}

func checkNotFound(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Expected to get NotFound error")
	}
	if _, typeCorrect := err.(*commonerrors.NotFound); !typeCorrect {
		t.Errorf("Expected to get error of type *commonerrors.NotFound, got %T", err)
	}
}
