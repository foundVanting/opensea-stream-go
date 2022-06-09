# OpenSea Stream API - golang SDK

A Golang SDK for receiving updates from the OpenSea Stream API - pushed over websockets. We currently support the following event types on a per-collection basis:

- item listed
- item sold
- item transferred
- item metadata updates
- item cancelled
- item received offer
- item received bid


Documentation: https://docs.opensea.io/reference/stream-api-overview

# Installation
This module requires Go 1.18 or later.

`go get github.com/foundVanting/opensea-stream-go`

# Getting Started

## Authentication

In order to make onboarding easy, we've integrated the OpenSea Stream API with our existing API key system. The API keys you have been using for the REST API should work here as well. If you don't already have one, request an API key from us [here](https://docs.opensea.io/reference/request-an-api-key).

## Simple example



```golang
func main() {
    client := opensea.NewStreamClient(types.MAINNET, "api-key", phx.LogInfo, func(err error) {
        fmt.Println("NewStreamClient err:", err)
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
```

You can also optionally pass in:

- a `network` if you would like to access testnet networks.
    - The default value is `Network.MAINNET`, which represents the following blockchains: Ethereum, Polygon mainnet, Klaytn mainnet, and Solana mainnet
    - Can also select `Network.TESTNET`, which represents the following blockchains: Rinkeby, Polygon testnet (Mumbai), and Klaytn testnet (Baobab).
- `apiUrl` if you would like to access another OpenSea Stream API endpoint. Not needed if you provide a network or use the default values.
- an `onError` callback to handle errors. The default behavior is to `console.error` the error.
- a `logLevel` to set the log level. The default is `LogLevel.INFO`.

## Available Networks

The OpenSea Stream API is available on the following networks:

### Mainnet

`wss://stream.openseabeta.com/socket`

Mainnet supports events from the following blockchains: Ethereum, Polygon mainnet, Klaytn mainnet, and Solana mainnet.

### Testnet

`wss://testnets-stream.openseabeta.com/socket`

Testnet supports events from the following blockchains: Rinkeby, Polygon testnet (Mumbai), and Klaytn testnet (Baobab).

To create testnet instance of the client, you can create it with the following arguments:
