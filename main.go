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
	GetYoutubeData() // runing in a >>> go func() <<<

	// tratará os dados
	// -- consultando black list (localizado em banco de dados local)
	// -- anexando categoria (localizado em banco de dados local)

	// envio via API para o endpoint informado
	// verifico se aguardo retorno ou não ( a definir )

	// outras funções
	// +++++++++
	// obter periódicamente a lista de categorias do youtube.
	// -- A função deverá enviar dados pra o rabbit para que o coletor busque os dados

	//
}
