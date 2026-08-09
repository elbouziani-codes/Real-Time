
import createUser from '../models/User.js';


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

