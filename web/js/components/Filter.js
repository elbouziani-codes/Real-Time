/**
 * Filter
 * All filter-related UI in one module:
 *  - Filter        (default) generic `.filter-group` wrapper
 *  - FilterHeader  icon + title header for a filter group
 *  - CategoryItem  single category checkbox row (Category model)
 *  - Categories    "Categories" filter group (array of Category models)
 *  - Sort          "Sort" filter group with a custom select (options array)
 *
 * All data comes from models — no category names, icons or sort labels
 * are hardcoded here.
 */

/**
 * Filter
 * Generic filter group wrapper (.filter-group).
 */
export default function Filter({ header = '', children = '' } = {}) {
    return `<div class="filter-group">${header}${children}</div>`;
}

/**
 * FilterHeader
 * Icon + title header used by filter groups.
 * Supports both existing header styles via class params:
 *  - Categories -> .category-header / .category-header__icon
 *  - Sort       -> .filter-group__header / .filter-group__icon
 */
export function FilterHeader({
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

const CHECK_SVG = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="20 6 9 17 4 12"></polyline>
    </svg>
`;

/**
 * CategoryItem
 * Single category checkbox row.
 * Receives a Category model: { name, colorClass, icon, count }.
 */
export function CategoryItem(category = {}) {
    const { name = '', colorClass = '', icon = '' } = category;

    return `
        <label class="category-item ${colorClass}">
            <input type="checkbox" class="category-item__checkbox">
            <span class="category-item__check">${CHECK_SVG}</span>
            <span class="category-item__icon">${icon}</span>
            <span class="category-item__text">${name}</span>
        </label>
    `;
}

const CATEGORIES_ICON = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
        <line x1="3" y1="9" x2="21" y2="9"></line>
        <line x1="9" y1="21" x2="9" y2="9"></line>
    </svg>
`;

/**
 * Categories
 * Filter group with a "Categories" header and one checkbox row per
 * Category model. Receives the categories from a model / caller.
 */
export function Categories({
    categories = [],
    title = 'Categories',
    icon = CATEGORIES_ICON,
} = {}) {
    const itemsHtml = categories.map((category) => CategoryItem(category)).join('');

    return Filter({
        header: FilterHeader({
            title,
            icon,
            headerClass: 'category-header',
            iconClass: 'category-header__icon',
        }),
        children: itemsHtml,
    });
}

const SORT_ICON = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="4" y1="6" x2="20" y2="6"></line>
        <line x1="6" y1="12" x2="18" y2="12"></line>
        <line x1="8" y1="18" x2="16" y2="18"></line>
    </svg>
`;

/**
 * Sort
 * Filter group with a "Sort" header and a custom select dropdown.
 * Receives the options array from a Filter model / caller.
 */
export function Sort({ options = [], title = 'Sort', icon = SORT_ICON } = {}) {
    const optionsHtml = options.map((option) => `<option>${option}</option>`).join('');

    return Filter({
        header: FilterHeader({ title, icon }),
        children: `
            <div class="filter-select-wrapper">
                <select class="filter-select">${optionsHtml}</select>
            </div>
        `,
    });
}
