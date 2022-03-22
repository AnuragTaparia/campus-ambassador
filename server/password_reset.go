package server

import (

	"time"
	"math/rand"
	"ca/dbmanage"
)

func GenerateCode() string{
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 5)
	for i := range b {
	    b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CheckValidity(email string){
	start := time.Now().Unix()
	go DeleteIfTimeout(start,email)
}

func DeleteIfTimeout(start int64,email string){
	for{
		now := time.Now().Unix()
		if now-start >= int64(300){
			dbmanage.DeleteInvalidResetCode(email)
			break
		}
	}
}