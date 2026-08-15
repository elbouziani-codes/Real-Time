import createUser from "../models/User.js";
export { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "./user.js";

export const DEFAULT_CONVERSATIONS = [
    createUser({ id: 1, name: "Sarah", avatarClass: "avatar--sarah", onlineStatus: "online", lastMessage: "See you tomorrow!", createdAt: "10m ago" }),
    createUser({ id: 2, name: "Ahmed", avatarClass: "avatar--ahmed", onlineStatus: "online", lastMessage: "Thanks for the help!", createdAt: "30m ago" }),
    createUser({ id: 3, name: "Fatima", avatarClass: "avatar--fatima", onlineStatus: "offline", lastMessage: "The design looks great!", createdAt: "1h ago" }),
    createUser({ id: 4, name: "Omar", avatarClass: "avatar--omar", onlineStatus: "online", lastMessage: "Let's work on that project", createdAt: "2h ago" }),
    createUser({ id: 5, name: "Layla", avatarClass: "avatar--layla", onlineStatus: "offline", lastMessage: "I will send you the files", createdAt: "4h ago" }),
    createUser({ id: 6, name: "Youssef", avatarClass: "avatar--youssef", onlineStatus: "away", lastMessage: "Great meeting you!", createdAt: "8h ago" }),
];

export const DEFAULT_MESSAGES = [
    { id: 1, mine: false, letter: "S", avatarClass: "avatar--sarah", content: "Hey Mohammed! How are you?", time: "2:30 PM" },
    { id: 2, mine: true, letter: "M", avatarClass: "avatar--mine", content: "Hi Sarah! I'm great, working on the forum project!", time: "2:32 PM" },
    { id: 3, mine: false, letter: "S", avatarClass: "avatar--sarah", content: "That sounds awesome! What tech stack are you using?", time: "2:33 PM" },
    { id: 4, mine: true, letter: "M", avatarClass: "avatar--mine", content: "Go for the backend with WebSocket, and vanilla JS frontend", time: "2:35 PM" },
    { id: 5, mine: false, letter: "S", avatarClass: "avatar--sarah", content: "Nice! I've been learning Go too. See you tomorrow! 😊", time: "2:50 PM" },
    { id: 6, mine: true, letter: "M", avatarClass: "avatar--mine", content: "See you! Let me know if you need any help with Go", time: "2:52 PM" },
];
