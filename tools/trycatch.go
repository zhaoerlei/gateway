package tools

import (
	"log"
)

func Safely(errMsg string) {
	if err := recover(); err != nil {
		log.Println(errMsg, "###########################", err)
	}
}
