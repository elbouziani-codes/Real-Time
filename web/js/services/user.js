import createUser from "../models/User.js";
import { fetchAllUsers, fetchUser, resetUsersPaging } from "../api/users.js";
import { isOnline } from "./online.js";
import { me } from "./me.js";

// The sidebar is two lists fed by the same pages:
//   DEFAULT_RECENT_MESSAGES - people who share messages with me, newest first
//   DEFAULT_USERS           - everyone else, alphabetical
export const DEFAULT_USERS = [];
export const DEFAULT_RECENT_MESSAGES = [];

function byName(a, b) {
    return String(a.name ?? "").localeCompare(String(b.name ?? ""), undefined, {
        sensitivity: "base",
        numeric: true,
    });
}

function normalizeUser(raw = {}) {
    return createUser({
        id: raw.ID,
        name: raw.NickName,
        handle: "@" + raw.NickName,
        onlineStatus: isOnline(raw.ID) ? "online" : "offline",
        lastMessage: raw.LastMessage ?? "",
        // LastMessageAt doubles as the conversation recency (0 = never talked).
        createdAt: Number(raw.LastMessageAt ?? 0),
    });
}

async function loadPage() {
    for (const raw of await fetchAllUsers()) {
        const user = normalizeUser(raw);
        (user.createdAt ? DEFAULT_RECENT_MESSAGES : DEFAULT_USERS).push(user);
    }
    sortUsers();
}

function sortUsers() {
    DEFAULT_RECENT_MESSAGES.sort((a, b) => b.createdAt - a.createdAt || byName(a, b));
    DEFAULT_USERS.sort(byName);
}

const findUser = (id) =>
    [...DEFAULT_RECENT_MESSAGES, ...DEFAULT_USERS].find((user) => String(user.id) === String(id));

function removeFrom(list, id) {
    const index = list.findIndex((user) => String(user.id) === String(id));
    return index >= 0 ? list.splice(index, 1)[0] : null;
}

export async function seedAllUsers() {
    resetUsersPaging();
    DEFAULT_USERS.length = 0;
    DEFAULT_RECENT_MESSAGES.length = 0;
    await loadPage();
    return { users: DEFAULT_USERS, recentMessages: DEFAULT_RECENT_MESSAGES };
}

// A message arrived from friendId: move them to the top of the conversations.
export function updateRecentConversation(friendId, { lastMessage = "", createdAt = Date.now() } = {}) {
    const key = String(friendId);
    if (!key) return;

    const user = removeFrom(DEFAULT_USERS, key)
        ?? removeFrom(DEFAULT_RECENT_MESSAGES, key)
        ?? createUser({ id: key });

    user.lastMessage = lastMessage;
    user.createdAt = Number(createdAt);
    DEFAULT_RECENT_MESSAGES.unshift(user); // newest first, so the front is enough
}

// A user came online: pin them to the front of the people list, fetching their
// profile when no page has brought them in yet.
export async function moveOnlineUserToFront(userId) {
    const key = String(userId);
    if (!key || key === String(me?.ID)) return;

    let user = findUser(key);
    if (!user) {
        const raw = await fetchUser(key);
        if (!raw) return;
        user = normalizeUser({ ...raw, LastMessageAt: 0 });
    }

    removeFrom(DEFAULT_USERS, key);
    removeFrom(DEFAULT_RECENT_MESSAGES, key);
    user.onlineStatus = "online";
    DEFAULT_USERS.unshift(user);
}
