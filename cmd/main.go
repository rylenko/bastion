package main

import (
	"fmt"
	"log"

	"github.com/rylenko/sapphire/pkg/shield"
	"github.com/rylenko/sapphire/pkg/shieldprovider"
)

func main() {
	sender, err := shield.NewSender(shieldprovider.New(), nil, nil)
	if err != nil {
		log.Fatal("new shield sender: ", err)
	}

	fmt.Println(sender)
}
