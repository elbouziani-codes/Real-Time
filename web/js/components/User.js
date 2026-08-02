import UserAvatar from './UserAvatar.js';

/**
 * User
 * One reusable component for every user row on the page.
 *
 * Variants:
 *  - 'item'         -> .user-item        (online / offline / search result users)
 *  - 'message'      -> .last-message-item (recent messages)
 *  - 'conversation' -> .conversation-item (chat conversations)
 *
 * Supported options: avatar (letter/image), status (online/offline/away),
 * badge, subtitle (@handle), last message text and timestamp.
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

export default function User({
    variant = 'item',
    active = false,
    name = '',
    subtitle = '',
    message = '',
    time = '',
    letter = '',
    image = '',
    colorClass = '',
    sizeClass = '',
    status = '',
    badge = '',
    className = '',
} = {}) {
    const config = VARIANTS[variant] || VARIANTS.item;
    const containerClass = [config.container, active ? 'active' : ''].filter(Boolean).join(' ');

    const avatarHtml = UserAvatar({
        letter,
        image,
        colorClass,
        sizeClass: sizeClass || config.sizeClass,
        className,
        status,
        badge,
    });

    let infoHtml = '';
    if (name) {
        infoHtml = `<div class="${config.infoClass}">`;
        infoHtml += `<span class="${config.nameClass}">${name}</span>`;
        if (subtitle) infoHtml += `<span class="${config.subtitleClass}">${subtitle}</span>`;
        if (message) infoHtml += `<span class="${config.messageClass}">${message}</span>`;
        infoHtml += `</div>`;
    }

    const timeHtml = time && config.timeClass ? `<span class="${config.timeClass}">${time}</span>` : '';

    return `<div class="${containerClass}">${avatarHtml}${infoHtml}${timeHtml}</div>`;
}
