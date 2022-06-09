package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/nshafer/phx"
	"opensea-stream-go/entity"
	"opensea-stream-go/opensea"
	"opensea-stream-go/types"
)

func main() {
	client := opensea.NewStreamClient(types.MAINNET, "apikey", phx.LogInfo, func(err error) {
		fmt.Println("opensea.NewStreamClient err:", err)
	})
	client.Connect()

	client.OnItemListed("ens", func(response any) {
		var itemListedEvent entity.ItemListedEvent
		err := mapstructure.Decode(response, &itemListedEvent)
		if err != nil {
			fmt.Println("mapstructure.Decode err:", err)
		}
		fmt.Printf("%+v\n", itemListedEvent)
	})
	select {}

}
