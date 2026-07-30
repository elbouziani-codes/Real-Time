const loginTab = document.getElementById("loginTab");
const registerTab = document.getElementById("registerTab");

const loginForm = document.getElementById("loginForm");
const registerForm = document.getElementById("registerForm");

loginTab.addEventListener("click", () => {

    loginTab.classList.add("active");
    registerTab.classList.remove("active");

    loginForm.classList.remove("hidden");
    registerForm.classList.add("hidden");
});

registerTab.addEventListener("click", () => {

    registerTab.classList.add("active");
    loginTab.classList.remove("active");

    registerForm.classList.remove("hidden");
    loginForm.classList.add("hidden");
});

loginForm.addEventListener("submit", (e) => {

    e.preventDefault();

    const data = {

        identifier:
            document.getElementById("loginIdentifier").value,

        password:
            document.getElementById("loginPassword").value
    };

    console.log("LOGIN", data);

    // fetch("/api/login")
});

registerForm.addEventListener("submit", (e) => {

    e.preventDefault();

    const data = {

        nickname:
            document.getElementById("nickname").value,

        firstName:
            document.getElementById("firstName").value,

        lastName:
            document.getElementById("lastName").value,

        age:
            document.getElementById("age").value,

        gender:
            document.getElementById("gender").value,

        email:
            document.getElementById("email").value,

        password:
            document.getElementById("password").value
    };

    console.log("REGISTER", data);

    // fetch("/api/register")
});