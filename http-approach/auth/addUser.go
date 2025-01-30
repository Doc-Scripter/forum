package auth

import (
	"fmt"
	"os"
)

func AddUser(username, password string) bool {
	file, err := os.OpenFile("auth/data/users.txt", os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("u-%s p-%s\n", username, password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
