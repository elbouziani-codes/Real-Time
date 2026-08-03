/**
 * Filter model — pure data structure (no rendering, no I/O, no seed data).
 * Seed/demo datasets live in js/services/seed.js.
 *
 * Fields:
 *  - categories  array of Category models
 *  - sorting     array of sort option labels
 *  - search      search placeholder / options
 */
export default function createFilter({
    categories = [],
    sorting = [],
    search = '',
} = {}) {
    return {
        categories,
        sorting,
        search,
    };
}
