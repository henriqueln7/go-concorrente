package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

//Crie um programa organizado como um pipeline de goroutines.
//Esse programa deve receber como argumento um caminho absoluto para um diretório.
//Uma goroutine deve navegar na árvore que tem como raiz o diretório passado como argumento.
//Essa goroutine deve passar para uma próxima goroutine do pipeline o nome dos arquivos encontrados na busca dos diretórios, ou seja, ignore os diretórios.
//Esta segunda goroutine deve ler o primeiro byte de conteúdo de cada um desses arquivos e escrever na saída padrão o nome dos arquivos que tem esse valor do primeiro byte sendo par.

func readDir(path string, filesOutChannel chan<- string) {
	err := filepath.WalkDir(path, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !dirEntry.IsDir() {
			filesOutChannel <- path
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	close(filesOutChannel)

}

func readFirstByte(filesInChannel <-chan string, joinCh chan int) {

	for filePath := range filesInChannel {
		file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("PUUF")
		}
		if len(file) > 0 {
			fmt.Printf("File: %s | First Byte: %b\n", filePath, file[0])
		}
	}

	joinCh <- 1
}

func main() {
	args := os.Args
	var path string

	if len(args) == 1 {
		path = "."
	} else {
		path = args[1]
	}

	start := time.Now()

	filesChannel := make(chan string, 200)
	joinCh := make(chan int)

	go readDir(path, filesChannel)
	go readFirstByte(filesChannel, joinCh)


	elapsed := time.Since(start)

	<-joinCh
	fmt.Println("Time elapsed: ", elapsed)
}
