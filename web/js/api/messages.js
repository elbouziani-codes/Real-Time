// The server finds the shared conversation from the authenticated user and
// selected user, so the browser never needs to cache a room id.
export async function fetchMessages({ friend, offset = 0 } = {}) {
    try {
        const response = await fetch("/api/getMessage", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ friend, offset }),
        });

        const text = await response.text();

        let body;
        try {
            body = JSON.parse(text);
        } catch {
            body = text;
        }

        return { code: response.status, body };
    } catch (error) {
        return { code: 500, body: "Error in request" };
    }
}
