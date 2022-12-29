package helper

import "fmt"

func PanicIfError(err error) {
	if err != nil {
		fmt.Println("Panic Called:", err)
		panic(err)
	}
}
