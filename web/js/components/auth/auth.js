function tabs(activeTab) {
    return `
        <div class="tabs">
            <button id="loginTab" class="tab-btn${activeTab == "login" ? " active" : ""}">Login</button>
            <button id="registerTab" class="tab-btn${activeTab == "register" ? " active" : ""}">Register</button>
        </div>
    `;
}

function field(type, id, label, placeholder, { required = true, minlength = "", autocomplete = "" } = {}) {
    const minAttr = minlength ? ` minlength="${minlength}"` : "";
    const autoAttr = autocomplete ? ` autocomplete="${autocomplete}"` : "";
    return `
        <label class="sr-only" for="${id}">${label}</label>
        <input type="${type}" class="auth-input" placeholder="${placeholder}" id="${id}" aria-label="${label}"${required ? " required" : ""}${minAttr}${autoAttr}>
    `;
}

function statusLine() {
    return `<p class="auth-error" role="alert" aria-live="polite"></p>`;
}

export function login() {
    return `
        ${tabs("login")}
        <!-- LOGIN FORM -->
        <form id="loginForm" class="auth-form" novalidate>
            <h2>Welcome Back</h2>
            <p class="auth-hint">Log in with your nickname or email.</p>
            ${field("text", "loginIdentifier", "Nickname or email", "Nickname or Email", { autocomplete: "username" })}
            ${field("password", "loginPassword", "Password", "Password", { autocomplete: "current-password" })}
            ${statusLine()}
            <button type="submit" class="auth-submit">Login</button>
        </form>
    `;
}

export function register() {
    return `
        ${tabs("register")}
        <!-- REGISTER FORM -->
        <form id="registerForm" class="auth-form" novalidate>
            <h2>Create Account</h2>
            ${field("text", "nickname", "Nickname", "Nickname (2-20 characters)", { minlength: "2", autocomplete: "username" })}
            ${field("text", "firstName", "First name", "First Name", { minlength: "2" })}
            ${field("text", "lastName", "Last name", "Last Name", { minlength: "2" })}
            <label class="sr-only" for="age">Age</label>
            <input type="number" class="auth-input" placeholder="Age" id="age" aria-label="Age" min="14" max="200">
            <label class="sr-only" for="gender">Gender</label>
            <select class="auth-input" id="gender" aria-label="Gender">
                <option value="">Gender</option>
                <option value="man">man</option>
                <option value="woman">woman</option>
            </select>
            ${field("email", "email", "Email", "Email", { autocomplete: "email" })}
            ${field("password", "password", "Password", "Password (8-20 characters)", { minlength: "8", autocomplete: "new-password" })}
            ${statusLine()}
            <button type="submit" class="auth-submit">Register</button>
        </form>
    `;
}
