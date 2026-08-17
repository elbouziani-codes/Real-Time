// Tracks which users are currently connected based on the backend's WebSocket
// presence events:
//   code 3 -> "add client <uuid> [|| <uuid> ...]"
//   code 4 -> "closed webSocket client <uuid>"
// The HTTP user list carries no online status, so these events are the only
// source of truth for the online/offline indicators.

const UUID_RE = /[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/g;

const onlineUsers = new Set();

export function isOnline(userId) {
    if (!userId) return false;
    return onlineUsers.has(String(userId));
}

function extractUserIds(content = "") {
    return String(content).match(UUID_RE) ?? [];
}

// Applies one presence frame from the WebSocket to the tracked set and the UI.
export function applyPresenceEvent(code, content) {
    if (code == 3) {
        for (const id of extractUserIds(content)) setUserOnline(id, true);
    } else if (code == 4) {
        for (const id of extractUserIds(content)) setUserOnline(id, false);
    }
}

// Marks a single user online/offline and updates every visible indicator that
// belongs to them (avatar dots, chat header status label).
export function setUserOnline(userId, online) {
    if (!userId) return;
    const key = String(userId);
    if (online) {
        onlineUsers.add(key);
    } else {
        onlineUsers.delete(key);
    }
    updatePresenceDom(key, online);
}

function updatePresenceDom(userId, online) {
    const selector = `[data-user-id="${CSS.escape(userId)}"]`;
    document.querySelectorAll(selector + " .user-avatar__status").forEach((dot) => {
        dot.classList.toggle("online", online);
        dot.classList.toggle("offline", !online);
    });
    document.querySelectorAll(selector + " .chat-status").forEach((label) => {
        label.classList.toggle("online", online);
        label.classList.toggle("offline", !online);
        label.textContent = online ? "Online" : "Offline";
    });
}

// Re-applies the known presence to the current DOM (after a page re-render, so
// rows that were drawn before the events arrived get their correct state).
export function syncPresenceDom() {
    for (const id of onlineUsers) {
        updatePresenceDom(id, true);
    }
}
