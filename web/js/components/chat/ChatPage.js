import ChatSidebar from './ChatSidebar.js';
import ChatMain from './ChatMain.js';


export default function ChatPage({
    conversations = [],
    messages = [],
    activeId = null,
} = {}) {
    const active = conversations.find((conversation) => String(conversation.id) === String(activeId)) || conversations[0] || {};

    return `
        <section class="chat-page">
            <div class="chat-container">
                ${ChatSidebar({ conversations, activeId: active.id ?? activeId })}
                ${ChatMain({ active, messages })}
            </div>
        </section>
    `;
}
