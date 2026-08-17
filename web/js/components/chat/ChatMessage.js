export default function ChatMessage({ mine = false, letter = '', avatarClass = '', name = '', content = '', time = '' } = {}) {
    const side = mine ? 'mine' : 'other';
    const nameHtml = name ? `<div class="message-sender">${name}</div>` : '';

    return `
        <div class="message message-${side}">
            <div class="message-avatar user-avatar user-avatar--sm ${avatarClass}">${letter}</div>
            <div class="message-bubble">
                ${nameHtml}
                <div class="message-content">${content}</div>
                <div class="message-time">${time}</div>
            </div>
        </div>
    `;
}
