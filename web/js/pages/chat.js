import { Navbar } from "./../components/navbar.js";
import ChatPage from "../components/chat/ChatPage.js";
import { me } from "../services/me.js";
import { getSortedConversations } from "../services/user.js";
import { currentChat } from "../services/messages.js";

export default function chat() {
    const conversations = getSortedConversations();

    return `
        ${Navbar(me)}
        ${ChatPage({ conversations, messages: [], activeId: currentChat?.UserB ?? null })}
    `;
}
