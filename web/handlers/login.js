import {login, register} from "./../js/components/auth/auth.js";
import { auth } from "./../js/pages/auth.js";
import { fetchLogin, fetchRegister } from "./../js/api/auth.js";




export function handleAuth()  {	
	const app = document.getElementById("app")
	app.innerHTML = auth(login);
	authListener(app);
};


export default function authListener(app) {
    const loginTab = document.getElementById("loginTab");
    const registerTab = document.getElementById("registerTab");

    loginTab?.addEventListener("click", () => {
		console.log("called? ");
        app.innerHTML = auth(login);
		authListener(app);
    });

    registerTab?.addEventListener("click", () => {	
        app.innerHTML = auth(register);
		authListener(app);
    });

    const loginForm = document.getElementById("loginForm");
    loginForm?.addEventListener("submit", loginSubmit);

    const registerForm = document.getElementById("registerForm");
    registerForm?.addEventListener("submit", registerSubmit);
}

function errorText(body) {
    if (typeof body == "string" && body.trim()) return body.trim();
    return "Something went wrong. Please try again.";
}

async function loginSubmit(e) {
    e.preventDefault();

    const form = e.currentTarget;
    const submit = form.querySelector(".auth-submit");
    const status = form.querySelector(".auth-error");
    status.textContent = "";

    const data = {
        identifier: document.getElementById("loginIdentifier").value.trim(),
        password: document.getElementById("loginPassword").value,
    };

    if (!data.identifier || !data.password) {
        status.textContent = "Enter your nickname/email and password.";
        return;
    }

    submit.disabled = true;
    try {
        const response = await fetchLogin(data);
        if (response.ok) {
           // try get in home again 
			console.log("good job")
            return;
        }
        status.textContent = errorText(response.body);
    } finally {
        submit.disabled = false;
    }
}

async function registerSubmit(e) {
    e.preventDefault();

    const form = e.currentTarget;
    const submit = form.querySelector(".auth-submit");
    const status = form.querySelector(".auth-error");
    status.textContent = "";

    const data = {
        nick_name: document.getElementById("nickname").value.trim(),
        last_name: document.getElementById("lastName").value.trim(),
        first_name: document.getElementById("firstName").value.trim(),
        email: document.getElementById("email").value.trim(),
        password: document.getElementById("password").value,
        gender: document.getElementById("gender").value,
        age: Number(document.getElementById("age").value),
    };

    if (!data.nick_name || !data.first_name || !data.last_name || !data.email || !data.password) {
        status.textContent = "Please fill in all the required fields.";
        return;
    }

    submit.disabled = true;
    try {
        const response = await fetchRegister(data);
        if (response.ok) {
           	// ok ??  
			console.log("good jbob")
            return;
        }
        status.textContent = errorText(response.body);
    } finally {
        submit.disabled = false;
    }
}
