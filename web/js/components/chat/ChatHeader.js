import UserAvatar from '../UserAvatar.js';
import { escapeHTML } from '../../utils/helpers.js';

export default function ChatHeader(conversation = {}) {
    const { name = '', online = false, id = '' } = conversation;
    const userId = id ?? '';
    const statusLabel = online ? 'Online' : 'Offline';
    const statusClass = online ? 'online' : 'offline';

    return `
        <div class="chat-header" data-user-id="${userId}">
            <div class="chat-user-info">
                ${UserAvatar(conversation, { sizeClass: 'user-avatar--lg', className: 'chat-avatar' })}
                <div>
                    <h4>${escapeHTML(name)}</h4>
                    <span class="chat-status ${statusClass}">${escapeHTML(statusLabel)}</span>
                </div>
            </div>
        </div>
    `;
}
