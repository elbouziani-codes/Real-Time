import SidebarSection from './SidebarSection.js';
import UsersList from './UsersList.js';

/**
 * RecentMessages
 * "💬 Recent Messages" heading + list of recent messages.
 * Reuses the User component — only the surrounding container differs.
 */
export default function RecentMessages({ messages = [] } = {}) {
    return SidebarSection({
        title: '💬 Recent Messages',
        titleClass: 'users-section-title',
        body: UsersList({ variant: 'message', users: messages }),
    });
}
