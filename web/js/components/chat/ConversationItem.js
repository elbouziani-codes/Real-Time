import UserAvatar from '../UserAvatar.js';


export default function ConversationItem(conversation = {}, { active = false } = {}) {
    const { name = '', lastMessage = '', createdAt = '' } = conversation;
    const userId = conversation.id?.Value ?? conversation.id ?? '';

    return `
        <div class="conversation-item${active ? ' active' : ''}" data-user-id="${userId}">
            ${UserAvatar(conversation, { sizeClass: 'user-avatar--md', className: 'conv-avatar' })}
            <div class="conv-info">
                <div class="conv-name">${name}</div>
                <div class="conv-preview">${lastMessage}</div>
            </div>
            <span class="conv-meta">${createdAt}</span>
        </div>
    `;
}
