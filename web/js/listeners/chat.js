import { me } from "../services/me.js";
import { getRoomId, selectChat, loadOlderMessages, sendMessage, sendTyping } from "../services/messages.js";
import { chat, setChat } from "./users.js";
import throttle from "../utils/helpers.js";

// Wires the /chat page: selects the conversation chosen on the home page, then
// handles in-page conversation switching, sending, and older-message paging.
export default function chatListener() {
    selectChat(chat);

    const conversationsList = document.querySelector(".conversations-list");
    conversationsList?.addEventListener("click", onConversationClick);

    const sendButton = document.querySelector(".btn-send");
    const chatInput = document.querySelector(".chat-input");
    sendButton?.addEventListener("click", onSend);
    chatInput?.addEventListener("keydown", (event) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            onSend();
        }
    });
    chatInput?.addEventListener("input", throttle(sendTyping, 2000));

    const scroller = document.querySelector(".messages-container");
    scroller?.addEventListener("scroll", throttle(onMessagesScroll, 200));
}

function onConversationClick(event) {
    const item = event.target.closest(".conversation-item");
    const friendId = item?.dataset.userId;
    if (!friendId) return;

    const roomId = getRoomId(friendId);
    setChat({
        UserA: me.ID?.Value ?? "",
        UserB: friendId,
        RoomID: roomId,
        RequestApi: roomId != null,
    });
    selectChat(chat);
}

function onSend() {
    const chatInput = document.querySelector(".chat-input");
    const content = chatInput?.value?.trim();
    if (!content) return;

    if (sendMessage(content)) {
        chatInput.value = "";
        chatInput.focus();
    }
}

function onMessagesScroll() {
    const scroller = document.querySelector(".messages-container");
    if (!scroller) return;
    if (scroller.scrollTop <= 30) {
        loadOlderMessages();
    }
}
