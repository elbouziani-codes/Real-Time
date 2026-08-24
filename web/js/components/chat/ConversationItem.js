import UserAvatar from '../UserAvatar.js';
import { escapeHTML } from '../../utils/helpers.js';
import formatDateTime from '../../utils/time.js';

export default function ConversationItem(conversation = {}, { active = false } = {}) {
    const { name = '', lastMessage = '', lastMessageAt = 0 } = conversation;
    const userId = conversation.id ?? '';
    const timeLabel = Number(lastMessageAt) > 0
        ? formatDateTime(lastMessageAt)
        : '';

    return `
        <div class="conversation-item${active ? ' active' : ''}" data-user-id="${userId}">
            ${UserAvatar(conversation, { sizeClass: 'user-avatar--md', className: 'conv-avatar' })}
            <div class="conv-info">
                <div class="conv-name">${escapeHTML(name)}</div>
                <div class="conv-preview">${escapeHTML(lastMessage)}</div>
            </div>
            <span class="conv-meta">${escapeHTML(timeLabel)}</span>
        </div>
    `;
}
