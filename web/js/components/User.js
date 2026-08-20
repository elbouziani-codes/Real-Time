import UserAvatar from './UserAvatar.js';
import { escapeHTML } from '../utils/helpers.js';
import formatDateTime from "../utils/time.js";

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
    return `<span class="${className}">${escapeHTML(value)}</span>`;
}

function renderTime(className, value) {
    if (!className || !value) return '';
    const numericValue = Number(value);
    const rendered = Number.isFinite(numericValue) && numericValue > 0
        ? formatDateTime(numericValue)
        : String(value);
    return `<span class="${className}">${escapeHTML(rendered)}</span>`;
}

function renderInfo(classes, { name, handle, lastMessage }) {
    if (!name) return '';

    const lines = [
        `<span class="${classes.nameClass}">${escapeHTML(name)}</span>`,
        renderSpan(classes.subtitleClass, handle),
        renderSpan(classes.messageClass, lastMessage),
    ].join('');

    return `<div class="${classes.infoClass}">${lines}</div>`;
}

export default function User(user = {}, { variant = 'item', active = false } = {}) {
    const { id = '', name = '', handle = '', lastMessage = '', createdAt = '' } = user;
    const userId = id?.Value ?? id ?? '';
    const classes = VARIANTS[variant] || VARIANTS.item;

    const containerClass = [classes.container, active ? 'active' : ''].filter(Boolean).join(' ');
    const avatarHtml = UserAvatar(user, { sizeClass: classes.sizeClass });
    const infoHtml = renderInfo(classes, { name, handle, lastMessage });
    const timeHtml = classes.timeClass ? renderTime(classes.timeClass, createdAt) : '';

    return `<div class="${containerClass}" data-user-id="${userId}">${avatarHtml}${infoHtml}${timeHtml}</div>`;
}
