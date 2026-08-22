import { fetchMessages } from "../api/messages.js";
import createMessage from "../models/Message.js";
import ChatMessage from "../components/chat/ChatMessage.js";
import ConversationItem from "../components/chat/ConversationItem.js";
import ChatHeader from "../components/chat/ChatHeader.js";
import { me } from "./me.js";
import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES, updateRecentConversation } from "./user.js";
import { disconnectSocket, sendWsRequest } from "../websocket/socket.js";
import { applyPresenceEvent, extractUserIds } from "./online.js";
import { moveOnlineUserToFront } from "./user.js";
import { escapeHTML } from "../utils/helpers.js";
import { showToast } from "../utils/toast.js";

// The backend serves messages in pages (LIMIT offset+10 OFFSET offset).
const PAGE_SIZE = 10;

// Messages of the currently selected room, oldest -> newest.
let messages = [];
export let currentChat = { UserA: "", UserB: "" };

let loading = false;
let hasMore = true;
let offset = 0;
// Bumped on every selectChat so a stale response of a previous room (or a
// previous selection of the same room) can never overwrite the current one.
let sessionId = 0;
const knownIds = new Set();

function uuidString(value) {
    return String(value ?? "");
}

function normalizeTimestamp(value) {
    const numeric = Number(value);
    if (!Number.isFinite(numeric) || numeric <= 0) return 0;
    return numeric < 1e12 ? numeric * 1000 : numeric;
}

// Maps the backend MessageOutput shape ({id, sender,
// chat_id, content, created_at}) onto the frontend message model.
function normalizeMessage(raw = {}) {
    return createMessage({
        id: raw.id ?? 0,
        roomId: raw.chat_id ?? 0,
        senderId: raw.sender ?? 0,
        content: raw.content ?? raw.Content ?? "",
        createdAt: normalizeTimestamp(raw.created_at ?? raw.CreatedAt ?? 0),
    });
}

function sortMessages(list = []) {
    return [...list].sort((a, b) => {
        const timeA = Number(a.createdAt ?? 0);
        const timeB = Number(b.createdAt ?? 0);
        return timeA - timeB;
    });
}

function findFriend(userId) {
    if (!userId) return null;
    return [...DEFAULT_RECENT_MESSAGES, ...DEFAULT_USERS].find(
        (user) => String(user.id) === String(userId),
    ) ?? null;
}

function formatTime(createdAt) {
    if (!createdAt) return "";
    const date = new Date(Number(createdAt));
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
    list.innerHTML = sortMessages(messages).map((message) => ChatMessage(toView(message))).join("");
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

// The server resolves the shared room from the two participants, so history
// loads even after a refresh.
export function selectChat(next) {
    sessionId += 1;
    loading = false;
    currentChat = {
        UserA: next?.UserA ?? "",
        UserB: next?.UserB ?? "",
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
    renderLoading();
    getMessages();
}

// Loads one page of messages for the selected user. With { older: true }, it
// prepends the next page while keeping the scroll position.
async function getMessages(options = {}) {
    const { older = false } = options;
    if (!currentChat?.UserB) return;
    if (loading) return;
    if (older && !hasMore) return;

    loading = true;
    const sid = sessionId;
    const requestOffset = older ? offset : 0;

    try {
        const response = await fetchMessages({
            friend: currentChat.UserB,
            offset: requestOffset,
        });

        // A newer selection (or a different room) superseded this request.
        if (sid !== sessionId || !currentChat?.UserB) return;

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
            messages = sortMessages([...fresh, ...messages.filter((message) => !message.temp)]);
        } else {
            // Optimistic (unsent-confirmed) messages are superseded by history.
            messages = sortMessages([...messages.filter((message) => !message.temp), ...fresh]);
        }

        if (list.length < PAGE_SIZE) hasMore = false;
        offset = messages.length;

        if (older) {
            const scroller = messagesScrollerEl();
            const previousTop = scroller?.scrollTop ?? 0;
            const previousHeight = scroller?.scrollHeight ?? 0;
            renderAllMessages();
            if (scroller) scroller.scrollTop = previousTop + (scroller.scrollHeight - previousHeight);
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
    if (!currentChat?.UserB || loading || !hasMore) return;
    if (messages.length === 0) return;
    getMessages({ older: true });
}

// Appends one new message (WebSocket delivery or optimistic own send) without
// reloading the whole history. Skips duplicates and only scrolls when the user
// is already near the bottom, otherwise a "New messages" pill is offered.
function appendMessage(message) {
    if (knownIds.has(message.id)) return;
    knownIds.add(message.id);
    messages = sortMessages([...messages, message]);
    hideTypingIndicator();

    const list = messagesListEl();
    const scroller = messagesScrollerEl();
    if (!list || !scroller) return;

    const nearBottom = scroller.scrollHeight - scroller.scrollTop - scroller.clientHeight < 80;
    renderAllMessages();
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
        roomId: "",
        senderId: currentChat.UserA,
        content: content.trim(),
        createdAt: Date.now(),
    });
    Object.assign(temp, { temp: true });
    updateRecentConversation(friend, {
        lastMessage: content.trim(),
        createdAt: temp.createdAt,
    });
    renderConversationSidebar();
    appendMessage(temp);
    return true;
}

// Handles a WebSocket frame. code 200 is a new message for whoever receives
// it; code 2 is a typing frame; codes 3/4 are presence events; 404 is an error.
export async function handleWsMessage(data) {
    if (!data || typeof data !== "object") return;

    const code = data.code;
    if (code == 200 && data.chat_id) {
        const senderId = uuidString(data.sender);
        const roomId = uuidString(data.chat_id);
        if (senderId && roomId) {
            updateRecentConversation(senderId, {
                lastMessage: data.content ?? "",
                createdAt: normalizeTimestamp(data.created_at ?? data.CreatedAt ?? Date.now()),
            });
            renderConversationSidebar();
            // A message for the currently open conversation.
            if (currentChat?.UserB && String(currentChat.UserB) === senderId) {
                appendMessage(normalizeMessage(data));
            } else {
                notifyNewMessage(senderId, data.content ?? "");
            }
        }
        return;
    }

    if (code == 1) {
        ("WebSocket replaced by another connection:", data.content);
        disconnectSocket();
        const { navigate } = await import("../router/router.js");
        navigate("/end");
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
        // Online users surface at the top of the people list; fetch anyone
        // no page has brought in yet.
        for (const id of code == 3 ? extractUserIds(data.content) : []) {
            await moveOnlineUserToFront(id);
            renderConversationSidebar();
        }
        return;
    }

    if (code == 404) {
        console.error("WebSocket request failed:", data.content);
    }
}

// The recipient is connected but looking somewhere else: raise a clickable
// toast that jumps straight into the conversation.
function notifyNewMessage(senderId, content) {
    const friend = findFriend(senderId);
    showToast({
        title: friend?.name ?? "New message",
        message: content,
        onClick: () => openConversation(senderId),
    });
}

async function openConversation(senderId) {
    if (window.location.pathname === "/chat") {
        selectChat({ UserA: me.ID, UserB: senderId });
        return;
    }
    // Another page: stage the chat so the /chat listener selects it on mount.
    const [{ navigate }, { setChat }] = await Promise.all([
        import("../router/router.js"),
        import("../listeners/users.js"),
    ]);
    setChat({ UserA: me.ID, UserB: senderId });
    navigate("/chat");
}

function renderConversationSidebar() {
    const list = document.querySelector(".conversations-list");
    const currentPath = window.location.pathname;
    if (!list || currentPath !== "/chat") return;

    const conversations = [...DEFAULT_RECENT_MESSAGES, ...DEFAULT_USERS];
    const activeId = currentChat?.UserB ?? null;
    list.innerHTML = conversations
        .map((conversation) =>
            ConversationItem(conversation, {
                active: (conversation.id) === (activeId),
            }),
        )
        .join("");
}
