package db

import (
	"context"

	log "github.com/sirupsen/logrus"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Magicking/faktur-daemon/common"
	"github.com/Magicking/faktur-daemon/merkle"
	"github.com/jinzhu/gorm"
)

type dbReceipt struct {
	gorm.Model
	Targethash  string
	Proofs      string // unique
	MerkleRoot  string
	Privatekey  string // unique
	Transaction dbTransaction
}

type dbTransaction struct {
	gorm.Model
	MerkleRoot      string
	TransactionHash string
	Status          int
}

const (
	NOT_SENT = iota
	RETRY
	WAITING_CONFIRMATION
	SENT
)

func SaveReceipt(ctx context.Context, proofs *merkle.Branch, targetHash merkle.Hashable, root ethcommon.Hash) error {
	db := common.DBFromContext(ctx)
	dbrcpt := dbReceipt{
		Targethash: targetHash.Hex(),
		Proofs:     proofs.String(),
		MerkleRoot: root.Hex(),
	}
	if err := db.Create(&dbrcpt).Error; err != nil {
		return err
	}
	return nil
}

func SaveTx(ctx context.Context, root ethcommon.Hash, txHash *ethcommon.Hash, state int) error {
	db := common.DBFromContext(ctx)
	var _txHash string
	if txHash != nil {
		_txHash = txHash.Hex()
	}
	dbtx := dbTransaction{
		MerkleRoot:      root.Hex(),
		TransactionHash: _txHash,
		Status:          state,
	}
	if err := db.Create(&dbtx).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTx(ctx context.Context, root ethcommon.Hash, txHash *ethcommon.Hash, state int) error {
	db := common.DBFromContext(ctx)
	cursor := db.Model(&dbTransaction{}).Where("merkle_root = ?", root.Hex())
	if cursor.Error != nil {
		return cursor.Error
	}
	if txHash != nil {
		cursor = cursor.Updates(&dbTransaction{Status: state, TransactionHash: txHash.Hex()})
	} else {
		cursor = cursor.Updates(&dbTransaction{Status: state})
	}
	if cursor.Error != nil {
		return cursor.Error
	}
	return nil
}

func FilterByState(ctx context.Context, state int) (dbtx []*dbTransaction, err error) {
	db := common.DBFromContext(ctx)
	dbtx = make([]*dbTransaction, 0)
	cursor := db.Where(&dbTransaction{Status: state}).Find(&dbtx)
	if cursor.Error != nil {
		return nil, err
	}
	if cursor.RecordNotFound() {
		return nil, nil
	}
	return dbtx, nil
}

// Get    Txs w/ STATUS
// Update Tx  w/ NEW_STATUS

func MigrateDatabase(ctx context.Context) {
	db := common.DBFromContext(ctx)
	if err := db.AutoMigrate(&dbReceipt{}).Error; err != nil {
		db.Close()
		log.Fatalf("Could not migrate models to database: %v", err)
	}
	if err := db.AutoMigrate(&dbTransaction{}).Error; err != nil {
		db.Close()
		log.Fatalf("Could not migrate models to database: %v", err)
	}
}
