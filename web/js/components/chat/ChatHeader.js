import UserAvatar from '../UserAvatar.js';
import { escapeHTML } from '../../utils/helpers.js';


export default function ChatHeader(conversation = {}) {
    const { name = '', onlineStatus = '', id = '' } = conversation;
    const userId = id?.Value ?? id ?? '';
    const statusLabel = onlineStatus ? onlineStatus.charAt(0).toUpperCase() + onlineStatus.slice(1) : '';

    return `
        <div class="chat-header" data-user-id="${userId}">
            <div class="chat-user-info">
                ${UserAvatar(conversation, { sizeClass: 'user-avatar--lg', className: 'chat-avatar' })}
                <div>
                    <h4>${escapeHTML(name)}</h4>
                    <span class="chat-status ${onlineStatus}">${escapeHTML(statusLabel)}</span>
                </div>
            </div>
        </div>
    `;
}
