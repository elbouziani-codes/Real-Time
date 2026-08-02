const LOGOUT_SVG = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
        <polyline points="16 17 21 12 16 7"></polyline>
        <line x1="21" y1="12" x2="9" y2="12"></line>
    </svg>
`;

/**
 * LogoutButton
 * Logout action button shown next to the navbar user profile.
 */
export default function LogoutButton({ label = 'Logout' } = {}) {
    return `
        <button class="logout-btn" aria-label="Logout" title="Logout">
            ${LOGOUT_SVG}
            <span class="logout-label">${label}</span>
        </button>
    `;
}
