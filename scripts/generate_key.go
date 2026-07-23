package main

import (
	"fmt"

	"github.com/kev1nandreas/go-rest-api-template/src/shared/auth"
)

func main() {
	fmt.Println(auth.GenerateRandomKey())
}
