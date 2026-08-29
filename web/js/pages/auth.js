import { navigate } from "./../router/router.js";


export function auth(fn) {
    return `
        <section class="auth">

        <div class="auth-left">
            <h1>💬 Real-Time Forum</h1>
            <p>Join discussions with developers around the world. Connect, share, and learn in real-time.</p>
        </div>

        <div class="auth-right">

            <div class="card">
                ${fn()}
            </div>

        </div>
    </section>
    `;
}

