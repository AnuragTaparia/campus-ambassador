package dbmanage

import(
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"ca/data"
	"ca/html"
	"strings"
)

func CreateNewUser(email string, password string, user_id string, isGoogleSignIn bool){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	var query string

	if isGoogleSignIn==true{
		query = fmt.Sprintf("INSERT INTO Users(Username,User_Id,Google_Signup) VALUES(%s,\"%s\",%t)",email,user_id,isGoogleSignIn)
	}else{
		query = fmt.Sprintf("INSERT INTO Users VALUES(%s,%s,\"%s\",%t)",email,password,user_id,isGoogleSignIn)
	}

	fmt.Println("QUERY: ",query)
	_,err = db.Exec(query)
	if err!=nil{
		fmt.Println("Error creating new user")
		panic(err)
	}

	query = fmt.Sprintf("INSERT INTO User_Profiles(User_Id,Email) VALUES(\"%s\",%s)",user_id,email)

	_,err = db.Exec(query)
	if err!=nil{
		fmt.Println("Error creating new user profile")
		panic(err)
	}
}

func CreateSession(user_id string, token string){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	fmt.Println("TOKEN",token)
	query := fmt.Sprintf("INSERT INTO Session_Tokens VALUES(\"%s\",\"%s\")",user_id,token)

	_,err = db.Exec(query)
	if err!=nil{
		fmt.Println("Error creating new session: ",err)
	}
}

func CheckGoogleSignup(user_id string) (bool,error){
	var check bool

	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	query := fmt.Sprintf("SELECT Google_Signup FROM Users where User_ID = \"%s\"",user_id)
	err = db.QueryRow(query).Scan(&check)
	if err!=nil{
		fmt.Println("Error checking google signup",err)
		return false,err
	}

	return check,nil
}

func CheckUserExists(user_id string, email string) bool{
	var user data.User
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	err = db.QueryRow("SELECT Username FROM Users where User_ID = "+user_id).Scan(&user.Email)
	if err != nil {
	    fmt.Println("Check User Error",err.Error()) // proper error handling instead of panic in your app
		return false
	}

	user.Email = fmt.Sprintf("\"%s\"",user.Email)
	fmt.Println("EMAIL",email,user.Email)
	if user.Email==email{
		return true
	}

	return false
}

func CheckUserExistsRegular(email string) (bool,string){
	var user data.User
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		panic(err)
	}

	fmt.Println("QUERY EMAIL ",email)
	email = strings.Replace(email,"\"","",-1)
	query := fmt.Sprintf("SELECT User_ID FROM Users where Username = \"%s\"",email)

	err = db.QueryRow(query).Scan(&user.UserID)
	if err != nil {
	    fmt.Println("Check User Error",err.Error()) // proper error handling instead of panic in your app
		return false,""
	}

	if user.UserID!=""{
		return true,user.UserID
	}

	return false,""
}

func RetrievePassword(email string) string{
	var user data.User
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return ""
	}

	err = db.QueryRow("SELECT Password FROM Users where Username = "+email).Scan(&user.Password)
	if err != nil {
	    fmt.Println("Check User Error",err.Error()) // proper error handling instead of panic in your app
		return ""
	}
	fmt.Println(user.Password)
	return user.Password
}

func GetUserIdFromSessionToken(token string) string{
	var user data.User
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return ""
	}

	token = fmt.Sprintf("\"%s\"",token)

	err = db.QueryRow("SELECT User_ID FROM Session_Tokens where Token = "+token).Scan(&user.UserID)
	if err != nil {
	    fmt.Println("Check User Error",err.Error()) // proper error handling instead of panic in your app
		return ""
	}

	return user.UserID
}

func GetUserIdFromEmail(email string) (string,error){
	var user_id string
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return "",err
	}

	email = fmt.Sprintf("\"%s\"",email)

	err = db.QueryRow("SELECT User_ID FROM Users where Username="+email).Scan(&user_id)
	if err != nil {
	    fmt.Println("Check User Error",err.Error())
		return "",err
	}

	return user_id,nil
}

func GetUserProfile(user_id string) html.ProfileParams{
	var user data.User
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return html.ProfileParams{}
	}

	query := fmt.Sprintf("SELECT Full_Name,Telephone,Address,Email FROM User_Profiles WHERE User_ID=\"%s\"",user_id)
	err = db.QueryRow(query).Scan(&user.Name,&user.Phone,&user.Address,&user.Email)
	if err!=nil{
		fmt.Println(err)
		return html.ProfileParams{}
	}

	profile := html.ProfileParams{
		Name: user.Name.String,
		Phone: user.Phone.String,
		Address: user.Address.String,
		Email: user.Email,
	}

	return profile
}

func GetEmailFromResetCode(code string) (string,error){
	var email string
	
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return "",err
	}

	query := fmt.Sprintf("SELECT Email FROM Reset WHERE Code=\"%s\"",code)
	err = db.QueryRow(query).Scan(&email)
	if err!=nil{
		fmt.Println("Error in GetEmailFromResetCode, ",err)
		return "",err
	}

	return email,nil

}

func DeleteInvalidResetCode(email string){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return 
	}

	query := fmt.Sprintf("DELETE FROM Reset WHERE Email=\"%s\"",email)
	_,err = db.Exec(query)
	if err!=nil{
		fmt.Println("Error creating new session: ",err)
	}

}

func RetrieveResetCode(email string) (string,error){
	var code string	

	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return "",err
	}

	query := fmt.Sprintf("SELECT Code FROM Reset WHERE Email=\"%s\"",email)
	err = db.QueryRow(query).Scan(&code)
	if err!=nil{
		return "",err
	}

	return code, nil
}

func SetResetCode(email,code string){
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return 
	}

	query := fmt.Sprintf("INSERT INTO Reset VALUES(\"%s\",\"%s\")",email,code)
	_,err = db.Exec(query)
	if err!=nil{
		fmt.Println(err)
		return
	}
}

func UpdatePassword(email,password string) error{
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return err
	}

	query := fmt.Sprintf("UPDATE Users SET Password=\"%s\" WHERE Username=\"%s\"",password,email)
	_,err = db.Exec(query)
	if err!=nil{
		fmt.Println(err)
		return err
	}

	return nil	
}

func EditProfile(user_id,fullname,phone,address,email string) bool{
	db,err := sql.Open("mysql","rishi:pwd1234@tcp(localhost:3306)/challengedb")
	defer db.Close()

	if err!=nil{
		fmt.Println(err)
		return false
	}

	allNull := true

	values := []string{fullname,phone,address,email}
	colNames := []string{"Full_Name","Telephone","Address","Email"}

	query := "UPDATE User_Profiles SET "
	for i:=0;i<len(values);i++{
		if values[i]!=""{
			values[i] = strings.Replace(values[i],`\`,`\\`,-1)
			values[i] = strings.Replace(values[i],"'",`\'`,-1)
			values[i] = strings.Replace(values[i],`"`,`\"`,-1)
			allNull=false
			query += colNames[i]+"="
			query += "\""+values[i]+"\""
			if i!=len(values)-1{
				query += ","
			}
		}
	}

	if allNull==true{
		return true
	}
	chars := []rune(query)
	if string(chars[len(chars)-1])==","{
		query = string(chars[:len(chars)-1])
	}
	query += fmt.Sprintf(" WHERE User_Id=\"%s\"",user_id)

	fmt.Println("QUERY: ",query)
	_,err = db.Exec(query)
	if err!=nil{
		fmt.Println("Error Updating profile", err)
		return false
	}

	return true
}