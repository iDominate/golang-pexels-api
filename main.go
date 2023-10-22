package main

import (
	"fmt"
	"log"

	c "github.com/iDominate/golang-pexels-api/Client"
	"github.com/iDominate/golang-pexels-api/env"
)

const PEXELS_TOKEN = "PEXELS_TOKEN"

func main() {
	env.InitEnv()
	client := &c.Client{}
	token, err := env.Getenv(PEXELS_TOKEN)

	if err != nil {
		log.Fatal(err.Error())
	}
	client = client.NewClient(token)
	res, err := client.GetRandomVideo()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(res)

}
