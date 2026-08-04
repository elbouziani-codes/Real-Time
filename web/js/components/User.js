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

// Class names per variant. An empty class means the variant omits that field.
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

/**
 * Renders a `<span>` only when both the variant provides a class and the
 * model provides a value — variants opt out of a field by leaving its class empty.
 */
function renderSpan(className, value) {
    if (!className || !value) return '';
    return `<span class="${className}">${value}</span>`;
}

/**
 * Info column: name plus the optional subtitle / message lines.
 * Rendered only for users that have a name.
 */
function renderInfo(classes, { name, handle, lastMessage }) {
    if (!name) return '';

    const lines = [
        `<span class="${classes.nameClass}">${name}</span>`,
        renderSpan(classes.subtitleClass, handle),
        renderSpan(classes.messageClass, lastMessage),
    ].join('');

    return `<div class="${classes.infoClass}">${lines}</div>`;
}

export default function User(user = {}, { variant = 'item', active = false } = {}) {
    const { name = '', handle = '', lastMessage = '', createdAt = '' } = user;
    const classes = VARIANTS[variant] || VARIANTS.item;

    const containerClass = [classes.container, active ? 'active' : ''].filter(Boolean).join(' ');
    const avatarHtml = UserAvatar(user, { sizeClass: classes.sizeClass });
    const infoHtml = renderInfo(classes, { name, handle, lastMessage });
    const timeHtml = renderSpan(classes.timeClass, createdAt);

    return `<div class="${containerClass}">${avatarHtml}${infoHtml}${timeHtml}</div>`;
}
