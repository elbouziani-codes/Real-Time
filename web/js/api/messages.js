// The server finds the shared conversation from the authenticated user and
// selected user, so the browser never needs to cache a room id. Keyset paging:
// before_at/before_id name the oldest message already held (both omitted for
// the newest page), so live traffic can never shift a page under the cursor.
export async function fetchMessages({ friend, before_at = 0, before_id = "" } = {}) {
    try {
        const response = await fetch("/api/getMessage", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ friend, before_at, before_id }),
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
        console.error("Failed to fetch messages:", error);
        return { code: 500, body: "Error in request" };
    }
}
