window.addEventListener('load',e=>{
	const signUpButton = document.getElementById('signUp');
	const signInButton = document.getElementById('signIn');
	const container = document.getElementById('container');

	signUpButton.addEventListener('click', () => {
		container.classList.add("right-panel-active");
	});

	signInButton.addEventListener('click', () => {
		container.classList.remove("right-panel-active");
	});
})

function confirmPassword(){
	var password = document.getElementById('password')
	var confirmPwd = document.getElementById('confirmPwd')
	var signupBtn = document.getElementById('signupBtn')
	var confirmMsg = document.getElementById('confirmMsg')

	if(password.value!=""){
		confirm = confirmPwd.value
		if(confirm==password.value){
			signupBtn.disabled = false
			confirmMsg.innerText = ""
		}
		else{
			confirmMsg.innerText = "Passwords do not match"
			signupBtn.disabled = true
		}
	}
	else{
		signupBtn.disabled = true
	}
}