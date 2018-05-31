package db

import (
	"context"
	"time"

	cmn "github.com/Magicking/faktur-daemon/common"
	"github.com/jinzhu/gorm"
)

type RootEntry struct {
	gorm.Model
	Root       string `gorm:"unique_index"`
	State      string
	TxHash     string
	LastUpdate time.Time
}

// Create a new entry
func CreateNewRoot(ctx context.Context, _root string) (*RootEntry, error) {
	db := cmn.DBFromContext(ctx)
	root := &RootEntry{Root: _root,
		State:      "new",
		TxHash:     "",
		LastUpdate: time.Now()}

	if err := db.Create(root).Error; err != nil {
		return nil, err
	}
	return root, nil
}

// Set internal state
func SetState(ctx context.Context, root *RootEntry, state string) error {
	db := cmn.DBFromContext(ctx)
	if err := db.Model(root).Updates(&RootEntry{State: state}).Error; err != nil {
		return err
	}
	return nil
}

// Finalize by setting the txhash
func FinalizeRoot(ctx context.Context, root *RootEntry, _txHash string) error {
	db := cmn.DBFromContext(ctx)
	if err := db.Model(root).Updates(&RootEntry{State: "final", TxHash: _txHash}).Error; err != nil {
		return err
	}
	return nil
}

func FilterByState(ctx context.Context, state string) ([]RootEntry, error) {
	db := cmn.DBFromContext(ctx)
	var roots []RootEntry
	cursor := db.Where(&RootEntry{State: state}).Find(&roots)
	if cursor.Error != nil {
		return nil, cursor.Error
	}
	if cursor.RecordNotFound() {
		return nil, nil
	}
	return roots, nil
}

func DBFixture(ctx context.Context) error {
	db := cmn.DBFromContext(ctx)
	if err := db.AutoMigrate(&RootEntry{}).Error; err != nil {
		db.Close()
		return err
	}
	return nil
}
