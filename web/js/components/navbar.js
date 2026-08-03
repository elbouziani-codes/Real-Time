import NavItems from './NavItems.js';
import UserAvatar from './UserAvatar.js';
import { CURRENT_USER } from '../services/seed.js';

/**
 * Navbar
 * Top navigation bar composed from smaller components:
 * Logo + NavItems + current user profile (avatar, name, logout).
 * The user block is rendered from a User model (defaults to CURRENT_USER).
 */

const LOGOUT_SVG = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
        <polyline points="16 17 21 12 16 7"></polyline>
        <line x1="21" y1="12" x2="9" y2="12"></line>
    </svg>
`;

export default function Navbar({user = CURRENT_USER } = {}) {
    const { name = '' } = user;

    return `
        <nav class="navbar">
            <div class="navbar-inner">
                <div class="logo">
                    <span class="logo-icon">💬</span>
                    <span class="logo-text">RealTime Forum</span>
                </div>
                ${NavItems()}
                <div class="user">
                    ${UserAvatar(user)}
                    <span class="user-name">${name}</span>
                    <button class="logout-btn" aria-label="Logout" title="Logout">
                        ${LOGOUT_SVG}<span class="logout-label">Logout</span>
                    </button>
                </div>
            </div>
        </nav>
    `;
}
