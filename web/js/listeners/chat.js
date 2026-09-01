import { me } from "../services/me.js";
import { selectChat, loadOlderMessages, sendMessage, sendTyping } from "../services/messages.js";
import { getChat, setChat } from "./users.js";
import throttle from "../utils/helpers.js";
import { seedAllUsers } from "../services/user.js";
import { renderConversationSidebar } from "./../services/messages.js";

// Wires the /chat page: selects the conversation chosen on the home page, then
// handles in-page conversation switching, sending, and older-message paging.
export function chatListener() {
    selectChat(getChat());

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
    chatInput?.addEventListener("input", throttle(sendTyping, 500));

    const scroller = document.querySelector(".messages-container");
    scroller?.addEventListener("scroll", throttle(onMessagesScroll, 200));
}



function onConversationClick(event) {
    const item = event.target.closest(".conversation-item");
    const friendId = item?.dataset.userId;
    if (!friendId) return;

    setChat({
        UserA: me.ID ?? "",
        UserB: friendId,
    });
    selectChat(getChat());
}

let windowScrollHandler = null;

export function onSideBareScroll() {
    const obj = document.querySelector(".conversations-list");
    if (!obj) return;

    const throttledLoadUsers = throttle(() => {
        seedAllUsers();
        renderConversationSidebar();
    }, 1000);

    // The window listener persists across SPA mounts: detach the previous one
    // so revisiting the home page never stacks duplicate scroll handlers.
    if (windowScrollHandler) {
        obj.removeEventListener("scroll", windowScrollHandler);
    }

    windowScrollHandler = () => {

        const scrollTop = obj.scrollTop;       
        const windowHeight = obj.clientHeight;
        const documentHeight = obj.scrollHeight; 
        if (scrollTop + windowHeight >= documentHeight -50) {
            throttledLoadUsers();
        }
    };
    obj.addEventListener("scroll", windowScrollHandler);
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
