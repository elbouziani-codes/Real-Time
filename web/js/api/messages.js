// Loads one page of messages for a chat room through the existing endpoint:
//  POST /api/getMessage   body: { id, me, freind, offset }
// id is the conversation uuid, me/freind the two participants, and offset the
// number of messages the client already holds (the backend serves the next 10
// older ones). The backend replies with { id, me: [MessageOutput...] } where
// the messages array rides under the "me" key.
export async function fetchMessages({ id, me, freind, offset = 0 } = {}) {
    try {
        const response = await fetch("/api/getMessage", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ id, me, freind, offset }),
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
