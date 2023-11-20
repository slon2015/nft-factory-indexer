package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nft-indexer-api/db"
	"nft-indexer-api/indexer"
	"nft-indexer-api/rpc"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func setupRouter(stor *db.EventStorage) *gin.Engine {
	r := gin.Default()

	r.GET("/collections", func(c *gin.Context) {
		c.JSON(http.StatusOK, stor.Collections)
	})

	r.GET("/mints", func(c *gin.Context) {
		c.JSON(http.StatusOK, stor.Mints)
	})

	return r
}

func setupConfig() {
    viper.AutomaticEnv()
	viper.SetConfigFile("./indexer.yaml")
	viper.SetConfigType("yaml")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindEnv("port")
	viper.BindEnv("address")
	viper.BindEnv("rpc.url")
	viper.BindEnv("creation-block")
	viper.BindEnv("indexer.parralel-requests")
	viper.BindEnv("indexer.blocks-per-request")
	viper.BindEnv("indexer.finalisation-blocks-count")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func setupIndexer(rpc *rpc.Rpc) (*indexer.Indexer, error) {
	indexer, err := indexer.NewIndexer(
		rpc, 
		viper.GetString("address"), 
		viper.GetUint64("creation-block"), 
		uint8(viper.GetUint("indexer.parralel-requests")),
		uint8(viper.GetUint("indexer.blocks-per-request")),
		uint8(viper.GetUint("indexer.finalisation-blocks-count")),
	)

	if err != nil {
		return nil, err
	}

	err = indexer.Start()

	if err != nil {
		return nil, err
	}

	return indexer, nil
}

func main() {
	setupConfig()

	rpc, err := rpc.NewEthRpc(context.Background(), viper.GetString("rpc.url"))

	if err != nil {
		log.Fatal(err)
	}

	indexer, err := setupIndexer(rpc)
	
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewEventStorage(indexer)
	r := setupRouter(store)

	defer indexer.Stop()
	defer store.Close()
	defer rpc.Close()

	// Listen and Server in 0.0.0.0:8080
	r.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
