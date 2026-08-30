// Holds the conversation being opened: the authenticated user and the other
// participant. The server resolves their shared room.
let chat = { UserA: "", UserB: "" };

export function setChat(next) {
    chat = { ...next };
}

export function getChat() {
    return chat;
}

export function resetChat() {
    chat = { UserA: "", UserB: "" };
}
