import { Categories, Sort } from './Filter.js';
import SidebarSection from './SidebarSection.js';
import UserList from './UserList.js';
import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES, DEFAULT_CATEGORIES, DEFAULT_FILTER } from '../services/seed.js';

/**
 * Online / Offline / Recent Messages sections.
 * Each is a titled SidebarSection wrapping a UserList —
 * the User component powers every user row.
 */
const OnlineSection = ({ users = [] }) =>
    SidebarSection({ title: '🟢 Online', body: UserList(users, { variant: 'item' }) });

const OfflineSection = ({ users = [] }) =>
    SidebarSection({ title: '🔴 Offline', body: UserList(users, { variant: 'item' }) });

const RecentMessagesSection = ({ users = [] }) =>
    SidebarSection({
        title: '💬 Recent Messages',
        titleClass: 'users-section-title',
        body: UserList(users, { variant: 'message' }),
    });

/**
 * Sidebar
 * Composes the full sidebar from models:
 *  - users          array of User models  -> split into Online / Offline by onlineStatus
 *  - categories     array of Category models
 *  - filters        Filter model (sorting options + search placeholder)
 *  - recentMessages array of User models (lastMessage + createdAt)
 *
 * Keeps the single `.users-section` wrapper exactly like the page.
 */
export default function Sidebar({
    users = DEFAULT_USERS,
    categories,
    filters = DEFAULT_FILTER,
    recentMessages = DEFAULT_RECENT_MESSAGES,
} = {}) {
    const onlineUsers = users.filter((user) => user.onlineStatus === 'online');
    const offlineUsers = users.filter((user) => user.onlineStatus !== 'online');
    // categories fall back to the Filter model's own categories
    const resolvedCategories = categories ?? filters.categories ?? DEFAULT_CATEGORIES;

    return `
        <aside class="sidebar">
            <h3 class="sidebar-title">🔍 Filter</h3>
            <input type="search" class="sidebar-search" placeholder="${filters.search}">
            ${Categories({ categories: resolvedCategories })}
            ${Sort({ options: filters.sorting })}
            <div class="users-section">
                ${OnlineSection({ users: onlineUsers })}
                ${OfflineSection({ users: offlineUsers })}
                ${RecentMessagesSection({ users: recentMessages })}
            </div>
        </aside>
    `;
}
