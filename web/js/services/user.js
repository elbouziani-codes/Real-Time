import createUser from "../models/User.js";
import { fetchAllUsers, fetchUser, resetUsersPaging } from "../api/users.js";

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
function normalizeUser(raw = {}, existing = {}) {
    const nickName = raw.NickName ?? raw.name ?? existing.name ?? "";
    return createUser({
        id: raw.ID ?? raw.id ?? existing.id,
        name: nickName,
        handle: nickName ? "@" + nickName : existing.handle,
        letter: existing.letter,
        avatarClass: existing.avatarClass,
        online: existing.online ?? false,
        lastMessage: raw.LastMessage ?? raw.lastMessage ?? existing.lastMessage ?? "",
        lastMessageAt: Number(raw.LastMessageAt ?? raw.lastMessageAt ?? existing.lastMessageAt ?? 0),
    });
}

function upsertUser(raw = {}) {
    const id = String(raw.ID ?? raw.id ?? "");
    if (!id) return null;

    const user = normalizeUser(raw, users.get(id));
    users.set(id, user);
    sortConversationList();
    return user;
}

// ─── Public API ─────────────────────────────────────────────────────────────

export function getUser(userId) {
    return users.get(String(userId)) ?? null;
}

export function getSortedConversations() {
    return sortedConversations;
}

function getAllUsers() {
    return [...users.values()];
}

// Initial load: fetch all users from the backend (already ordered by
// lastMessageAt DESC, nick_name ASC). Populates the users Map and builds
// the sorted conversations list.
export async function seedAllUsers() {
    //resetUsersPaging();
    //users.clear();
    for (const raw of await fetchAllUsers()) {
        upsertUser(raw);
    }
    // Apply buffered presence statuses; keep entries whose user hasn't been
    // paged in yet so a later scroll-seeded page still gets their status.
    for (const [key, online] of pendingStatus) {
        const user = users.get(key);
        if (!user) continue;
        user.online = online;
        pendingStatus.delete(key);
    }
    sortConversationList();
    return { users: getAllUsers(), conversations: sortedConversations };
}

// Fills in the profile of a user that exists in the Map but has no name yet
// (created as a fallback from a message/presence before being paged in).
async function hydrateProfile(userId) {
    const raw = await fetchUser(userId);
    const user = users.get(String(userId));
    if (!raw || !user || user.name) return;
    Object.assign(user, normalizeUser(raw), {
        online: user.online,
        lastMessage: user.lastMessage,
        lastMessageAt: user.lastMessageAt,
    });
    sortConversationList();
    try {
        const { renderConversationSidebar } = await import("./messages.js");
        renderConversationSidebar();
    } catch {}
}

// Applies a buffered presence status to a user that has just been created
// as a fallback (before any page-in had their profile).
function applyPendingStatus(userId) {
    if (!pendingStatus.has(userId)) return;
    const user = users.get(userId);
    if (!user) return;
    user.online = pendingStatus.get(userId);
    pendingStatus.delete(userId);
    updatePresenceDom(userId, user.online);
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
        applyPendingStatus(key);
        hydrateProfile(key);
    }
    user.lastMessage = lastMessage;
    user.lastMessageAt = Number(createdAt);
    sortConversationList();
}

// ─── Presence (online / offline) ────────────────────────────────────────────
// Updates ONLY the online flag and the DOM indicator. NEVER sorts.

// Presence frames can arrive before seedAllUsers() has run (e.g. refresh on
// "/"), so statuses seen for unknown users are buffered here and applied
// during seeding.
const pendingStatus = new Map();

export function updateUserStatus(userId, online) {
    const key = String(userId);
    const user = users.get(key);
    if (!user) {
        pendingStatus.set(key, online);
        return;
    }
    user.online = online;
    updatePresenceDom(key, online);
}

export async function ensureUser(userId) {
    const key = String(userId);
    if (!key) return null;
    if (users.has(key)) return users.get(key);

    const raw = await fetchUser(key);
    return raw ? upsertUser(raw) : null;
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


export function resetUsers() {
    users.clear();
    sortedConversations = [];
    resetUsersPaging();
}
