import {handleWsMessage} from "./../services/messages.js"

let socket = null;

function wsUrl() {
    const protocol = location.protocol === "https:" ? "wss:" : "ws:";
    return `${protocol}//${location.host}/api/ws`;
}

export function connectSocket() {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) return;

    socket = new WebSocket(wsUrl());

    socket.onopen = () => {};

    socket.onmessage = (event) => {
        let data;
        try {
            data = JSON.parse(event.data);
        } catch (error) {
            console.error("Invalid WebSocket message:", error);
            return;
        }
        handleWsMessage(data)

    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    socket.onclose = () => {
        socket = null;
    };
}

// Closes the socket (used on logout) and stops the automatic reconnect loop so
// a signed-out tab does not keep a stale connection alive.
export function disconnectSocket() {
    if (!socket) return;
    socket.close();
    socket = null;
}

// Sends one request in the existing protocol:
//  { request_type, mod, id, content, destination }
// Returns false when the socket is not open so callers can react.
export function sendWsRequest(request) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        console.warn("WebSocket is not connected; request not sent");
        return false;
    }
    socket.send(JSON.stringify(request));
    return true;
}

