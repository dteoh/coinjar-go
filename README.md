# coinjar-go

[![Build Status](https://travis-ci.org/dteoh/coinjar-go.svg?branch=master)](https://travis-ci.org/dteoh/coinjar-go)

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

## APIs Implemented

* Account

        client.Account()

* Bitcoin Addresses
    * List

            client.BitcoinAddresses() // Only first 100
            client.ListBitcoinAddresses(limit, offset int)

    * Retrieve

            client.BitcoinAddress(address string)

* Contacts
    * List

            client.Contacts() // Only first 100
            client.ListContacts(limit, offset int)

    * Retrieve

            client.Contact(uuid string)

* Payments
    * List

            client.Payments() // Only first 100
            client.ListPayments(limit, offset int)

    * Retrieve

            client.Payment(uuid string)

* Transactions
    * List

            client.Transactions() // Only first 100
            client.ListTransactions(limit, offset int)

    * Retrieve

            client.Transaction(uuid string)

* Fair Rate

        client.FairRate(currency string)

## TODOs

* Implement missing APIs
* Make collection APIs automatically retrieve all records

