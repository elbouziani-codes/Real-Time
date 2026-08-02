import SidebarSection from './SidebarSection.js';
import UsersList from './UsersList.js';

/**
 * OnlineUsers
 * "🟢 Online" heading + list of online users.
 * Uses the same User component as every other user row.
 */
export default function OnlineUsers({ users = [] } = {}) {
    return SidebarSection({
        title: '🟢 Online',
        body: UsersList({ variant: 'item', users }),
    });
}
