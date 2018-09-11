package db

import (
	"context"

	log "github.com/sirupsen/logrus"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Magicking/faktur-daemon/common"
	"github.com/Magicking/faktur-daemon/merkle"
	"github.com/jinzhu/gorm"
)

type DbReceipt struct {
	gorm.Model
	Targethash    string
	Proofs        string // unique
	MerkleRoot    string
	TransactionId int
	Transaction   *Dbtransaction
}

type Dbtransaction struct {
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

func GetTxByRoot(ctx context.Context, root ethcommon.Hash) (*Dbtransaction, error) {
	db := common.DBFromContext(ctx)

	var tx Dbtransaction
	cursor := db.Where(Dbtransaction{MerkleRoot: root.Hex()}).First(&tx)
	if cursor.RecordNotFound() {
		return nil, nil
	}
	if cursor.Error != nil {
		return nil, cursor.Error
	}
	return &tx, nil
}

func GetReceiptsByHash(ctx context.Context, targetHash ethcommon.Hash) ([]DbReceipt, error) {
	db := common.DBFromContext(ctx)

	var rcpts []DbReceipt
	cursor := db.Preload("Transaction").Where(DbReceipt{Targethash: targetHash.Hex()}).Find(&rcpts)
	if cursor.RecordNotFound() {
		return nil, nil
	}
	if cursor.Error != nil {
		return nil, cursor.Error
	}
	return rcpts, nil
}

func SaveReceipt(ctx context.Context, proofs *merkle.Branch, targetHash merkle.Hashable, root ethcommon.Hash, tx *Dbtransaction) error {
	db := common.DBFromContext(ctx)
	dbrcpt := DbReceipt{
		Targethash:  targetHash.Hex(),
		Proofs:      proofs.String(),
		MerkleRoot:  root.Hex(),
		Transaction: tx,
	}
	if err := db.Create(&dbrcpt).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTx(ctx context.Context, root ethcommon.Hash, txHash *ethcommon.Hash, state int) error {
	db := common.DBFromContext(ctx)
	//TODO Facktorize & optimise
	var tx Dbtransaction
	cursor := db.Where("merkle_root = ?", root.Hex()).Find(&tx)
	if cursor.RecordNotFound() {
		var _txHash string
		if txHash != nil {
			_txHash = txHash.Hex()
		}
		dbtx := Dbtransaction{
			MerkleRoot:      root.Hex(),
			TransactionHash: _txHash,
			Status:          state,
		}
		if err := db.Create(&dbtx).Error; err != nil {
			return err
		}
		return nil
	}
	cursor = db.Model(&Dbtransaction{}).Where("merkle_root = ?", root.Hex())
	if cursor.Error != nil {
		return cursor.Error
	}
	if txHash != nil {
		cursor = cursor.Updates(&Dbtransaction{Status: state, TransactionHash: txHash.Hex()})
	} else {
		cursor = cursor.Updates(&Dbtransaction{Status: state})
	}
	if cursor.Error != nil {
		return cursor.Error
	}
	return nil
}

func FilterByState(ctx context.Context, state int) (dbtx []*Dbtransaction, err error) {
	db := common.DBFromContext(ctx)
	dbtx = make([]*Dbtransaction, 0)
	cursor := db.Where(&Dbtransaction{Status: state}).Find(&dbtx)
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
	if err := db.AutoMigrate(&DbReceipt{}).Error; err != nil {
		db.Close()
		log.Fatalf("Could not migrate models to database: %v", err)
	}
	if err := db.AutoMigrate(&Dbtransaction{}).Error; err != nil {
		db.Close()
		log.Fatalf("Could not migrate models to database: %v", err)
	}
}
