import { Navbar } from "./../components/navbar.js";
import ChatPage from "../components/chat/ChatPage.js";
import { me } from "../services/me.js";
import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "../services/user.js";
import { currentChat } from "../services/messages.js";

export default function chat() {
    const conversations = [...DEFAULT_RECENT_MESSAGES, ...DEFAULT_USERS];

    return `
        ${Navbar(me)}
        ${ChatPage({ conversations, messages: [], activeId: currentChat?.UserB ?? null })}
    `;
}
