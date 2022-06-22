package main

import (
	"fmt"
	"github.com/foundVanting/opensea-stream-go/entity"
	"github.com/foundVanting/opensea-stream-go/opensea"
	"github.com/foundVanting/opensea-stream-go/types"
	"github.com/mitchellh/mapstructure"
	"github.com/nshafer/phx"
)

func main() {
	client := opensea.NewStreamClient(types.MAINNET, "apikey", phx.LogInfo, func(err error) {
		fmt.Println("opensea.NewStreamClient err:", err)
	})
	client.Connect()

	client.OnItemListed("collection-slug", func(response any) {
		var itemListedEvent entity.ItemListedEvent
		err := mapstructure.Decode(response, &itemListedEvent)
		if err != nil {
			fmt.Println("mapstructure.Decode err:", err)
		}
		fmt.Printf("%+v\n", itemListedEvent)
	})

	client.OnItemReceivedBid("collection-slug", func(response any) {
		var itemReceivedBidEvent entity.ItemReceivedBidEvent
		err := mapstructure.Decode(response, &itemReceivedBidEvent)
		if err != nil {
			fmt.Println("mapstructure.Decode err:", err)
		}
		fmt.Printf("%+v\n", itemReceivedBidEvent)
	})
	client.OnItemReceivedBid("collection-slug", func(response any) {
		var itemReceivedBidEvent entity.ItemReceivedBidEvent
		err := mapstructure.Decode(response, &itemReceivedBidEvent)
		if err != nil {
			fmt.Println("mapstructure.Decode err:", err)
		}
		fmt.Printf("%+v\n", itemReceivedBidEvent)
	})
	client.OnItemReceivedOffer("collection-slug", func(response any) {
		var itemReceivedOfferEvent entity.ItemReceivedOfferEvent
		err := mapstructure.Decode(response, &itemReceivedOfferEvent)
		if err != nil {
			fmt.Println("mapstructure.Decode err:", err)
		}
		fmt.Printf("%+v\n", itemReceivedOfferEvent)
	})

	select {}
}
