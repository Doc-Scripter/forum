package auth

import (
	"bufio"
	"log"
	"strings"
)

func CredentialExists(credential string) bool {
	Userfile, _ := FileData()
	scanner := bufio.NewScanner(Userfile)

	for scanner.Scan() {
		line := scanner.Text()
		u := strings.Split(line, " ")[0][2:]
		if credential == u {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	return false
}
