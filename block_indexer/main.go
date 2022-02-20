package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/viney-shih/goroutines"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	UserName      string = "root"
	Password      string = "3860"
	Addr          string = "127.0.0.1"
	Port          int    = 3306
	Database      string = "letmetest"
	MaxLifetime   int    = 10
	MaxOpenConns  int    = 10
	MaxIdleConns  int    = 10
	POOL_SIZE     int    = 10
	N             int64  = 14236000
	TEST_ENDPOINT string = "https://mainnet.infura.io/v3/c54650de28294071b4980a9217f2a724"
)

type Block struct {
	Hash        string             `gorm:"varchar(64);primaryKey" json:"hash"`
	Difficulty  string             `json:"difficulty"`
	Extra       []byte             `json:"extra"`
	GasLimit    uint64             `json:"gasLimit"`
	GasUsed     uint64             `json:"gasUsed"`
	Nonce       uint64             `json:"nonce"`
	Number      string             `json:"number"`
	ParentHash  common.Hash        `json:"parentHash"`
	ReceiptHash common.Hash        `json:"receiptHash"`
	UncleHash   common.Hash        `json:"uncleHash"`
	Size        common.StorageSize `json:"size"`
	Time        uint64             `json:"time"`
	CreatedAt   time.Time          `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}
type Transaction struct {
	Hash        string          `gorm:"varchar(64);primaryKey" json:"hash"`
	BlockHash   string          `json:"block_hash"`
	TxFrom      common.Address  `json:"tx_from"`
	BlockNumber string          `json:"block_number"`
	Gas         uint64          `json:"gas"`
	GasPrice    string          `json:"gas_price"`
	Nonce       uint64          `json:"nonce"`
	TxTo        *common.Address `json:"tx_to"`
	Value       string          `json:"value"`
	Type        uint8           `json:"type"`
	CreatedAt   time.Time       `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time       `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

func index_block(conn *gorm.DB, chainID *big.Int, idx *big.Int, client *ethclient.Client) {
	block, _ := client.BlockByNumber(context.TODO(), idx)
	b := Block{
		Hash:        block.Hash().String(),
		Difficulty:  block.Difficulty().String(),
		Extra:       block.Extra(),
		GasLimit:    block.GasLimit(),
		GasUsed:     block.GasUsed(),
		Nonce:       block.Nonce(),
		Number:      block.Number().String(),
		ParentHash:  block.ParentHash(),
		ReceiptHash: block.ReceiptHash(),
		UncleHash:   block.UncleHash(),
		Size:        block.Size(),
		Time:        block.Time(),
	}
	conn.Create(&b)
	txs := []Transaction{}
	for _, tx := range block.Transactions() {
		from, _ := types.Sender(types.LatestSignerForChainID(chainID), tx)
		txs = append(txs,
			Transaction{
				Hash:        tx.Hash().String(),
				BlockHash:   block.Hash().String(),
				TxFrom:      from,
				BlockNumber: block.Number().String(),
				Gas:         tx.Gas(),
				GasPrice:    tx.GasPrice().String(),
				Nonce:       tx.Nonce(),
				TxTo:        tx.To(),
				Value:       tx.Value().String(),
				Type:        tx.Type(),
			})
	}
	conn.Create(&txs)
}
func main() {

	client, err := ethclient.Dial(TEST_ENDPOINT)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("connection to mysql failed:", err)
	}
	db, err := conn.DB()
	if err != nil {
		log.Fatal("get db failed:", err)
	}
	db.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)

	conn.Debug().AutoMigrate(&Block{})
	conn.Debug().AutoMigrate(&Transaction{})

	last_number := big.NewInt(N - 1)
	p := goroutines.NewPool(POOL_SIZE)
	defer p.Release()
	for {
		header, _ := client.HeaderByNumber(context.TODO(), nil)
		if last_number.String() != header.Number.String() {
			for i := big.NewInt(0).Add(last_number, big.NewInt(1)); i.Cmp(header.Number) <= 0; i.Add(i, big.NewInt(1)) {
				p.Schedule(func() {
					index_block(conn, chainID, i, client)
				})
			}
			last_number = header.Number
		}
		log.Println("fully synced")
		time.Sleep(5 * time.Second)
	}

}
