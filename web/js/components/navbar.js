import NavItems from './NavItems.js';
import UserAvatar from './UserAvatar.js';
import { navigate } from "./../router/router.js";
import { logoutMe } from "./../services/me.js";
import { escapeHTML } from "../utils/helpers.js";

const LOGOUT_SVG = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
        <polyline points="16 17 21 12 16 7"></polyline>
        <line x1="21" y1="12" x2="9" y2="12"></line>
    </svg>
`;

export function Navbar(user) {
    const { NickName = '' } = user;

    return `
        <nav class="navbar">
            <div class="navbar-inner">
                <div class="logo homePage-logo">
                    <span class="logo-icon">💬</span>
                    <span class="logo-text">RealTime Forum</span>
                </div>
                ${NavItems()}
                <div class="user">
                    ${UserAvatar(user)}
                    <span class="user-name">${escapeHTML(user.NickName)}</span>
                    <button class="logout-btn" aria-label="Logout" title="Logout">
                        ${LOGOUT_SVG}<span class="logout-label">Logout</span>
                    </button>
                </div>
            </div>
        </nav>
    `;
}

export function navBarListener() {
    const homePageLogo = document.querySelector(".homePage-logo");
    const homePage = document.querySelector(".homePage");
    const chatPage = document.querySelector(".chatPage");
    const logoutBtn = document.querySelector(".logout-btn");

    homePage?.addEventListener("click", () => navigate("/"));
    homePageLogo?.addEventListener("click", () => navigate("/"));
    chatPage?.addEventListener("click", () => navigate("/chat"));

    logoutBtn?.addEventListener("click", async () => {
        logoutBtn.disabled = true;
        const done = await logoutMe();
        if (!done) {
            logoutBtn.disabled = false;
            return;
        }
        navigate("/login");
    });
}
