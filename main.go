package main

import "fmt"

func initialize() {

	// obtenho todas as variaveis globais

}

func main() {

	fmt.Println("============================================")
	fmt.Println("                  CONSUMER                  ")
	fmt.Println("============================================")
	fmt.Println("")

	// função recebe as dados do rabbit
	GetVideoDataFromBroker() // runing in a >>> go func() <<<

}
