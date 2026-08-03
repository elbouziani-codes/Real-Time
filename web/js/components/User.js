import UserAvatar from './UserAvatar.js';

/**
 * User
 * One reusable component for every user row on the page.
 * All data comes from a User model; the optional second argument
 * only carries presentation options (variant / active state).
 *
 * Variants:
 *  - 'item'         -> .user-item        (online / offline / search result users)
 *  - 'message'      -> .last-message-item (recent messages)
 *  - 'conversation' -> .conversation-item (chat conversations)
 *
 * Model fields used: name, handle, letter, avatarClass, onlineStatus,
 * lastMessage, unreadCount, createdAt.
 */
const VARIANTS = {
    item: {
        container: 'user-item',
        infoClass: 'user-item-info',
        nameClass: 'user-item-name',
        subtitleClass: 'user-item-status',
        messageClass: '',
        timeClass: '',
        sizeClass: 'user-avatar--sm',
    },
    message: {
        container: 'last-message-item',
        infoClass: 'last-msg-info',
        nameClass: 'last-msg-name',
        subtitleClass: '',
        messageClass: 'last-msg-text',
        timeClass: 'last-msg-time',
        sizeClass: 'user-avatar--sm',
    },
    conversation: {
        container: 'conversation-item',
        infoClass: 'conv-info',
        nameClass: 'conv-name',
        subtitleClass: '',
        messageClass: 'conv-preview',
        timeClass: 'conv-meta',
        sizeClass: 'user-avatar--md',
    },
};

export default function User(user = {}, { variant = 'item', active = false } = {}) {
    const {
        name = '',
        handle = '',
        lastMessage = '',
        createdAt = '',
    } = user;

    const config = VARIANTS[variant] || VARIANTS.item;
    const containerClass = [config.container, active ? 'active' : ''].filter(Boolean).join(' ');

    const avatarHtml = UserAvatar(user, { sizeClass: config.sizeClass });

    let infoHtml = '';
    if (name) {
        infoHtml = `<div class="${config.infoClass}">`;
        infoHtml += `<span class="${config.nameClass}">${name}</span>`;
        if (config.subtitleClass && handle) infoHtml += `<span class="${config.subtitleClass}">${handle}</span>`;
        if (config.messageClass && lastMessage) infoHtml += `<span class="${config.messageClass}">${lastMessage}</span>`;
        infoHtml += '</div>';
    }

    const timeHtml = config.timeClass && createdAt ? `<span class="${config.timeClass}">${createdAt}</span>` : '';

    return `<div class="${containerClass}">${avatarHtml}${infoHtml}${timeHtml}</div>`;
}
