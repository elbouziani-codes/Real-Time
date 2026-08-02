import Search from './Search.js';
import Categories from './Categories.js';
import Sort from './Sort.js';
import OnlineUsers from './OnlineUsers.js';
import OfflineUsers from './OfflineUsers.js';
import RecentMessages from './RecentMessages.js';

/**
 * Sidebar
 * Composes the full sidebar: title, search, categories, sort
 * and the shared users section (online / offline / recent messages).
 * Keeps the single `.users-section` wrapper exactly like the page.
 */
export default function Sidebar({
    searchPlaceholder = 'Search posts...',
    categories,
    sortOptions,
    onlineUsers = [],
    offlineUsers = [],
    recentMessages = [],
} = {}) {
    return `
        <aside class="sidebar">
            <h3 class="sidebar-title">🔍 Filter</h3>
            ${Search({ placeholder: searchPlaceholder })}
            ${Categories({ categories })}
            ${Sort({ options: sortOptions })}
            <div class="users-section">
                ${OnlineUsers({ users: onlineUsers })}
                ${OfflineUsers({ users: offlineUsers })}
                ${RecentMessages({ messages: recentMessages })}
            </div>
        </aside>
    `;
}
