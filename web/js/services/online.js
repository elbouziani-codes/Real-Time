// Presence events from the WebSocket are parsed here and delegated to the
// users service. This module owns the backend protocol format (codes 3/4 and
// the UUID regex); the users service owns the online flag and DOM updates.

import { updateUserStatus, extractUserIds } from "./user.js";

// Applies one presence frame from the WebSocket.
//   code 3 → user came online
//   code 4 → user went offline
export function applyPresenceEvent(code, content) {
    const online = code == 3;
    for (const id of extractUserIds(content)) {
        updateUserStatus(id, online);
    }
}
