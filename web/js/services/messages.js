import { fetchMessages } from "../api/messages.js";
import createMessage from "../models/Message.js";
import ChatMessage from "../components/chat/ChatMessage.js";
import ConversationItem from "../components/chat/ConversationItem.js";
import ChatHeader from "../components/chat/ChatHeader.js";
import { me } from "./me.js";
import { getUser, getSortedConversations, updateLastMessage } from "./user.js";
import { disconnectSocket, sendWsRequest } from "../websocket/socket.js";
import { applyPresenceEvent } from "./online.js";
import { escapeHTML, debounce } from "../utils/helpers.js";
import { showToast } from "../utils/toast.js";
import { navigate } from "../router/router.js";
import { setChat } from "../listeners/users.js";
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

const hider = debounce(hideTypingIndicator, 1500);

function normalizeMessage(raw = {}) {
    return createMessage({
        id: raw.id ?? 0,
        roomId: raw.chat_id ?? 0,
        senderId: raw.sender ?? 0,
        content: raw.content ?? raw.Content ?? "",
        createdAt: raw.created_at ?? raw.CreatedAt ?? 0,
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

export function renderConversationSidebar() {
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


function hideTypingIndicator() {
    messagesListEl()?.querySelector(".typing-indicator")?.remove();
}


function showTypingIndicator() {
    const list = messagesListEl();
    hider();
    if (!list || list.querySelector(".typing-indicator")) return;
    const friend = getUser(currentChat?.UserB);
    const name = friend?.name ? escapeHTML(friend.name) : "Someone";
    list.insertAdjacentHTML(
        "beforeend",
        `<div class="typing-indicator" role="status" aria-label="${name} is typing"><span class="typing-dot"></span><span class="typing-dot"></span><span class="typing-dot"></span><span class="typing-name">${name} is typing…</span></div>`,
    );
    const scroller = messagesScrollerEl();
    if (scroller) scroller.scrollTop = scroller.scrollHeight;
    hider();

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

export function resetMessages() {
    sessionId += 1;
    messages = [];
    loading = false;
    hasMore = true;
    offset = 0;
    currentChat = { UserA: "", UserB: "" };
    knownIds.clear();
    hideTypingIndicator();
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

let lastTempId = null;

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
    lastTempId = temp.id;

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
        disconnectSocket();
        const { navigate } = await import("../router/router.js");
        navigate("/end");
        return;
    }

    // Typing indicator (code 2)
    if (code == 2) {
        if(currentChat?.UserB && String(currentChat.UserB) === data.sender) {
            showTypingIndicator();
        }
        return;
    }

    // Presence events: user online (3) / offline (4)
    if (code == 3 || code == 4) {
        await applyPresenceEvent(code, data.content);
        renderConversationSidebar();
        return;
    }

    if (code == 404) {
        if (lastTempId) {
            messages = messages.filter((m) => m.id !== lastTempId);
            knownIds.delete(lastTempId);
            lastTempId = null;
            renderAllMessages();
        }
        return;
    }


}

function handleIncomingMessage(data) {
    const senderId = data.sender;
    const roomId = data.chat_id;
    if (!senderId || !roomId) return;

    const partnerId = senderId === String(me?.ID) ? String(currentChat?.UserB) : senderId;
    if (!partnerId) return;

    updateLastMessage(partnerId, {
        lastMessage: data.content ?? "",
        createdAt: data.created_at ?? data.CreatedAt ?? Date.now(),
    });
    renderConversationSidebar();

    if (window.location.pathname === "/chat" && String(currentChat?.UserB) === partnerId) {
        appendMessage(normalizeMessage(data));
    } else {
        notifyNewMessage(partnerId, data.content ?? "");
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

function openConversation(senderId) {
    const chat = { UserA: me.ID, UserB: senderId };

    if (location.pathname === "/chat") {
        selectChat(chat);
        return;
    }

    setChat(chat);
    navigate("/chat");
}
