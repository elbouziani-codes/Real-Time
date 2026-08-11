import { Categories } from "./Filter.js";
import SidebarSection from "./SidebarSection.js";
import UserList from "../UserList.js";
import { DEFAULT_CATEGORIES } from "../../services/categories.js";

import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "../../services/seed.js";


function renderRecentMessages(messages = []) {
  return SidebarSection({
    title: "💬 List Users",
    titleClass: "users-section-title",
    body: UserList(messages, { variant: "message", listClass: "users-list" }),
  });
}


function renderUsers(users = []) {
  return UserList(users, { variant: "item", listClass: "users-section" });
}


export default function Sidebar({
  users = DEFAULT_USERS,
  categories = DEFAULT_CATEGORIES,
  recentMessages = DEFAULT_RECENT_MESSAGES,
} = {}) {
  return `
        <aside class="sidebar">
            <h3 class="sidebar-title">Filter</h3>
            ${Categories({ categories })}
            <div class="users-section">
                ${renderRecentMessages(recentMessages)}
                ${renderUsers(users)}
            </div>
        </aside>
    `;
}
