/**
 * Category model — pure data structure (no rendering, no I/O, no seed data).
 * Seed/demo datasets live in js/services/seed.js.
 *
 * Fields:
 *  - id          unique identifier
 *  - name        category label (e.g. "Technology")
 *  - icon        icon markup shown next to the label
 *  - count       number of posts in this category
 *  - colorClass  CSS color modifier (e.g. category-item--tech)
 */
export default function createCategory({
    id = 0,
    name = '',
    icon = '',
    count = 0,
    colorClass = '',
} = {}) {
    return {
        id,
        name,
        icon,
        count,
        colorClass,
    };
}
