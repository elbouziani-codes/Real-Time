import { login, register } from "./../components/login.js";
import {postLogin, postRegister }from  "./../api/auth.js"
import { navigate } from "./../router/router.js";

function loginAndRegister() {
    if (window.location.pathname == "/register") {
        return register();
    } else if (window.location.pathname == "/login") {
        return login();
    } else {
        navigate("/");
    }
}

export function auth() {
    return `
        <section class="auth">

        <div class="auth-left">
            <h1>💬 Real-Time Forum</h1>
            <p>Join discussions with developers around the world. Connect, share, and learn in real-time.</p>
        </div>

        <div class="auth-right">

            <div class="card">
                ${loginAndRegister()}
            </div>

        </div>
    </section>
    `;
}

export function authListener() {
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

    let ok = postLogin(data);
    console.log(ok)
    console.log("LOGIN", data);
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
    let ok = postRegister(data)
    console.log(ok)
    console.log("REGISTER", data);
}
