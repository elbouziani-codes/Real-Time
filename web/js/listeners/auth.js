
import {navigate} from "./../router/router.js"
import {fetchLogin, fetchRegister }from  "./../api/auth.js"
export default function authListener() {
    const loginTab = document.getElementById("loginTab");

    const registerTab = document.getElementById("registerTab");

    loginTab?.addEventListener("click", () => {
    navigate("/login");
    });

    registerTab?.addEventListener("click", () => {
    navigate("/register");
    });

    const loginForm = document.getElementById("loginForm");

    loginForm?.addEventListener("submit", loginSubmit);

    const registerForm = document.getElementById("registerForm");

    registerForm?.addEventListener("submit", registerSubmit);
}





function loginSubmit(e) {
    e.preventDefault();

    const data = {
    identifier: document.getElementById("loginIdentifier").value,

    password: document.getElementById("loginPassword").value,
    };

    let response = fetchLogin(data);
}

function registerSubmit(e) {
    e.preventDefault();

    const data = {
    nick_name: document.getElementById("nickname").value,

    last_name: document.getElementById("lastName").value,

    first_name: document.getElementById("firstName").value,

    email: document.getElementById("email").value,

    password: document.getElementById("password").value,

    gender: document.getElementById("gender").value,

    age: Number(document.getElementById("age").value, 10),

};
    let ok = fetchRegister(data)
}
