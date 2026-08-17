
// Singleton WebSocket connection to the backend's /api/ws endpoint (the same
// origin the app is served from, so the session cookie travels with the
// handshake). One socket is shared by every feature and reconnects with a short
// delay when it drops.
let socket = null;
let reconnectTimer = null;
let shouldReconnect = true;
const subscribers = new Set();

function wsUrl() {
    const protocol = location.protocol === "https:" ? "wss:" : "ws:";
    return `${protocol}//${location.host}/api/ws`;
}

export function connectSocket() {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) return;

    shouldReconnect = true;
    socket = new WebSocket(wsUrl());

    socket.onopen = () => {
        console.log("WebSocket connected");
    };

    socket.onmessage = (event) => {
        let data;
        try {
            data = JSON.parse(event.data);
        } catch (error) {
            console.error("Invalid WebSocket message:", error);
            return;
        }
        subscribers.forEach((listener) => listener(data));
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    socket.onclose = () => {
        console.log("WebSocket closed");
        socket = null;
        if (shouldReconnect) {
            reconnectTimer = setTimeout(connectSocket, 3000);
        }
    };
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

// Subscribes a listener to every incoming frame. Returns an unsubscribe fn.
export function onWsMessage(listener) {
    subscribers.add(listener);
    return () => subscribers.delete(listener);
}
