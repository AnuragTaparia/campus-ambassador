package dbmanage

import(
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
)

func Begin(){
	err := createDatabase()
	if err!=nil{
		return
	}

	createUsersTable()
	createUserProfileTable()
	createSessionTokenTable()
	createResetPasswordCode()
}

func createDatabase() error{
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return errors.New("Something went wrong connecting to mysql. Please check the logs")
	}

	_,err = db.Exec("CREATE DATABASE IF NOT EXISTS challengedb")
	if err!=nil{
		fmt.Println(err)
		return errors.New("Something went wrong creating the database. Please check the logs")
	}

	fmt.Println("Database Created")
	return nil
}

func createUsersTable(){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	_,err = db.Exec("CREATE TABLE IF NOT EXISTS Users(Username varchar(50), Password text, User_ID varchar(50),Google_Signup Boolean default false ,PRIMARY KEY (Username))")
	if err!=nil{
		panic(err)
	}
}

func createUserProfileTable(){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	_,err = db.Exec("CREATE TABLE IF NOT EXISTS User_Profiles(User_ID varchar(50), Full_Name text, Telephone varchar(15),Address text,Email varchar(50), FOREIGN KEY(Email) REFERENCES Users(Username) ON DELETE CASCADE ON UPDATE CASCADE )")
	if err!=nil{
		panic(err)
	}

}

func createResetPasswordCode(){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	_,err = db.Exec("CREATE TABLE IF NOT EXISTS Reset(Email varchar(50), Code varchar(10))")	
	if err!=nil{
		panic(err)
	}
}

func createSessionTokenTable(){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	_,err = db.Exec("CREATE TABLE IF NOT EXISTS Session_Tokens(User_ID varchar(50), Token text)")
	if err!=nil{
		panic(err)
	}
}