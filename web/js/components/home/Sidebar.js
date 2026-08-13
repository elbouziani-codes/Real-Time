import { Categories, LikedFilter } from "./Filter.js";
import UserList from "../UserList.js";
import { DEFAULT_CATEGORIES } from "../../services/categories.js";
import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "../../services/seed.js";

export default function Sidebar({
    users = DEFAULT_USERS,
    categories = DEFAULT_CATEGORIES,
    recentMessages = DEFAULT_RECENT_MESSAGES,
} = {}) {
    return `
        <aside class="sidebar">
            <h3 class="sidebar-title">Filter</h3>
            ${Categories({ categories })}
            ${LikedFilter()}
            <div class="users-section">
                <h4 class="users-section-title">💬 List Users</h4>
                ${UserList(recentMessages, { variant: "message", listClass: "users-list" })}
                ${UserList(users, { variant: "item", listClass: "users-section" })}
            </div>
        </aside>
    `;
}
