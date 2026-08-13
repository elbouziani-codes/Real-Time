import ChatSidebar from './ChatSidebar.js';
import ChatMain from './ChatMain.js';
import { DEFAULT_CONVERSATIONS, DEFAULT_MESSAGES } from '../../services/seed.js';


export default function ChatPage({
    conversations = DEFAULT_CONVERSATIONS,
    messages = DEFAULT_MESSAGES,
} = {}) {
    const active = conversations[0] || {};

    return `
        <section class="chat-page">
            <div class="chat-container">
                ${ChatSidebar({ conversations, activeId: active.id })}
                ${ChatMain({ active, messages })}
            </div>
        </section>
    `;
}
