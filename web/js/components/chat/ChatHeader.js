import UserAvatar from '../UserAvatar.js';


export default function ChatHeader(conversation = {}) {
    const { name = '', onlineStatus = '' } = conversation;
    const statusLabel = onlineStatus ? onlineStatus.charAt(0).toUpperCase() + onlineStatus.slice(1) : '';

    return `
        <div class="chat-header">
            <div class="chat-user-info">
                ${UserAvatar(conversation, { sizeClass: 'user-avatar--lg', className: 'chat-avatar' })}
                <div>
                    <h4>${name}</h4>
                    <span class="chat-status ${onlineStatus}">${statusLabel}</span>
                </div>
            </div>
        </div>
    `;
}
