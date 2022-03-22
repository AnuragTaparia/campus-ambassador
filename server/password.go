package server

import (
  "fmt"
  "golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
  var passwordBytes = []byte(password)
  hashedPasswordBytes, err := bcrypt.
    GenerateFromPassword(passwordBytes, bcrypt.MinCost)

  return string(hashedPasswordBytes), err
}

func doPasswordsMatch(hashedPassword, currPassword string) bool {
  fmt.Println("HASHED: ",hashedPassword)
  fmt.Println("CURRENT: ",currPassword)
  err := bcrypt.CompareHashAndPassword(
    []byte(hashedPassword), []byte(currPassword))
  return err == nil
}

func EncryptPassword(password string) (string,error){
  // Hash password
  var hashedPassword, err = hashPassword(password)

  if err != nil {
    println(fmt.Println("Error hashing password"))
    return "",err
  }

  fmt.Println("Password Hash:", hashedPassword)
  return hashedPassword,nil

}

func DecryptPassword(password string, hashedPassword string) bool{
	return doPasswordsMatch(hashedPassword, password)	
}