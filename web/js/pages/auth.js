import { login, register } from "../components/auth/auth.js";
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

