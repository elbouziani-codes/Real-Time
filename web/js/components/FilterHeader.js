/**
 * FilterHeader
 * Icon + title header used by filter groups.
 * Supports both existing header styles via class params:
 *  - Categories -> .category-header / .category-header__icon
 *  - Sort       -> .filter-group__header / .filter-group__icon
 */
export default function FilterHeader({
    title = '',
    icon = '',
    headerClass = 'filter-group__header',
    iconClass = 'filter-group__icon',
} = {}) {
    return `
        <div class="${headerClass}">
            <span class="${iconClass}">${icon}</span>
            <span>${title}</span>
        </div>
    `;
}
