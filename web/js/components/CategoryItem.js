const CHECK_SVG = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="20 6 9 17 4 12"></polyline>
    </svg>
`;

/**
 * CategoryItem
 * Single category checkbox row.
 * Data-driven: label, color modifier class and icon markup come from props.
 */
export default function CategoryItem({ label = '', colorClass = '', icon = '' } = {}) {
    return `
        <label class="category-item ${colorClass}">
            <input type="checkbox" class="category-item__checkbox">
            <span class="category-item__check">${CHECK_SVG}</span>
            <span class="category-item__icon">${icon}</span>
            <span class="category-item__text">${label}</span>
        </label>
    `;
}
