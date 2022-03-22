package server

import(
	"net/http"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"ca/html"
	"ca/dbmanage"
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/google/uuid"
	"time"
	"github.com/gomodule/redigo/redis"
	"strings"
	"errors"
)

type googleAuthResp struct{
	UserID string `json:"id"`
	Email string `json:"email"`
}

var cache redis.Conn

var (
	googleOauthConfigSignUp = &oauth2.Config{
		RedirectURL: "http://34f5-167-172-142-10.ngrok.io/signup",
		ClientID: os.Getenv("Google_Oauth_Client_Id"),
		ClientSecret: os.Getenv("Google_Oauth_Client_Secret"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	googleOauthConfigLogin = &oauth2.Config{
		RedirectURL: "http://34f5-167-172-142-10.ngrok.io/login",
		ClientID: os.Getenv("Google_Oauth_Client_Id"),
		ClientSecret: os.Getenv("Google_Oauth_Client_Secret"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	randomState = "random"
)


func Begin(){
	initCache()

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/userlogin",handleUserLoginPage)
	mux.HandleFunc("/googlesignup",handleGoogleSignup)
	mux.HandleFunc("/googlelogin",handleGoogleLogin)
	mux.HandleFunc("/signup",handleCallbackSignup) //google auth callback
	mux.HandleFunc("/login",handleCallbackLogin) //google auth callback
	mux.HandleFunc("/regularlogin",handleRegularLogin)
	mux.HandleFunc("/regularsignup",handleRegularSignup)
	mux.HandleFunc("/editprofile",handleEditProfilePage)
	mux.HandleFunc("/mainprofile",handleMainProfilePage)
	mux.HandleFunc("/edit",handleEditProfile)
	mux.HandleFunc("/forgotpwd",handleForgotPwdPage)
	mux.HandleFunc("/requestcode",handleRequestCode)
	mux.HandleFunc("/resetpassword",handleResetPassword)

	mux.Handle("/css/login.css",http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))	
	mux.Handle("/js/login.js",http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	mux.Handle("/css/profile.css",http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	mux.Handle("/js/profile.js",http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))


	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080",handler)
}

func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	cache = conn
}


func handleUserLoginPage(w http.ResponseWriter, r *http.Request) {

	cookieExists,sessionToken,err := checkCookieExists(w,r)
	if err!=nil{
		fmt.Println(err)
		return
	}

	if cookieExists==false{

		fmt.Println("Cookie Not Found")
		p := html.LoginParams{
			Title: "Login/Signup",
		}
		html.Login(w, p)
		return
	}


	// fmt.Println("Cookie Found", sessionToken)
	user_id := dbmanage.GetUserIdFromSessionToken(sessionToken)
	
	p := dbmanage.GetUserProfile(user_id)
	p.Title = "Profile"
	html.MainProfile(w,p)
}

func handleForgotPwdPage(w http.ResponseWriter, r *http.Request){
	p := html.ResetPageParams{
		Title: "Reset Password",
	}
	html.ResetMainPage(w, p)
}

func handleResetPassword(w http.ResponseWriter, r *http.Request){
	if r.Method==http.MethodPost{

		r.ParseForm()

		email := r.FormValue("email")
		code := r.FormValue("code")
		password := r.FormValue("password")

		trueEmail,err := dbmanage.GetEmailFromResetCode(code)
		if err!=nil{
			serveLogin("Login/Signup","Something went wrong on the server. Please try again",w)
			return
		}
		fmt.Println(trueEmail,email)
		if trueEmail!=email{
			serveLogin("Login/Signup","Unauthorised attempt to change passwords",w)
			return
		}

		user_id,err := dbmanage.GetUserIdFromEmail(email)
		if err!=nil{
			serveLogin("Login/Signup","Something went wrong on the server. Please try again",w)
			return
		}

		isGoogleSignup,err := dbmanage.CheckGoogleSignup(user_id)
		if isGoogleSignup{
			serveLogin("Login/Signup","Cannot change passwords for a account signed up with Google",w)
			return
		}

		dbCode,err := dbmanage.RetrieveResetCode(email)
		if err!=nil{
			fmt.Println(err)
			return
		}
		fmt.Println("DBCODE: ",dbCode,code)
		if dbCode!=code{
			dbmanage.DeleteInvalidResetCode(email)
			serveLogin("Login/Signup","Verification codes do not match",w)
			return
		}

		encPwd,err := EncryptPassword(password)
		if err!=nil{
			fmt.Println(err)
			return
		}

		err = dbmanage.UpdatePassword(email,encPwd)
		if err!=nil{
			fmt.Println(err)
			serveLogin("Login/Signup","Updae password failed",w)
			return
		}

		dbmanage.DeleteInvalidResetCode(email)
		serveLogin("Login/Signup","Password has been reset",w)
		return

	}else{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleRequestCode(w http.ResponseWriter, r *http.Request){
	if r.Method==http.MethodPost{

		r.ParseForm()

		email := r.FormValue("email")
		userExists,_ := dbmanage.CheckUserExistsRegular(email)

		if userExists==false{
			serveLogin("Login/Signup","That is not a registered email address",w)
			return
		}

		user_id,err := dbmanage.GetUserIdFromEmail(email)
		if err!=nil{
			serveLogin("Login/Signup","Something went wrong on the server. Please try again",w)
			return
		}

		isGoogleSignup,err := dbmanage.CheckGoogleSignup(user_id)
		if isGoogleSignup{
			serveLogin("Login/Signup","Cannot change passwords for a account signed up with Google",w)
			return
		}

		code := GenerateCode()

		err = SendCode(code,email)
		if err!=nil{
			fmt.Println("Error in sending email: ",err)
			return
		}
		dbmanage.SetResetCode(email,code)
		CheckValidity(email)
		p := html.ResetPageParams{
			Title: "Reset Password",
			Email: email,
			Message: "An email with a verification code has been sent to you",
		}
		html.ResetPassword(w, p)

	}else{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleEditProfile(w http.ResponseWriter, r *http.Request){
	if r.Method==http.MethodPost{
		cookieExists,sessionToken,err := checkCookieExists(w,r)
		if err!=nil{
			return
		}

		if cookieExists==false{
			serveLogin("Login/Signup","",w)
			return
		}

		user_id := dbmanage.GetUserIdFromSessionToken(sessionToken)
		r.ParseForm()

		email := r.FormValue("email")
		fullname := r.FormValue("fullname")
		address := r.FormValue("address")
		phone := r.FormValue("phone")

		success:=dbmanage.EditProfile(user_id,fullname,phone,address,email)
		if success{
			p := dbmanage.GetUserProfile(user_id)
			p.Title = "Profile"
			html.MainProfile(w,p)
		}
	}

}

func handleEditProfilePage(w http.ResponseWriter, r *http.Request){
	cookieExists,sessionToken,err := checkCookieExists(w,r)
	if err!=nil{
		return
	}

	if cookieExists==false{
		serveLogin("Login/Signup","",w)
		return
	}else{
		user_id := dbmanage.GetUserIdFromSessionToken(sessionToken)
		googleSignedIn,_ := dbmanage.CheckGoogleSignup(user_id)

		p := dbmanage.GetUserProfile(user_id)
		p.Title = "Edit Profile"
		p.EditEmail = googleSignedIn
		html.EditProfile(w,p)
		return
	}
}
func handleMainProfilePage(w http.ResponseWriter, r *http.Request){
	cookieExists,sessionToken,err := checkCookieExists(w,r)
	if err!=nil{
		return
	}

	if cookieExists==false{
		serveLogin("Login/Signup","",w)
		return
	}else{
		user_id := dbmanage.GetUserIdFromSessionToken(sessionToken)
		p := dbmanage.GetUserProfile(user_id)
		p.Title = "Profile"
		html.MainProfile(w,p)
		return
	}
}

func handleGoogleSignup(w http.ResponseWriter, r *http.Request){
	url := googleOauthConfigSignUp.AuthCodeURL(randomState)
	http.Redirect(w,r,url, http.StatusTemporaryRedirect)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request){
	url := googleOauthConfigLogin.AuthCodeURL(randomState)
	http.Redirect(w,r,url, http.StatusTemporaryRedirect)
}


func handleRegularLogin(w http.ResponseWriter, r *http.Request){
	if r.Method==http.MethodPost{
		r.ParseForm()

		email := r.FormValue("email")
		email = fmt.Sprintf("\"%s\"",email)
		pwd := r.FormValue("pwd")

		userExists,user_id := dbmanage.CheckUserExistsRegular(email)
		if userExists{
			dbPwd := dbmanage.RetrievePassword(email)
			pwdMatch := DecryptPassword(pwd,dbPwd)

			if pwdMatch{
				createCookie(user_id,w)
				p:=dbmanage.GetUserProfile(user_id)
				p.Title = "Profile"

				html.MainProfile(w,p)
			}else{
				serveLogin("Login/Signup","Incorrect Credentials",w)
			}

			return
		}else{
			// Account User Does not exist
			serveLogin("Login/Signup","This account does not exist",w)
		}

		return
	}else{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleRegularSignup(w http.ResponseWriter, r *http.Request){
	if r.Method==http.MethodPost{

		r.ParseForm()

		email := r.FormValue("email")
		email = fmt.Sprintf("\"%s\"",email)
		pwd := r.FormValue("pwd")

		userExists,_ := dbmanage.CheckUserExistsRegular(email)
		if userExists{
			serveLogin("Login/Signup","This account already exists",w)
			return
		}else{
			encPwd,err := EncryptPassword(pwd)
			if err!=nil{
				fmt.Println(err)
				return
			}
			user_id := uuid.New().String()
			user_id = strings.Replace(user_id,"-","",-1)

			
			encPwd = fmt.Sprintf("\"%s\"",encPwd)

			dbmanage.CreateNewUser(email,encPwd,user_id,false)
			createCookie(user_id,w)	

			// Replace with edit profile page
			googleSignedIn,_ := dbmanage.CheckGoogleSignup(user_id)
			p := dbmanage.GetUserProfile(user_id)
			p.Title = "Edit Profile"
			p.EditEmail = googleSignedIn
			html.EditProfile(w,p)
		}

	}else{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleCallbackLogin(w http.ResponseWriter, r *http.Request){
	authResp := makeGoogleAuthRequests(w,r,googleOauthConfigLogin)	

	user_id := authResp.UserID
	email := authResp.Email
	// password := ""
	// isGoogleSignin := true

	email = fmt.Sprintf("\"%s\"",email)

	userExists := dbmanage.CheckUserExists(user_id,email)

	if userExists{
		createCookie(user_id,w)
		
		p:=dbmanage.GetUserProfile(user_id)
		p.Title = "Profile"
		
		html.MainProfile(w,p)
	}else{
		serveLogin("Login/Signup","This account does not exist",w)
	}

	// http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
}

func handleCallbackSignup(w http.ResponseWriter, r *http.Request){
	authResp := makeGoogleAuthRequests(w,r,googleOauthConfigSignUp)	

	user_id := authResp.UserID
	email := authResp.Email
	password := ""
	isGoogleSignin := true

	email = fmt.Sprintf("\"%s\"",email)
	
	userExists := dbmanage.CheckUserExists(user_id,email)
	if userExists{
		serveLogin("Login/Signup","That account already exists",w)
		return
	}

	dbmanage.CreateNewUser(email,password,user_id,isGoogleSignin)

	createCookie(user_id,w)
	
	// Replace with edit profile page
	googleSignedIn,_ := dbmanage.CheckGoogleSignup(user_id)
	p := dbmanage.GetUserProfile(user_id)
	p.Title = "Edit Profile"
	p.EditEmail = googleSignedIn
	html.EditProfile(w,p)
	// http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
}

func createCookie(user_id string, w http.ResponseWriter){
	sessionToken := uuid.New().String()
	sessionToken = strings.Replace(sessionToken,"-","",-1)

	fmt.Println("Session Token: ", sessionToken)

	_, err := cache.Do("SETEX", sessionToken, "1200", user_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "golang_challenge_session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(1200* time.Second),
	})

	dbmanage.CreateSession(user_id,sessionToken)
}

func makeGoogleAuthRequests(w http.ResponseWriter, r *http.Request,googleOauthConfig *oauth2.Config) googleAuthResp{
	if r.FormValue("state") != randomState{
		fmt.Println("State not valid")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return googleAuthResp{}
	}

	token,err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err!=nil{
		fmt.Println("Could not get token: ",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return googleAuthResp{}
	}

	resp,err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken)

	if err!=nil{
		fmt.Println("Could not create request: ",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return googleAuthResp{}
	}

	defer resp.Body.Close()
	content,err := ioutil.ReadAll(resp.Body)
	
	if err!=nil{
		fmt.Println("Could not parse response: ",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return googleAuthResp{}
	}

	var authResp googleAuthResp

	err = json.Unmarshal(content, &authResp)

	if err!=nil{
		fmt.Println("Could not JSON parse response: ",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return googleAuthResp{}
	}

	return authResp
}

func checkCookieExists(w http.ResponseWriter, r *http.Request) (bool,string,error){
	c,err := r.Cookie("golang_challenge_session_token")
	if err != nil{
		if err==http.ErrNoCookie{
			return false, "",nil
		}

		w.WriteHeader(http.StatusBadRequest)
		return false, "", errors.New("Bad Request Error")
	}

	sessionToken := c.Value

	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		fmt.Println("Fetching Cache failed ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return false,"",errors.New("Internal Server Error")
	}

	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return false,"",errors.New("Unauthorised request")
	}

	return true,sessionToken,nil
}

func serveLogin(title,msg string,w http.ResponseWriter){
	p := html.LoginParams{
		Title: title,
		Message: msg,
	}
	html.Login(w, p)
}