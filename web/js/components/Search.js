/**
 * Search
 * Sidebar search input.
 */
export default function Search({ placeholder = 'Search posts...' } = {}) {
    return `<input type="search" class="sidebar-search" placeholder="${placeholder}">`;
}
