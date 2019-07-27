package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
)

type Transaction struct {
	CCNum      string  `bson:"ccnum"`
	Date       string  `bson:"date"`
	Amount     float32 `bson:"amount"`
	Cvv        string  `bson:"cvv"`
	Expiration string  `bson:"exp"`
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Panicln(err)
	}
	defer session.Close()

	results := make([]Transaction, 0)
	if err := session.DB("store").C("transactions").Find(nil).All(&results); err != nil {
		log.Panicln(err)
	}
	for _, txn := range results {
		fmt.Println(txn.CCNum, txn.Date, txn.Amount, txn.Cvv, txn.Expiration)
	}
}
