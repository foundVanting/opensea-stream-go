package opensea

import (
	"fmt"
	"github.com/foundVanting/opensea-stream-go/types"
	"github.com/nshafer/phx"
	"log"
	"net/url"
)

type StreamClient struct {
	socket   *phx.Socket
	channels map[string]*phx.Channel
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
	socket.Logger = phx.NewSimpleLogger(logLevel)
	return &StreamClient{
		socket:   socket,
		channels: make(map[string]*phx.Channel),
	}
}

func (s StreamClient) Connect() error {
	fmt.Println("Connecting to socket")
	return s.socket.Connect()
}
func (s *StreamClient) Disconnect() error {
	//s.socket.OnError()
	fmt.Println("Succesfully disconnected from socket")
	s.channels = make(map[string]*phx.Channel)
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
	s.channels[topic] = channel
	return
}
func (s StreamClient) getChannel(topic string) (channel *phx.Channel) {
	var ok bool
	channel, ok = s.channels[topic]
	if !ok {
		channel = s.createChannel(topic)
	}
	return channel
}

func (s StreamClient) on(eventType types.EventType, collectionSlug string, callback func(payload any)) func() {
	s.socket.Connect()
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
			delete(s.channels, collectionSlug)
			fmt.Printf("Succesfully left channel %s listening for %s\n", topic, eventType)
		})
	}

}

func collectionTopic(slug string) string {
	return fmt.Sprintf("collection:%s", slug)
}
func (s StreamClient) OnItemListed(collectionSlug string, Callback func(itemListedEvent any)) {
	s.on(types.ItemListed, collectionSlug, Callback)
}

func (s StreamClient) OnItemSold(collectionSlug string, Callback func(itemListedEvent any)) {
	s.on(types.ItemSold, collectionSlug, Callback)
}
func (s StreamClient) OnItemTransferred(collectionSlug string, Callback func(itemListedEvent any)) {
	s.on(types.ItemTransferred, collectionSlug, Callback)
}
func (s StreamClient) OnItemCancelled(collectionSlug string, Callback func(itemListedEvent any)) {
	s.on(types.ItemCancelled, collectionSlug, Callback)
}
