import createUser from "../models/User.js";
import { fetchAllUsers, resetUsersPaging } from "../api/users.js";

// ─── Single source of truth ─────────────────────────────────────────────────
// Every user lives here. The `online` property controls ONLY the status dot.
// The `lastMessageAt` property controls ONLY conversation ordering.
const users = new Map();

// Derived, sorted list for the conversation sidebar. Rebuilt by
// getSortedConversations() and consumed by rendering code.
let sortedConversations = [];

// ─── Sorting (one place, never uses online) ─────────────────────────────────
function byName(a, b) {
    return String(a.name ?? "").localeCompare(String(b.name ?? ""), undefined, {
        sensitivity: "base",
        numeric: true,
    });
}

function sortConversationList() {
    const all = [...users.values()];
    all.sort((a, b) => {
        const aHas = a.lastMessageAt > 0;
        const bHas = b.lastMessageAt > 0;

        // Both have conversations → newest first
        if (aHas && bHas) {
            if (a.lastMessageAt !== b.lastMessageAt) return b.lastMessageAt - a.lastMessageAt;
            return byName(a, b);
        }

        // One has conversations → goes above
        if (aHas) return -1;
        if (bHas) return 1;

        // Neither has conversations → alphabetical
        return byName(a, b);
    });
    sortedConversations = all;
}

// ─── Normalization ──────────────────────────────────────────────────────────
function normalizeUser(raw = {}) {
    return createUser({
        id: raw.ID,
        name: raw.NickName,
        handle: "@" + raw.NickName,
        online: false,
        lastMessage: raw.LastMessage ?? "",
        lastMessageAt: Number(raw.LastMessageAt ?? 0),
    });
}

// ─── Public API ─────────────────────────────────────────────────────────────

export function getUser(userId) {
    return users.get(String(userId)) ?? null;
}

export function getSortedConversations() {
    return sortedConversations;
}

export function getAllUsers() {
    return [...users.values()];
}

// Initial load: fetch all users from the backend (already ordered by
// lastMessageAt DESC, nick_name ASC). Populates the users Map and builds
// the sorted conversations list.
export async function seedAllUsers() {
    //resetUsersPaging();
    //users.clear();
    for (const raw of await fetchAllUsers()) {
        const user = normalizeUser(raw);
        users.set(String(user.id), user);
    }
    sortConversationList();
    return { users: getAllUsers(), conversations: sortedConversations };
}

// A new message arrived from or to friendId: update their lastMessageAt and
// re-sort so the conversation moves to the top.
export function updateLastMessage(friendId, { lastMessage = "", createdAt = Date.now() } = {}) {
    const key = String(friendId);
    if (!key) return;

    let user = users.get(key);
    if (!user) {
        user = createUser({ id: key });
        users.set(key, user);
    }
    user.lastMessage = lastMessage;
    user.lastMessageAt = Number(createdAt);
    sortConversationList();
}

// ─── Presence (online / offline) ────────────────────────────────────────────
// Updates ONLY the online flag and the DOM indicator. NEVER sorts.

export function updateUserStatus(userId, online) {
    const user = users.get(String(userId));
    if (!user) return;
    user.online = online;
    updatePresenceDom(String(userId), online);
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

// ─── UUID extraction (backend protocol) ─────────────────────────────────────
const UUID_RE = /[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/g;

export function extractUserIds(content = "") {
    return String(content).match(UUID_RE) ?? [];
}


export function resetUsers(){
    resetUsersPaging()
}