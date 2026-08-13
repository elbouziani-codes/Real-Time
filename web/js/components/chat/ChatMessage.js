
export default function ChatMessage({ mine = false, letter = '', avatarClass = '', content = '', time = '' } = {}) {
    const side = mine ? 'mine' : 'other';

    return `
        <div class="message message-${side}">
            <div class="message-avatar user-avatar user-avatar--sm ${avatarClass}">${letter}</div>
            <div class="message-bubble">
                <div class="message-content">${content}</div>
                <div class="message-time">${time}</div>
            </div>
        </div>
    `;
}
