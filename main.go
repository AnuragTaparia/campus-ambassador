package main

import(
	"ca/dbmanage"
	"ca/server"
)

func main(){
	dbmanage.Begin()
	server.Begin()
}