import SidebarSection from './SidebarSection.js';
import UsersList from './UsersList.js';

/**
 * OfflineUsers
 * "🔴 Offline" heading + list of offline users.
 * Uses the same User component as every other user row.
 */
export default function OfflineUsers({ users = [] } = {}) {
    return SidebarSection({
        title: '🔴 Offline',
        body: UsersList({ variant: 'item', users }),
    });
}
