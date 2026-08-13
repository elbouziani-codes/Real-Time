import UserAvatar from './UserAvatar.js';


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
};


function renderSpan(className, value) {
    if (!className || !value) return '';
    return `<span class="${className}">${value}</span>`;
}

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
