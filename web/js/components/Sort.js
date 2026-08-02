import Filter from './Filter.js';
import FilterHeader from './FilterHeader.js';

const SORT_ICON = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="4" y1="6" x2="20" y2="6"></line>
        <line x1="6" y1="12" x2="18" y2="12"></line>
        <line x1="8" y1="18" x2="16" y2="18"></line>
    </svg>
`;

const DEFAULT_OPTIONS = ['Latest', 'Most liked', 'Most commented'];

/**
 * Sort
 * Filter group with a "Sort" header and a custom select dropdown.
 * Accepts custom options; falls back to the page defaults.
 */
export default function Sort({ options = DEFAULT_OPTIONS, title = 'Sort', icon = SORT_ICON } = {}) {
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
