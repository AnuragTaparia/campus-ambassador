function logout(){
	document.cookie = "golang_challenge_session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
	window.location.href = "/"
}