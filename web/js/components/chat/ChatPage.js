import ChatSidebar from './ChatSidebar.js';
import ChatMain from './ChatMain.js';


export default function ChatPage({
    conversations = [],
    messages = [],
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
