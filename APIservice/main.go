package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

const (
	TEST_ENDPOINT string = "https://mainnet.infura.io/v3/c54650de28294071b4980a9217f2a724"
)

type M map[string]interface{}

func main() {
	router := gin.Default()
	client, err := ethclient.Dial(TEST_ENDPOINT)
	if err != nil {
		log.Fatal("Whoops! client", err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("Whoops! chainID", err)
	}
	router.GET("/blocks", func(c *gin.Context) {
		limit := c.Query("limit")
		n, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			panic(err)
		}
		myMapSlice := get_n_block(client, n)
		c.IndentedJSON(200, gin.H{"blocks": myMapSlice})

	})
	router.GET("/blocks/:id", func(c *gin.Context) {
		bid := c.Param("id")
		n, err := strconv.ParseInt(bid, 10, 64)
		if err != nil {
			panic(err)
		}
		result := get_block(client, n)
		c.IndentedJSON(200, result)

	})
	router.GET("/transaction/:txHash", func(c *gin.Context) {
		txHash_str := c.Param("txHash")
		txHash := common.HexToHash(txHash_str)
		result := get_tx(client, txHash, chainID)
		c.IndentedJSON(200, result)

	})
	router.Run(":8080")
}
func get_tx(client *ethclient.Client, txHash common.Hash, chainID *big.Int) M {
	tx, _, err := client.TransactionByHash(context.TODO(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	from, err := types.Sender(types.LatestSignerForChainID(chainID), tx)
	if err != nil {
		log.Fatal(err)
	}
	receipt, err := client.TransactionReceipt(context.TODO(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	m := M{
		"tx_hash": tx.Hash(),
		"from":    from,
		"to":      tx.To(),
		"nonce":   tx.Nonce(),
		"data":    "0x" + hex.EncodeToString(tx.Data()),
		"value":   tx.Value(),
	}
	var logs []M
	for _, my_log := range receipt.Logs {
		logs = append(logs, M{
			"index": my_log.Index,
			"data":  "0x" + hex.EncodeToString(my_log.Data),
		})
	}
	m["logs"] = logs
	return m

}
func get_block(client *ethclient.Client, n int64) M {
	block, err := client.BlockByNumber(context.TODO(), big.NewInt(n))
	if err != nil {
		log.Fatal(err)
	}
	m := M{
		"block_num":   block.Number().Uint64(),
		"block_hash":  block.Hash().Hex(),
		"block_time":  block.Time(),
		"parent_hash": block.ParentHash(),
	}
	var t_hashes []common.Hash
	for _, tx := range block.Transactions() {
		t_hashes = append(t_hashes, tx.Hash())
	}
	m["transactions"] = t_hashes
	return m
}
func get_n_block(client *ethclient.Client, limit int64) []M {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	start := big.NewInt(0).Sub(header.Number, big.NewInt(limit))
	end := header.Number
	ctx := context.TODO()
	var myMapSlice []M
	for i := new(big.Int).Set(start); i.Cmp(end) <= 0; i.Add(i, big.NewInt(1)) {
		block, err := client.BlockByNumber(ctx, i)
		if err != nil {
			log.Fatal(err)
		}
		m := M{
			"block_num":   block.Number().Uint64(),
			"block_hash":  block.Hash().Hex(),
			"block_time":  block.Time(),
			"parent_hash": block.ParentHash(),
		}
		myMapSlice = append(myMapSlice, m)
	}
	return myMapSlice
}
