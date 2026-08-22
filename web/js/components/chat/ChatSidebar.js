import ConversationList from './ConversationList.js';


export default function ChatSidebar({ conversations = [], activeId = null } = {}) {
    return `
        <aside class="chat-sidebar">
            <div class="chat-sidebar-header">
                <h3>💬 Conversations</h3>
            </div>
            ${ConversationList({ conversations, activeId })}
        </aside>
    `;
}
