package main

import (
	"fmt"
	"log"
	"os"

	client "github.com/living-etc/sightseer.go/clients/git"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	gitClient := client.NewGitClient(wd)

	main, err := gitClient.Show("main")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hash:    ", main.Hash)
	fmt.Println("Subject: ", main.Subject)
	fmt.Println("Body:    ", main.Body)
	fmt.Println("Author:  ", main.AuthorName)
	fmt.Println("Date:    ", main.AuthorDate)
}
