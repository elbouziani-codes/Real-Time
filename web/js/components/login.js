

export  function login(){
    return `
                <div class="tabs">
                    <button id="loginTab" class="tab-btn active">Login</button>
                    <button id="registerTab" class="tab-btn">Register</button>
                </div>
        <!-- LOGIN FORM -->
                <form id="loginForm" class="auth-form">
                    <h2>Welcome Back</h2>
                    <p class="auth-hint">Demo: use "mohammed" / any password</p>
                    <input type="text" class="auth-input" placeholder="Nickname or Email" id="loginIdentifier" required>
                    <input type="password" class="auth-input" placeholder="Password" id="loginPassword" required>
                    <button type="submit" class="auth-submit">Login</button>
                </form>
    `
}

export  function register(){
    return `
                <div class="tabs">
                    <button id="loginTab" class="tab-btn">Login</button>
                    <button id="registerTab" class="tab-btn active">Register</button>
                </div>
        <!-- REGISTER FORM -->
                <form id="registerForm" class=" auth-form">
                    <h2>Create Account</h2>
                    <input type="text" class="auth-input" placeholder="Nickname" id="nickname" required>
                    <input type="text" class="auth-input" placeholder="First Name" id="firstName" required>
                    <input type="text" class="auth-input" placeholder="Last Name" id="lastName" required>
                    <input type="number" class="auth-input" placeholder="Age" id="age">
                    <select class="auth-input" id="gender">
                        <option value="">Gender</option>
                        <option>man</option>
                        <option>woman</option>
                    </select>
                    <input type="email" class="auth-input" placeholder="Email" id="email" required>
                    <input type="password" class="auth-input" placeholder="Password (min 6 chars)" id="password" required>
                    <button type="submit" class="auth-submit">Register</button>
                </form>
    `
}