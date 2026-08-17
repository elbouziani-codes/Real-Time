import { fetchMessages } from "../api/messages.js";
import createMessage from "../models/Message.js";
import ChatMessage from "../components/chat/ChatMessage.js";
import ChatHeader from "../components/chat/ChatHeader.js";
import { me } from "./me.js";
import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "./user.js";
import { sendWsRequest } from "../websocket/socket.js";
import { applyPresenceEvent } from "./online.js";
import { escapeHTML } from "../utils/helpers.js";

// The backend serves messages in pages (LIMIT offset+10 OFFSET offset).
const PAGE_SIZE = 10;

// friendId -> roomId learned from WebSocket chat_id fields (and restored from
// localStorage across reloads). The users list does not carry the conversation
// id, so this map is how the frontend knows which room to ask the HTTP API for.
export const roomIds = new Map();

const ROOM_IDS_KEY = "realtime-forum:room-ids:";
let roomIdsLoadedFor = null;

function currentUserKey() {
    return ROOM_IDS_KEY + (me?.ID?.Value ?? "anon");
}

function loadRoomIds() {
    const key = currentUserKey();
    if (roomIdsLoadedFor === key) return;
    roomIdsLoadedFor = key;
    roomIds.clear();
    try {
        const raw = localStorage.getItem(key);
        if (!raw) return;
        const parsed = JSON.parse(raw);
        if (parsed && typeof parsed == "object") {
            for (const [friendId, roomId] of Object.entries(parsed)) {
                if (friendId && roomId) roomIds.set(friendId, roomId);
            }
        }
    } catch (error) {
        console.warn("Failed to restore chat rooms:", error);
    }
}

function persistRoomId(friendId, roomId) {
    if (!friendId || !roomId) return;
    roomIds.set(friendId, roomId);
    try {
        const map = {};
        roomIds.forEach((value, key) => { map[key] = value; });
        localStorage.setItem(currentUserKey(), JSON.stringify(map));
    } catch (error) {
        console.warn("Failed to persist chat rooms:", error);
    }
}

// Looks up the room id for a friend, restoring the per-user cache on first use.
export function getRoomId(friendId) {
    if (!friendId) return null;
    loadRoomIds();
    return roomIds.get(friendId) ?? null;
}

// Messages of the currently selected room, oldest -> newest.
export let messages = [];
export let currentChat = { UserA: "", UserB: "", RoomID: null };

let loading = false;
let hasMore = true;
let offset = 0;
// Bumped on every selectChat so a stale response of a previous room (or a
// previous selection of the same room) can never overwrite the current one.
let sessionId = 0;
const knownIds = new Set();

function uuidString(value) {
    return String(value?.Value ?? value ?? "");
}

// Maps the backend MessageOutput shape ({id:{Value}, sender:{Value},
// chat_id:{Value}, content, created_at}) onto the frontend message model.
function normalizeMessage(raw = {}) {
    return createMessage({
        id: raw.id?.Value ?? raw.id ?? 0,
        roomId: raw.chat_id?.Value ?? raw.chat_id ?? 0,
        senderId: raw.sender?.Value ?? raw.sender ?? 0,
        content: raw.content ?? raw.Content ?? "",
        createdAt: raw.created_at ?? raw.CreatedAt ?? 0,
    });
}

function findFriend(userId) {
    if (!userId) return null;
    return [...DEFAULT_RECENT_MESSAGES, ...DEFAULT_USERS].find(
        (user) => String(user.id?.Value ?? user.id) === String(userId),
    ) ?? null;
}

function formatTime(createdAt) {
    if (!createdAt) return "";
    const date = new Date(Number(createdAt) * 1000);
    if (Number.isNaN(date.getTime())) return "";
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

// The ChatMessage component renders {mine, letter, avatarClass, name, content,
// time}. mine is decided by the authenticated user id (currentChat.UserA),
// never by a DOM class. All user content is escaped here, once.
function toView(message) {
    const mine = String(message.senderId) === String(currentChat?.UserA);
    const friend = findFriend(currentChat?.UserB);
    const letter = mine
        ? (me?.NickName?.[0] ?? "").toUpperCase()
        : (friend?.letter ?? "");
    const avatarClass = mine ? "avatar--mine" : (friend?.avatarClass ?? "");
    const senderName = mine ? (me?.NickName ?? "You") : (friend?.name ?? "");

    return {
        mine,
        letter,
        avatarClass,
        name: escapeHTML(senderName),
        content: escapeHTML(message.content ?? ""),
        time: formatTime(message.createdAt),
    };
}

function messagesListEl() {
    return document.querySelector(".messages-list");
}

function messagesScrollerEl() {
    return document.querySelector(".messages-container");
}

function scrollToBottom() {
    const scroller = messagesScrollerEl();
    if (scroller) scroller.scrollTop = scroller.scrollHeight;
}

function renderAllMessages() {
    const list = messagesListEl();
    if (!list) return;
    list.innerHTML = messages.map((message) => ChatMessage(toView(message))).join("");
}

function prependMessages(fresh) {
    const scroller = messagesScrollerEl();
    const list = messagesListEl();
    if (!scroller || !list) return;

    const previousHeight = scroller.scrollHeight;
    list.insertAdjacentHTML("afterbegin", fresh.map((message) => ChatMessage(toView(message))).join(""));
    // Keep the user's visual position instead of jumping to the bottom.
    scroller.scrollTop += scroller.scrollHeight - previousHeight;
}

function renderHeader() {
    const header = document.querySelector(".chat-header");
    if (!header) return;
    header.innerHTML = ChatHeader(findFriend(currentChat?.UserB) ?? {});
}

function renderActiveConversation() {
    document.querySelectorAll(".conversation-item").forEach((item) => {
        const friendId = item.dataset.userId;
        item.classList.toggle("active", friendId && String(friendId) === String(currentChat?.UserB));
    });
}

function renderLoading() {
    const list = messagesListEl();
    if (!list) return;
    list.innerHTML = `<div class="loading" role="status">Loading messages...</div>`;
}

function renderEmpty() {
    const list = messagesListEl();
    if (!list) return;
    list.innerHTML = `<div class="empty-state">No messages yet — start the conversation</div>`;
}

function renderPlaceholder() {
    const list = messagesListEl();
    if (!list) return;
    list.innerHTML = `<div class="chat-placeholder"><div class="placeholder-icon">💬</div><h3>Select a conversation to start chatting</h3></div>`;
}

function renderError(message) {
    const list = messagesListEl();
    if (!list) return;
    list.innerHTML = `<div class="empty-state">${escapeHTML(message)}</div>`;
}

function showNewMessagesButton() {
    const container = messagesScrollerEl();
    if (!container) return;
    if (container.querySelector(".chat-new-messages")) return;

    const button = document.createElement("button");
    button.className = "chat-new-messages";
    button.textContent = "↓ New messages";
    button.style.cssText =
        "position:sticky;bottom:12px;align-self:center;padding:8px 14px;border:none;" +
        "border-radius:20px;background:var(--primary);color:var(--text);cursor:pointer;" +
        "box-shadow:var(--shadow-md);flex-shrink:0;";
    button.addEventListener("click", () => {
        scrollToBottom();
        button.remove();
    });
    container.appendChild(button);
}

// --- Typing indicator -----------------------------------------------------

let typingHideTimer = null;

function hideTypingIndicator() {
    messagesListEl()?.querySelector(".typing-indicator")?.remove();
}

function showTypingIndicator() {
    const list = messagesListEl();
    if (!list || list.querySelector(".typing-indicator")) return;
    list.insertAdjacentHTML(
        "beforeend",
        `<div class="typing-indicator" role="status" aria-label="The other person is typing"><span class="typing-dot"></span><span class="typing-dot"></span><span class="typing-dot"></span></div>`,
    );
    const scroller = messagesScrollerEl();
    if (scroller) scroller.scrollTop = scroller.scrollHeight;
    clearTimeout(typingHideTimer);
    typingHideTimer = setTimeout(hideTypingIndicator, 4000);
}

// Sends the "typing" request over the existing WebSocket protocol so the other
// participant can show a typing indicator. Callers throttle it.
export function sendTyping() {
    const friend = currentChat?.UserB;
    if (!friend) return;
    sendWsRequest({
        request_type: "typing",
        mod: "",
        id: "",
        content: "",
        destination: friend,
    });
}

// --- Conversation selection & history -------------------------------------

// Selects a conversation: resets the previous room's state, then loads its
// history when the room id is known (otherwise the chat is still empty).
export function selectChat(next) {
    sessionId += 1;
    loading = false;
    // A previously selected conversation may have learned its room id from a
    // WebSocket frame or localStorage since the `chat` state was saved, so fall
    // back to the persisted room map when the passed value is still unknown.
    const roomId = next?.RoomID ?? next?.roomId ?? getRoomId(next?.UserB);
    currentChat = {
        UserA: next?.UserA ?? "",
        UserB: next?.UserB ?? "",
        RoomID: roomId,
    };
    messages = [];
    offset = 0;
    hasMore = true;
    knownIds.clear();
    hideTypingIndicator();

    renderHeader();
    renderActiveConversation();

    if (!currentChat.UserB) {
        renderPlaceholder();
        return;
    }
    if (currentChat.RoomID) {
        renderLoading();
        getMessages(currentChat.RoomID);
    } else {
        renderEmpty();
    }
}

// Loads one page of messages for the current room.
//   getMessages(roomId)              -> first page, render, scroll to bottom
//   getMessages(roomId, {older:true})-> next older page, prepend, keep position
export async function getMessages(roomId, options = {}) {
    const { older = false } = options;
    if (!currentChat || String(currentChat.RoomID) !== String(roomId)) return;
    if (loading) return;
    if (older && !hasMore) return;

    loading = true;
    const sid = sessionId;
    const requestOffset = older ? offset : 0;

    try {
        const response = await fetchMessages({
            id: roomId,
            me: currentChat.UserA,
            freind: currentChat.UserB,
            offset: requestOffset,
        });

        // A newer selection (or a different room) superseded this request.
        if (sid !== sessionId || !currentChat || String(currentChat.RoomID) !== String(roomId)) return;

        if (response.code != 200) {
            console.error("Failed to load messages:", response.code, response.body);
            renderError("Could not load messages. Please try again.");
            return;
        }

        const list = Array.isArray(response.body?.me) ? response.body.me : [];
        const fresh = [];
        for (const raw of list) {
            const normalized = normalizeMessage(raw);
            if (knownIds.has(normalized.id)) continue;
            knownIds.add(normalized.id);
            fresh.push(normalized);
        }

        if (older) {
            // The backend returns the page newest-first; flip it so prepending
            // keeps the whole list oldest -> newest. Optimistic (temp) copies
            // of the user's own sends are dropped: their authoritative versions
            // ride inside the fetched page (they are the newest rows), so
            // keeping both would duplicate them.
            fresh.reverse();
            messages = [...fresh, ...messages.filter((message) => !message.temp)];
        } else {
            // Optimistic (unsent-confirmed) messages are superseded by history.
            messages = [...messages.filter((message) => !message.temp), ...fresh];
        }

        if (list.length < PAGE_SIZE) hasMore = false;
        offset = messages.length;

        if (older) {
            prependMessages(fresh);
        } else if (messages.length === 0) {
            renderEmpty();
        } else {
            renderAllMessages();
            scrollToBottom();
        }
    } catch (error) {
        console.error("Failed to load messages:", error);
        renderError("Could not load messages. Please try again.");
    } finally {
        loading = false;
    }
}

// Infinite scroll entry point: fetch the next older page while preserving the
// scroll position. Guards against duplicate requests for the same page.
export function loadOlderMessages() {
    if (!currentChat?.RoomID || loading || !hasMore) return;
    if (messages.length === 0) return;
    getMessages(currentChat.RoomID, { older: true });
}

// Appends one new message (WebSocket delivery or optimistic own send) without
// reloading the whole history. Skips duplicates and only scrolls when the user
// is already near the bottom, otherwise a "New messages" pill is offered.
function appendMessage(message) {
    if (knownIds.has(message.id)) return;
    knownIds.add(message.id);
    messages.push(message);
    hideTypingIndicator();

    const list = messagesListEl();
    const scroller = messagesScrollerEl();
    if (!list || !scroller) return;

    const nearBottom = scroller.scrollHeight - scroller.scrollTop - scroller.clientHeight < 80;
    list.insertAdjacentHTML("beforeend", ChatMessage(toView(message)));

    if (nearBottom) {
        scrollToBottom();
    } else {
        showNewMessagesButton();
    }
}

// Sends a message over the existing WebSocket protocol. The backend never
// echoes a message back to its sender, so the sender's copy is appended
// optimistically with a temp id (dropped on the next history reload).
export function sendMessage(content) {
    const friend = currentChat?.UserB;
    if (!friend || !content?.trim()) return false;

    const ok = sendWsRequest({
        request_type: "message",
        mod: "Create",
        id: "",
        content: content.trim(),
        destination: friend,
    });

    if (!ok) {
        console.error("Message not sent: WebSocket is not connected");
        return false;
    }

    const temp = createMessage({
        id: "temp-" + Date.now(),
        roomId: currentChat.RoomID,
        senderId: currentChat.UserA,
        content: content.trim(),
        createdAt: Math.floor(Date.now() / 1000),
    });
    Object.assign(temp, { temp: true });
    appendMessage(temp);
    return true;
}

// Handles a WebSocket frame. code 200 is a new message for whoever receives
// it; code 2 is a typing frame; codes 3/4 are presence events; 404 is an error.
export function handleWsMessage(data) {
    if (!data || typeof data !== "object") return;

    const code = data.code;

    if (code == 200 && data.chat_id?.Value) {
        const senderId = uuidString(data.sender);
        const roomId = uuidString(data.chat_id);
        if (senderId && roomId) {
            persistRoomId(senderId, roomId);
            // A message for the currently open conversation — even when the
            // room id was still unknown (a fresh conversation), attach it now.
            if (currentChat?.UserB && String(currentChat.UserB) === senderId) {
                currentChat.RoomID = roomId;
                appendMessage(normalizeMessage(data));
            }
        }
        return;
    }

    if (code == 1) {
        console.log("WebSocket replaced by another connection:", data.content);
        return;
    }

    if (code == 2) {
        if (currentChat?.UserB && String(currentChat.UserB) === uuidString(data.sender)) {
            showTypingIndicator();
        }
        return;
    }

    if (code == 3 || code == 4) {
        applyPresenceEvent(code, data.content);
        return;
    }

    if (code == 404) {
        console.error("WebSocket request failed:", data.content);
    }
}
