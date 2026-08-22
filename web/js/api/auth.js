export async function fetchMe() {
    try {
        const response = await fetch("/api/me", {
            method: "GET",
        });
        if (response.status == 200) {
            return { code: response.status, body: await response.json() };
        }
        return { code: response.status, body: "errors" };
    } catch (error) {
        console.error("Failed to get profile:", error);
        return { code: 0, body: "Network error" };
    }
}

export async function fetchLogin(credentials) {
    try {
        const response = await fetch("/api/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(credentials),
        });

        const text = await response.text();

        let data;

        try {
            data = JSON.parse(text);
        } catch {
            data = text;
        }
        return { ok: response.ok, body: data, code: response.status };
    } catch (error) {
        console.error("login request failed:", error);
        return { ok: false, body: "Network error", code: 0 };
    }
}

export async function fetchRegister(credentials) {
    try {
        const response = await fetch("/api/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
            body: JSON.stringify(credentials),
        });

        const text = await response.text();

        let data;
        try {
            data = JSON.parse(text);
        } catch {
            data = text;
        }

        return { ok: response.ok, body: data, code: response.status };
    } catch (error) {
        console.error("Register request failed:", error);
        return { ok: false, body: "Network error", code: 0 };
    }
}

export async function fetchLogout() {
    try {
        const response = await fetch("/api/logout", {
            method: "POST",
            credentials: "include",
        });
        return { code: response.status, body: await response.text() };
    } catch (error) {
        console.error("Logout request failed:", error);
        return { code: 500, body: "Error in request" };
    }
}
