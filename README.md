# Donut Challenge

This repo contains code for a recruitment challenge from [Donut](trydonut.com).
It exposes two actions utilising the [Coinbase
API](https://docs.pro.coinbase.com/): one for getting the current rate between
two currencies, and one for placing an order.

## Requirements

I have developed this using Go `v1.11.4 darwin/amd64` and have not tested any
other versions or distributions.

## Installation

```
go get github.com/willclarktech/donut-challenge
```

## Usage

### `GetCurrentRate`

Just call `donutchallenge.GetCurrentRate` with your desired currencies (all
supported currencies are exported as named constants for convenience).

```go
import (
	donutchallenge "github.com/willclarktech/donut-challenge"
	"fmt"
)

func main() {
	rate, _ := donutchallenge.GetCurrentRate(donutchallenge.ETH, donutchallenge.USD)
	fmt.Println("The rate is:", rate)
}
```

### `PlaceOrder`

Placing an order requires authentication with the Coinbase API as described
[here](https://docs.pro.coinbase.com/#creating-a-request). This library assumes
the following environmental variables have been set:

- `CB_ACCESS_KEY`
- `CB_ACCESS_SECRET`
- `CB_ACCESS_PASSPHRASE`

A signature and timestamp will be generated when needed by the library.
Parameters must be provided to the function using the exported struct.

```go
import (
	donutchallenge "github.com/willclarktech/donut-challenge"
	"fmt"
)

func main() {
	orderParams := donutchallenge.CoinbaseOrderParams{
		ProductID: "BTC-USD",
		Side:      donutchallenge.BUY,
		Size:      "0.03",
		Price:     "3200.11",
	}
	err := donutchallenge.PlaceOrder(orderParams)
	if err != nil {
		fmt.Println("There was an error placing the order:", err.Error())
	}
}
```

The task only specified that the order should be placed, hence this function
only returns an error. In any real-life context some data will be needed from
the response, e.g. the order ID. This would be extracted in the same way as for `GetCurrentRate`.

## Caveats

This is my third day using Go, so a number of issues remain, including:

1. I have not had time to research libraries/frameworks, including testing frameworks beyond the built-in `go test`. I assume there are some interesting options that would be worth exploring for testing and/or HTTP client functionality.
1. When specifying the type of the parameters for `PlaceOrder` I really wanted a union type, hence my workaround using type aliases such as `Currency`. There is probably a more idiomatic way of doing this.
1. `golint` raises warnings because of a lack of comments in `constants.go`. I see little value in adding comments for each of these constants, but perhaps that is a firm norm in the Go community. More likely it confirms my approach is not idiomatic.
1. I have structured the files (and code within those files) in a way that seems to make sense to me, but my views on project structure generally take a while to crystallise, so it is likely this could be improved.
1. The tests are fairly minimal. I have brief integration tests for the main exported actions, and unit tests for most of the smaller functions. Some of the test tables could be filled out more with extra edge cases etc. However, I found my usual (Node.js, FP-ish) approach to testing did not translate well at all so I would need more time to research testing philosophy in Go to write a suite I was happy with.
1. One of the problems I faced when testing was related to error propagation. Without the use of stubs/spies etc., or making my functions very pure and passing functions around a lot, it was difficult to test unhappy paths. In the case of creating a signature from a base64-encoded secret, I could pass an input which I knew would throw an error when run through the real function, but otherwise simulating errors seemed to require more work than I had time for during this challenge. So I have had to rely on manually confirming that I am always checking for errors and returning where appropriate.
1. Placing an order currently requires the necessary credentials to be provided via environmental variables. It would obviously be nice to have other options for securely providing these, such as from a file.
1. The parameters for placing an order have many logical interactions as described in Coinbase’s API documentation. This provides an opportunity for helpful client-side validation. E.g. if the `type` is `market` then either `size` or `funds` must be specified. If neither is specified, then there is no need to wait for an API response before informing the user that the request is invalid. As written, my `PlaceOrder` function delegates all validation to the API since there are a lot of interactions, which would take a lot of code/time for a fairly trivial implementation.
1. If this was a real-life case, I would spend more time looking at existing implementations such as [go-gdax](https://github.com/preichenberger/go-gdax) (or in fact just using an implementation like that). However, I figured a library with as much structure as that (e.g. for the client) was probably overkill for this assignment, and working things out on my own was better for learning at this early stage in any case.
1. Since many things were new for me, I had to do a lot of spiking. My commits would normally exhibit a more atomic TDD process.
1. My commits appear to be unverified! I’m trying out [Keybase](https://keybase.io/) and there’s some additional setup I need to do before my commits will show as verified on GitHub.
