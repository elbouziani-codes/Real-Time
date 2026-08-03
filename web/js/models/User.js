/**
 * User model — pure data structure (no rendering, no I/O, no seed data).
 * Seed/demo datasets live in js/services/seed.js.
 *
 * Fields:
 *  - id            unique identifier
 *  - name          display name
 *  - handle        @handle shown as subtitle
 *  - letter        avatar letter (derived from name when omitted)
 *  - avatarClass   CSS color modifier (e.g. avatar--sarah)
 *  - onlineStatus  'online' | 'offline' | 'away' | ''
 *  - lastMessage   preview text for recent-message/conversation rows
 *  - unreadCount   unread badge number
 *  - createdAt     timestamp / relative time label
 */
export default function createUser({
    id = 0,
    name = '',
    handle = '',
    letter = '',
    avatarClass = 'avatar--mine',
    onlineStatus = '',
    lastMessage = '',
    unreadCount = 0,
    createdAt = '',
} = {}) {
    return {
        id,
        name,
        handle,
        letter: letter || (name ? name.charAt(0).toUpperCase() : ''),
        avatarClass,
        onlineStatus,
        lastMessage,
        unreadCount,
        createdAt,
    };
}
