import ChatMessage from './ChatMessage.js';


export default function ChatMessages({ messages = [] } = {}) {
    return `
        <div class="messages-container">
            <div class="messages-list">
                ${messages.map((message) => ChatMessage(message)).join('')}
            </div>
        </div>
    `;
}
