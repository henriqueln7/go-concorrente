package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//Crie um programa que recebe um número inteiro n como argumento e cria n goroutines.
//Cada uma dessas goroutines deve dormir por um tempo aleatório de no máximo 5 segundos.
//A main-goroutine deve esperar todas as goroutines filhas terminarem de executar para em seguida escrever na saída padrão o valor de n.

// Função auxiliar pra pegar o ID de uma Goroutine. Peguei em algum canto aleatório da internet
func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("É necessário passar a quantidade de GoRoutines a serem criadas")
		fmt.Println("Ex.:: go run fork-sleep-join 10")
		panic("Faltou o número de GoRoutines")
	}
	n, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	joinCh := make(chan int, n)

	rand.Seed(time.Now().Unix())

	//Criando as n goroutines
	for i := 0; i < n; i++ {
		go func(joinCh chan int) {
			timeToSleep := rand.Intn(6)
			fmt.Printf("[%d] Dormindo por %d segundos \n", goid(), timeToSleep)
			time.Sleep(time.Duration(timeToSleep) * time.Second)
			fmt.Printf("[%d] Acordei\n", goid())
			joinCh <- 1
		}(joinCh)
	}

	for i := 0; i < n; i++ {
		<-joinCh
	}

	fmt.Println(n)
}
