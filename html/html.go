package html

import(
	"embed"
	"io"
	"html/template"
)

var files embed.FS

var(
	login = parse("static/login.html")
	mainProfile = parse("static/mainprofile.html")
	editProfile = parse("static/editProfile.html")
	resetMainPage = parse("static/resetpasswordhome.html")
	resetPassword = parse("static/resetpassword.html")
)

type LoginParams struct {
	Title string
	Message string
}

type ResetPageParams struct{
	Title string
	Email string
	Message string
}

func Login(w io.Writer, p LoginParams) error{
	return login.Execute(w, p)
}

type ProfileParams struct {
	Title string
	Name string
	Address string
	Phone string
	Email string
	EditEmail bool
}

func ResetMainPage(w io.Writer, p ResetPageParams) error{
	return resetMainPage.Execute(w, p)
}
func ResetPassword(w io.Writer, p ResetPageParams) error{
	return resetPassword.Execute(w, p)
}

func MainProfile(w io.Writer, p ProfileParams) error{
	return mainProfile.Execute(w, p)
}

func EditProfile(w io.Writer, p ProfileParams) error{
	return editProfile.Execute(w, p)
}


func parse(file string) *template.Template {
	return template.Must(template.ParseFiles(file))
}
