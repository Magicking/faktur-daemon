package db

import (
	"context"
	"log"

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

/*
type dbTransaction struct {
	gorm.Model
	MerkleRoot      string
	TransactionHash string
	Status          int
}*/

func SaveTx(ctx context.Context, root, txHash ethcommon.Hash) error {
	db := common.DBFromContext(ctx)
	dbtx := dbTransaction{
		MerkleRoot:      root.Hex(),
		TransactionHash: txHash.Hex(),
		Status:          NOT_SENT,
	}
	if err := db.Create(&dbtx).Error; err != nil {
		return err
	}
	return nil
}

func GetTxsFilter(ctx context.Context, state int) (*dbTransaction, error) {
	db := common.DBFromContext(ctx)
	var dbtx dbTransaction
	if err := db.Where(&dbTransaction{Status: state}).Last(&dbtx).Error; err != nil {
		return nil, err
	}
	return &dbtxnil, nil
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
