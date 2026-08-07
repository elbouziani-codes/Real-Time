/**
 * Seed service — the app's demo/mock data.
 * Models stay pure data structures; this service owns the actual datasets
 * (a stand-in for what the backend would return). Swap this module for a
 * real API service without touching components or models.
 */
import createUser from '../models/User.js';
import createPost from '../models/Post.js';
import createCategory from '../models/Category.js';
import createFilter from '../models/Filter.js';
import fetchPost from "./../api/posts.js"

/** Current logged-in user shown in the navbar. */
export const CURRENT_USER = createUser({
    id: 0,
    name: 'Mohammed',
    handle: '@mohammed',
    avatarClass: 'avatar--mine',
    onlineStatus: 'online',
});

/** Sidebar users (online + offline). */
export const DEFAULT_USERS = [
    createUser({ id: 1, name: 'Sarah Ahmed', handle: '@sarah', avatarClass: 'avatar--sarah', onlineStatus: 'online' }),
    createUser({ id: 2, name: 'Ahmed Benali', handle: '@ahmed', avatarClass: 'avatar--ahmed', onlineStatus: 'online' }),
    createUser({ id: 3, name: 'Omar El Idrissi', handle: '@omar', avatarClass: 'avatar--omar', onlineStatus: 'online' }),
    createUser({ id: 4, name: 'Fatima Zahra', handle: '@fatima', avatarClass: 'avatar--fatima', onlineStatus: 'offline' }),
    createUser({ id: 5, name: 'Layla Hassan', handle: '@layla', avatarClass: 'avatar--layla', onlineStatus: 'offline' }),
];

/** Sidebar "Recent Messages" rows (users with a last message + time). */
export const DEFAULT_RECENT_MESSAGES = [
    createUser({ id: 1, name: 'Sarah', avatarClass: 'avatar--sarah', lastMessage: 'See you tomorrow!', createdAt: '10m ago' }),
    createUser({ id: 2, name: 'Ahmed', avatarClass: 'avatar--ahmed', lastMessage: 'Thanks for the help!', createdAt: '30m ago' }),
    createUser({ id: 3, name: 'Fatima', avatarClass: 'avatar--fatima', lastMessage: 'The design looks great!', createdAt: '1h ago' }),
    createUser({ id: 4, name: 'Omar', avatarClass: 'avatar--omar', lastMessage: "Let's work on that project", createdAt: '2h ago' }),
];

/** Feed posts shown by default on the Home page. */

/** Category filter options shown in the sidebar. */
export const DEFAULT_CATEGORIES = [
    createCategory({
        id: 1,
        name: 'Technology',
        count: 12,
        colorClass: 'category-item--tech',
        icon: `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
                <line x1="8" y1="21" x2="16" y2="21"></line>
                <line x1="12" y1="17" x2="12" y2="21"></line>
            </svg>`,
    }),
    createCategory({
        id: 2,
        name: 'Programming',
        count: 8,
        colorClass: 'category-item--code',
        icon: `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 18 22 12 16 6"></polyline>
                <polyline points="8 6 2 12 8 18"></polyline>
            </svg>`,
    }),
    createCategory({
        id: 3,
        name: 'Gaming',
        count: 5,
        colorClass: 'category-item--gaming',
        icon: `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="6" y1="11" x2="10" y2="11"></line>
                <line x1="8" y1="9" x2="8" y2="13"></line>
                <line x1="15" y1="12" x2="15.01" y2="12"></line>
                <line x1="18" y1="10" x2="18.01" y2="10"></line>
                <path d="M17.32 5H6.68a4 4 0 0 0-3.978 3.59c-.006.052-.01.101-.017.152C2.604 9.416 2 14.456 2 16a3 3 0 0 0 3 3c1 0 1.5-.5 2-1l1.414-1.414A2 2 0 0 1 9.828 16h4.344a2 2 0 0 1 1.414.586L17 18c.5.5 1 1 2 1a3 3 0 0 0 3-3c0-1.545-.604-6.584-.685-7.258-.007-.05-.011-.1-.017-.151A4 4 0 0 0 17.32 5z"></path>
            </svg>`,
    }),
    createCategory({
        id: 4,
        name: 'News',
        count: 3,
        colorClass: 'category-item--news',
        icon: `
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v16a2 2 0 0 1-2 2Zm0 0a2 2 0 0 1-2-2v-9h2"></path>
                <path d="M18 14h-8"></path>
                <path d="M15 18h-5"></path>
                <line x1="10" y1="6" x2="8" y2="6"></line>
                <line x1="14" y1="6" x2="12" y2="6"></line>
                <line x1="18" y1="6" x2="16" y2="6"></line>
                <line x1="10" y1="10" x2="8" y2="10"></line>
                <line x1="14" y1="10" x2="12" y2="10"></line>
                <line x1="18" y1="10" x2="16" y2="10"></line>
            </svg>`,
    }),
];


/** Sort options shown in the sidebar select. */
export const DEFAULT_SORT_OPTIONS = ['Latest', 'Most liked', 'Most commented'];

/** Default filter configuration used by the sidebar. */
export const DEFAULT_FILTER = createFilter({
    categories: DEFAULT_CATEGORIES,
    sorting: DEFAULT_SORT_OPTIONS,
    search: 'Search posts...',
});
