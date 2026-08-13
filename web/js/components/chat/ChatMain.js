import ChatHeader from './ChatHeader.js';
import ChatMessages from './ChatMessages.js';
import ChatInput from './ChatInput.js';


export default function ChatMain({ active = {}, messages = [] } = {}) {
    return `
        <main class="chat-main">
            <div class="active-chat">
                ${ChatHeader(active)}
                ${ChatMessages({ messages })}
                ${ChatInput()}
            </div>
        </main>
    `;
}
