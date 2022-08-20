package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

//Crie um programa que recebe um número inteiro n como argumento e cria n goroutines.
//Cada uma dessas goroutines deve dormir por um tempo aleatório de no máximo 5 segundos.
//Depois que  acordar, cada thread deve sortear um outro número aleatório s (entre 0 e 10).
//Somente depois de todas as n goroutines terminarem suas escolhas (ou seja, ao fim da primeira fase), começamos a segunda fase.
//Nesta segunda fase, a n-ésima goroutine criada deve dormir pelo tempo s escolhido pela goroutine n - 1 (faça a contagem de maneira modular, ou seja, a primeira goroutine dorme conforme o número sorteado pela última).

func worker(id int, joinCh chan int, chIn <-chan int, chOut chan<- int, wg *sync.WaitGroup) {
	timeToSleep := rand.Intn(6)
	fmt.Printf("[%d] Dormindo por %d segundos \n", id, timeToSleep)
	time.Sleep(time.Duration(timeToSleep) * time.Second)
	fmt.Printf("[%d] Acordei\n", id)

	wg.Done()
	wg.Wait()

	timeToNextGoroutineSleep := rand.Intn(11)
	fmt.Printf("[%d] Vou colocar a próxima goroutine pra dormir %d segundos :D\n", id, timeToNextGoroutineSleep)
	chOut <- timeToNextGoroutineSleep
	sleepFor := <-chIn

	fmt.Printf("[%d]Vou dormir por %d segundos \n", id, sleepFor)
	time.Sleep(time.Duration(sleepFor) * time.Second)
	fmt.Printf("[%d] Acordei :D\n", id)
	joinCh <- 1
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

	wg := sync.WaitGroup{}

	chans := make([]chan int, n)

	for i := 0; i < len(chans); i++ {
		chans[i] = make(chan int, 1)
	}

	wg.Add(n) //Barreira para segurar n coroutines
	for i := 0; i < n; i++ {

		previousIndex := (i + n - 1) % n
		nextIndex := (i + 1) % n

		fmt.Printf("[%d] previousIndex: %d\n", i, previousIndex)
		fmt.Printf("[%d] nextIndex: %d\n", i, nextIndex)

		go worker(i, joinCh, chans[nextIndex], chans[previousIndex], &wg) //?????????
	}

	for i := 0; i < n; i++ {
		<-joinCh
	}

	fmt.Println(n)
}
