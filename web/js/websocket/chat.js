import { onWsMessage } from "./socket.js";
import { handleWsMessage } from "../services/messages.js";

let subscribed = false;

// Routes incoming WebSocket frames to the chat messages service. Called once so
// a single subscriber lives for the whole session.
export function initChatSocket() {
    if (subscribed) return;
    subscribed = true;
    onWsMessage(handleWsMessage);
}
