# coinjar-go

An unofficial golang package to interact with the CoinJar API.


## Usage

    package main

    import (
    	"fmt"
    	"github.com/dteoh/coinjar-go/coinjar"
    )

    func main() {
    	client := coinjar.NewClient("your api key")
    	account, _ := client.Account()
    	fmt.Println(account)
    }

