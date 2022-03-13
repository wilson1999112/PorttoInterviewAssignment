package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	UserName     string = "root"
	Password     string = "3860"
	Addr         string = "127.0.0.1"
	Port         int    = 3306
	Database     string = "new_schema"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
	PoolSize     int    = 16
	NumberStart  int    = 14378628
	DbBufferSize int    = 64
	StableSetNum int    = 20
	TestEndPoint string = "https://mainnet.infura.io/v3/cde840a2236a4727b3eb9d3c9032cde2"
)

type Block struct {
	ID         uint64 `gorm:"primary_key;auto_increment"`
	Hash       common.Hash
	ParentHash common.Hash
	Number     int
	Time       uint64
	Stable     bool
}
type Transaction struct {
	ID        uint64 `gorm:"primary_key;auto_increment"`
	BlockHash common.Hash
	Hash      common.Hash
	TxFrom    common.Address
	TxTo      common.Address
	Nounce    uint64
	Stable    bool
}
type Blocks struct {
	blocks [DbBufferSize]Block
	index  int
}
type Txs struct {
	transactions [DbBufferSize]Transaction
	index        int
}

var wg = sync.WaitGroup{}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func (bs *Blocks) writeDb(conn *gorm.DB, chainID *big.Int, b *types.Block, stable bool) {
	numberString := b.Number().String()
	number, _ := strconv.Atoi(numberString)
	bs.blocks[bs.index] = Block{
		Hash:       b.Hash(),
		Stable:     stable,
		ParentHash: b.ParentHash(),
		Number:     number,
		Time:       b.Time(),
	}
	bs.index++
	if bs.index == DbBufferSize {
		conn.Create(&bs.blocks)
		bs.index = 0
	}
	txsBuffer := &Txs{}
	for _, tx := range b.Transactions() {
		txsBuffer.writeDb(conn, chainID, tx, b.Hash(), stable)
	}
	txsBuffer.flush(conn)
}
func (bs *Blocks) flush(conn *gorm.DB) {
	if bs.index > 0 {
		conn.Create(bs.blocks[:bs.index])
		bs.index = 0
	}
}
func (txs *Txs) writeDb(conn *gorm.DB, chainID *big.Int, tx *types.Transaction, blockHash common.Hash, stable bool) {
	from, _ := types.Sender(types.LatestSignerForChainID(chainID), tx)
	var to common.Address
	if tx.To() != nil {
		to = *tx.To()
	}
	txs.transactions[txs.index] = Transaction{
		BlockHash: blockHash,
		Hash:      tx.Hash(),
		Stable:    stable,
		TxFrom:    from,
		TxTo:      to,
		Nounce:    tx.Nonce(),
	}
	txs.index++
	if txs.index == DbBufferSize {
		conn.Create(&txs.transactions)
		txs.index = 0
	}
}
func (txs *Txs) flush(conn *gorm.DB) {
	if txs.index > 0 {
		conn.Create(txs.transactions[:txs.index])
		txs.index = 0
	}
}
func getLatestBlockNumber(client *ethclient.Client) int {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		time.Sleep(5 * time.Second)
		return getLatestBlockNumber(client)
	}
	numberString := header.Number.String()
	latestNumber, _ := strconv.Atoi(numberString)
	return latestNumber
}
func getBlockByNumber(client *ethclient.Client, index int) *types.Block {
	block, err := client.BlockByNumber(context.TODO(), big.NewInt(int64(index)))
	if err != nil {
		time.Sleep(5 * time.Second)
		return getBlockByNumber(client, index)
	}
	return block
}
func blockReader(conn *gorm.DB, chainID *big.Int, client *ethclient.Client, ch <-chan int) {
	blockBuffer := &Blocks{}
	for i := range ch {
		block := getBlockByNumber(client, i)
		blockBuffer.writeDb(conn, chainID, block, true)
	}
	blockBuffer.flush(conn)
	wg.Done()
}
func blockTracker(conn *gorm.DB, chainID *big.Int, client *ethclient.Client, i int, ch <-chan struct{}) {
	blockBuffer := &Blocks{}
	block := getBlockByNumber(client, i)
	oldHash := block.Hash()
	blockBuffer.writeDb(conn, chainID, block, false)
	blockBuffer.flush(conn)
Loop:
	for {
		select {
		case <-ch:
			if oldHash != block.Hash() {
				// delete the old data
				conn.Where("hash = ?", oldHash).Delete(&Block{})
				conn.Where("block_hash = ?", oldHash).Delete(&Transaction{})
				blockBuffer.writeDb(conn, chainID, block, true)
				blockBuffer.flush(conn)
				fmt.Printf("detect different blocks %v: %v, %v\n", i, oldHash, block.Hash())
			} else {
				// update stable=true
				conn.Model(&Block{}).Where("hash = ?", oldHash).Update("stable", true)
				conn.Model(&Transaction{}).Where("block_hash = ?", oldHash).Update("stable", true)
			}
			break Loop
		default:
			block = getBlockByNumber(client, i)
			time.Sleep(5 * time.Second)
		}
	}
	fmt.Printf("block: %v is stable now\n", i)
}
func main() {

	client, err := ethclient.Dial(TestEndPoint)
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

	// here I only read the stable blocks
	stableNumber := NumberStart - 1
	latestNumber := getLatestBlockNumber(client)
	fmt.Printf("stableNumber: %v, latestNumber: %v\n", stableNumber, latestNumber)
	if stableNumber < latestNumber-StableSetNum {
		realPoolSize := min(PoolSize, latestNumber-stableNumber-StableSetNum)
		wg.Add(realPoolSize)
		ch := make(chan int, realPoolSize)
		for i := 0; i < realPoolSize; i++ {
			go blockReader(conn, chainID, client, ch)
		}
		for i := stableNumber + 1; i <= latestNumber-StableSetNum; i++ {
			ch <- i
		}
		close(ch)
		wg.Wait()
		stableNumber = latestNumber - StableSetNum
	}

	// here read the latest unstable blocks and track them util it becomes a stable block
	queue := []int{}
	doneMap := make(map[int]chan struct{})
	for {
		latestNumber := getLatestBlockNumber(client)
		fmt.Printf("stableNumber: %v, latestNumber: %v\n", stableNumber, latestNumber)
		for i := len(queue); i < latestNumber-stableNumber; i++ {
			ch := make(chan struct{})
			doneMap[stableNumber+i] = ch
			queue = append(queue, stableNumber+i)
			go blockTracker(conn, chainID, client, stableNumber+i, ch)
		}
		for len(queue) > StableSetNum {
			doneMap[queue[0]] <- struct{}{}
			queue = queue[1:]
			stableNumber++
		}
		time.Sleep(5 * time.Second)
	}

}
