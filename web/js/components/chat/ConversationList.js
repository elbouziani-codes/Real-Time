import ConversationItem from './ConversationItem.js';


export default function ConversationList({ conversations = [], activeId = null } = {}) {
    return `
        <div class="conversations-list">
            ${conversations
                .map((conversation) => ConversationItem(conversation, { active: conversation.id === activeId }))
                .join('')}
        </div>
    `;
}
