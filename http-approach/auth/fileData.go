package auth

import (
	"fmt"
	"log"
	"os"
)

func FileData() (*os.File, *os.File) {

	//open the files in the directory data
	Userfile, err := os.Open("auth/data/users.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	Membersfile, err := os.OpenFile("auth/data/members.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	return Userfile, Membersfile
}
