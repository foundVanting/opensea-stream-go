package types

type EventType string

const (
	ItemMetadataUpdated EventType = "item_metadata_updated"
	ItemListed          EventType = "item_listed"
	ItemSold            EventType = "item_sold"
	ItemTransferred     EventType = "item_transferred"
	ItemReceivedOffer   EventType = "item_received_offer"
	ItemReceivedBid     EventType = "item_received_bid"
	ItemCancelled       EventType = "item_cancelle"
)

type Network string

const (
	MAINNET Network = "mainnet"
	TESTNET Network = "testnet"
)
