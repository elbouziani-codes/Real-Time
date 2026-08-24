import { fetchMessages } from "../api/messages.js";
import createMessage from "../models/Message.js";
import ChatMessage from "../components/chat/ChatMessage.js";
import ConversationItem from "../components/chat/ConversationItem.js";
import ChatHeader from "../components/chat/ChatHeader.js";
import { me } from "./me.js";
import { getUser, getSortedConversations, updateLastMessage } from "./user.js";
import { disconnectSocket, sendWsRequest } from "../websocket/socket.js";
import { applyPresenceEvent } from "./online.js";
import { escapeHTML } from "../utils/helpers.js";
import { showToast } from "../utils/toast.js";

// ─── Constants ──────────────────────────────────────────────────────────────
const PAGE_SIZE = 10;

// ─── Message state (per conversation) ───────────────────────────────────────
let messages = [];
let loading = false;
let hasMore = true;
let offset = 0;
let sessionId = 0;
const knownIds = new Set();

// ─── Currently selected conversation ────────────────────────────────────────
export let currentChat = { UserA: "", UserB: "" };

// ─── Helpers ────────────────────────────────────────────────────────────────

function uuidString(value) {
    return String(value ?? "");
}

function normalizeTimestamp(value) {
    const numeric = Number(value);
    if (!Number.isFinite(numeric) || numeric <= 0) return 0;
    return numeric < 1e12 ? numeric * 1000 : numeric;
}

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

function formatTime(createdAt) {
    if (!createdAt) return "";
    const date = new Date(Number(createdAt));
    if (Number.isNaN(date.getTime())) return "";
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

// ─── DOM helpers ────────────────────────────────────────────────────────────

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

// ─── Message rendering ──────────────────────────────────────────────────────

function toView(message) {
    const mine = String(message.senderId) === String(currentChat?.UserA);
    const friend = getUser(currentChat?.UserB);
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

function renderAllMessages() {
    const list = messagesListEl();
    if (!list) return;
    list.innerHTML = sortMessages(messages)
        .map((message) => ChatMessage(toView(message)))
        .join("");
}

function renderHeader() {
    const header = document.querySelector(".chat-header");
    if (!header) return;
    header.innerHTML = ChatHeader(getUser(currentChat?.UserB) ?? {});
}

function renderActiveConversation() {
    document.querySelectorAll(".conversation-item").forEach((item) => {
        const friendId = item.dataset.userId;
        item.classList.toggle(
            "active",
            friendId && String(friendId) === String(currentChat?.UserB),
        );
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

// ─── Conversation sidebar ───────────────────────────────────────────────────

function renderConversationSidebar() {
    const list = document.querySelector(".conversations-list");
    if (!list || window.location.pathname !== "/chat") return;

    const conversations = getSortedConversations();
    const activeId = currentChat?.UserB ?? null;
    list.innerHTML = conversations
        .map((conversation) =>
            ConversationItem(conversation, {
                active: String(conversation.id) === String(activeId),
            }),
        )
        .join("");
}

// ─── New messages button ────────────────────────────────────────────────────

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

// ─── Typing indicator ───────────────────────────────────────────────────────

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

// ─── Conversation selection & history loading ────────────────────────────────

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
    loadMessages();
}

async function loadMessages(options = {}) {
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
            messages = sortMessages([...fresh, ...messages.filter((m) => !m.temp)]);
        } else {
            messages = sortMessages([...messages.filter((m) => !m.temp), ...fresh]);
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

export function loadOlderMessages() {
    if (!currentChat?.UserB || loading || !hasMore) return;
    if (messages.length === 0) return;
    loadMessages({ older: true });
}

// ─── Message append (WebSocket + optimistic) ────────────────────────────────

function appendMessage(message) {
    if (knownIds.has(message.id)) return;
    knownIds.add(message.id);
    messages = sortMessages([...messages, message]);
    hideTypingIndicator();

    const scroller = messagesScrollerEl();
    if (!scroller) return;

    const nearBottom =
        scroller.scrollHeight - scroller.scrollTop - scroller.clientHeight < 80;
    renderAllMessages();
    if (nearBottom) {
        scrollToBottom();
    } else {
        showNewMessagesButton();
    }
}

// ─── Send message ───────────────────────────────────────────────────────────

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

    updateLastMessage(friend, {
        lastMessage: content.trim(),
        createdAt: temp.createdAt,
    });
    renderConversationSidebar();
    appendMessage(temp);
    return true;
}

// ─── WebSocket event router ─────────────────────────────────────────────────

export async function handleWsMessage(data) {
    if (!data || typeof data !== "object") return;

    const code = data.code;

    // New message (code 200)
    if (code == 200 && data.chat_id) {
        handleIncomingMessage(data);
        return;
    }

    // Connection replaced by another tab (code 1)
    if (code == 1) {
        console.warn("WebSocket replaced by another connection:", data.content);
        disconnectSocket();
        const { navigate } = await import("../router/router.js");
        navigate("/end");
        return;
    }

    // Typing indicator (code 2)
    if (code == 2) {
        if (currentChat?.UserB && String(currentChat.UserB) === uuidString(data.sender)) {
            showTypingIndicator();
        }
        return;
    }

    // Presence events: user online (3) / offline (4)
    if (code == 3 || code == 4) {
        applyPresenceEvent(code, data.content);
        renderConversationSidebar();
        return;
    }

    // Error (code 404)
    if (code == 404) {
        console.error("WebSocket request failed:", data.content);
    }
}

function handleIncomingMessage(data) {
    const senderId = uuidString(data.sender);
    const roomId = uuidString(data.chat_id);
    if (!senderId || !roomId) return;

    // Update conversation ordering (moves sender to top)
    updateLastMessage(senderId, {
        lastMessage: data.content ?? "",
        createdAt: normalizeTimestamp(data.created_at ?? data.CreatedAt ?? Date.now()),
    });
    renderConversationSidebar();

    // If the message belongs to the currently open conversation, append it.
    // Otherwise, show a toast notification.
    if (currentChat?.UserB && String(currentChat.UserB) === senderId) {
        appendMessage(normalizeMessage(data));
    } else {
        notifyNewMessage(senderId, data.content ?? "");
    }
}

function notifyNewMessage(senderId, content) {
    const friend = getUser(senderId);
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
    const [{ navigate }, { setChat }] = await Promise.all([
        import("../router/router.js"),
        import("../listeners/users.js"),
    ]);
    setChat({ UserA: me.ID, UserB: senderId });
    navigate("/chat");
}
