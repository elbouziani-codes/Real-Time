import { Categories } from './Filter.js';
import SidebarSection from './SidebarSection.js';
import UserList from './UserList.js';
import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES, DEFAULT_CATEGORIES, DEFAULT_FILTER } from '../services/seed.js';

/**
 * Titled block holding the last-message rows.
 * The rows use the 'message' variant inside the page's `.users-list` wrapper.
 */
function renderRecentMessages(messages = []) {
    return SidebarSection({
        title: '💬 List Users',
        titleClass: 'users-section-title',
        body: UserList(messages, { variant: 'message', listClass: 'users-list' }),
    });
}

/**
 * Untitled block of user rows (name + @handle + status dot),
 * wrapped in the page's nested `.users-section` container.
 */
function renderUsers(users = []) {
    return UserList(users, { variant: 'item', listClass: 'users-section' });
}

/**
 * Sidebar
 * Composes the full sidebar from models:
 *  - users          array of User models  -> `.user-item` rows
 *  - categories     array of Category models
 *  - filters        Filter model (categories fallback)
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
    // categories fall back to the Filter model's own categories
    const resolvedCategories = categories ?? filters.categories ?? DEFAULT_CATEGORIES;

    return `
        <aside class="sidebar">
            <h3 class="sidebar-title">Filter</h3>
            ${Categories({ categories: resolvedCategories })}
            <div class="users-section">
                ${renderRecentMessages(recentMessages)}
                ${renderUsers(users)}
            </div>
        </aside>
    `;
}
