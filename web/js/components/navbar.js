
import NavItems from './NavItems.js';
import UserAvatar from './UserAvatar.js';
import LogoutButton from './LogoutButton.js';

/**
 * Navbar
 * Top navigation bar composed from smaller components:
 * Logo + NavItems + current user profile (avatar, name, logout).
 */


export default function Navbar({navItems, user = {}} = {}) {
    const {letter = 'M', colorClass = 'avatar--mine', status = 'online', name = ''} = user;

    return `
        <nav class="navbar">
            <div class="navbar-inner">
                <div class="logo">
                    <span class="logo-icon">💬</span>
                    <span class="logo-text">RealTime Forum</span>
                </div>
                ${NavItems({ items: navItems })}
                <div class="user">
                    ${UserAvatar({ letter, colorClass, status })}
                    <span class="user-name">${name}</span>
                    ${LogoutButton()}
                </div>
            </div>
        </nav>
    `;
}
