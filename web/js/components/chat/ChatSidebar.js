import ConversationList from './ConversationList.js';
import UserList from '../UserList.js';


export default function ChatSidebar({ conversations = [], activeId = null, suggestions = [] } = {}) {
    return `
        <aside class="chat-sidebar">
            <div class="chat-sidebar-header">
                <h3>💬 Conversations</h3>
                <button class="btn-icon chat-sidebar__toggle" title="New chat" aria-expanded="false" aria-controls="chat-new-drawer" type="button">➕</button>
            </div>
            ${ConversationList({ conversations, activeId })}
            <div class="chat-sidebar__drawer" id="chat-new-drawer">
                <div class="chat-sidebar__drawer-head">
                    <h4>New chat</h4>
                    <span>People you have not messaged yet</span>
                </div>
                ${suggestions.length ? UserList(suggestions, { variant: "item", listClass: "chat-sidebar__suggestions" }) : '<div class="empty-state">No more users to start a new conversation with.</div>'}
            </div>
        </aside>
    `;
}
