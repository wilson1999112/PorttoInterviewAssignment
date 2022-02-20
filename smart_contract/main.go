package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"wilson/testproject/api"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

const (
	PRIVATE_KEY   string = "89e65031242e5711bda8780517f69fce3f6eb7a166a72c50eaa2be083615e188"
	TEST_ENDPOINT string = "http://127.0.0.1:7545"
)

func main() {
	client, err := ethclient.Dial(TEST_ENDPOINT)
	if err != nil {
		panic(err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	address, _, _, err := api.DeployApi(
		get_auth_from_pk(
			client,
			chainID,
			PRIVATE_KEY,
		),
		client,
	)
	if err != nil {
		panic(err)
	}
	conn, err := api.NewApi(common.HexToAddress(address.Hex()), client)
	if err != nil {
		panic(err)
	}

	// here is router part
	router := gin.Default()
	router.GET("/owner", func(c *gin.Context) {
		res, err := conn.Owner(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"owner": res})
	})
	router.GET("/winer", func(c *gin.Context) {
		res, err := conn.Winer(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"winer": res})
	})
	router.POST("/set/:num", func(c *gin.Context) {
		num := c.Param("num")
		n, _ := strconv.Atoi(num)
		auth := request_handler(client, chainID, c.Request.Body)
		res, err := conn.SetNum(auth, big.NewInt(int64(n)))
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"res": res})
	})
	router.POST("/guess/:num", func(c *gin.Context) {
		num := c.Param("num")
		n, _ := strconv.Atoi(num)
		auth := request_handler(client, chainID, c.Request.Body)
		res, err := conn.GuessNum(auth, big.NewInt(int64(n)))
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"res": res})
	})
	router.Run(":8080")

}
func request_handler(client *ethclient.Client, chainID *big.Int, body io.Reader) *bind.TransactOpts {
	var v map[string]interface{}
	err := json.NewDecoder(body).Decode(&v)
	if err != nil {
		panic(err)
	}
	return get_auth_from_pk(client, chainID, v["private_key"].(string))
}
func get_auth_from_pk(client *ethclient.Client, chainID *big.Int, private_key string) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		panic(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), get_public_address(privateKey))
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(1000000)
	return auth
}
func get_public_address(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
