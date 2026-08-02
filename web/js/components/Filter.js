/**
 * Filter
 * Generic filter group wrapper (.filter-group).
 * Used by Categories, Sort and any future filter block.
 */
export default function Filter({ header = '', children = '' } = {}) {
    return `<div class="filter-group">${header}${children}</div>`;
}
