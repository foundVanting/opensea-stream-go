package opensea

import (
	"fmt"
	"github.com/foundVanting/opensea-stream-go/types"
	"github.com/nshafer/phx"
	"log"
	"net/url"
	"sync"
)

type StreamClient struct {
	socket   *phx.Socket
	channels *sync.Map //map[string]*phx.Channel
}

func NewStreamClient(network types.Network, token string, logLevel phx.LoggerLevel, onError func(error)) *StreamClient {
	m := map[types.Network]string{
		types.MAINNET: "wss://stream.openseabeta.com/socket",
		types.TESTNET: "wss://testnets-stream.openseabeta.com/socket",
	}
	socketUrl := fmt.Sprintf("%s?token=%s", m[network], token)

	endPoint, _ := url.Parse(socketUrl)
	socket := phx.NewSocket(endPoint)

	socket.OnError(onError)
	socket.OnClose(func() {
		err := socket.Reconnect()
		if err != nil {
			onError(err)
		}
	})
	socket.Logger = phx.NewSimpleLogger(logLevel)
	return &StreamClient{
		socket:   socket,
		channels: &sync.Map{},
	}
}

func (s StreamClient) Connect() error {
	fmt.Println("Connecting to socket")
	return s.socket.Connect()
}
func (s *StreamClient) Disconnect() error {
	//s.socket.OnError()
	fmt.Println("Succesfully disconnected from socket")
	s.channels = &sync.Map{}
	return s.socket.Disconnect()
}
func (s *StreamClient) createChannel(topic string) (channel *phx.Channel) {
	channel = s.socket.Channel(topic, nil)
	join, err := channel.Join()
	if err != nil {
		fmt.Println(err)
		return
	}
	join.Receive("ok", func(response any) {
		log.Println("Joined channel:", channel.Topic(), response)
	})
	join.Receive("error", func(response any) {
		log.Println("failed 2 joined channel:", channel.Topic(), response)
	})
	s.channels.Store(topic, channel)
	return
}
func (s StreamClient) getChannel(topic string) (channel *phx.Channel) {

	value, ok := s.channels.Load(topic)
	if ok {
		channel = value.(*phx.Channel)
		return
	}

	channel = s.createChannel(topic)

	return channel
}

func (s StreamClient) on(eventType types.EventType, collectionSlug string, callback func(payload any)) func() {
	topic := collectionTopic(collectionSlug)
	fmt.Printf("Fetching channel %s\n", topic)
	channel := s.getChannel(topic)
	fmt.Printf("Subscribing to %s events on %s\n", eventType, topic)
	channel.On(string(eventType), callback)
	return func() {
		fmt.Printf("Unsubscribing from %s events on %s\n", eventType, topic)
		leave, err := channel.Leave()
		if err != nil {
			fmt.Println("channel.Leave err:", err)
		}
		leave.Receive("ok", func(response any) {
			s.channels.Delete(topic)
			fmt.Printf("Succesfully left channel %s listening for %s\n", topic, eventType)
		})
	}

}

func collectionTopic(slug string) string {
	return fmt.Sprintf("collection:%s", slug)
}
func (s StreamClient) OnItemListed(collectionSlug string, Callback func(itemListedEvent any)) func() {
	return s.on(types.ItemListed, collectionSlug, Callback)
}

func (s StreamClient) OnItemSold(collectionSlug string, Callback func(itemSoldEvent any)) func() {
	return s.on(types.ItemSold, collectionSlug, Callback)
}
func (s StreamClient) OnItemTransferred(collectionSlug string, Callback func(itemTransferredEvent any)) func() {
	return s.on(types.ItemTransferred, collectionSlug, Callback)
}
func (s StreamClient) OnItemCancelled(collectionSlug string, Callback func(itemCancelledEvent any)) func() {
	return s.on(types.ItemCancelled, collectionSlug, Callback)
}
func (s StreamClient) OnItemReceivedBid(collectionSlug string, Callback func(itemReceivedBidEvent any)) func() {
	return s.on(types.ItemReceivedBid, collectionSlug, Callback)
}
func (s StreamClient) OnItemReceivedOffer(collectionSlug string, Callback func(itemReceivedOfferEvent any)) func() {
	return s.on(types.ItemReceivedOffer, collectionSlug, Callback)
}
func (s StreamClient) OnItemMetadataUpdated(collectionSlug string, Callback func(itemMetadataUpdatedEvent any)) func() {
	return s.on(types.ItemMetadataUpdated, collectionSlug, Callback)
}
