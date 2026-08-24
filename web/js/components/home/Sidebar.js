import { Categories, LikedFilter } from "./Filter.js";
import { DEFAULT_CATEGORIES } from "../../services/categories.js";

export default function Sidebar({ categories = DEFAULT_CATEGORIES } = {}) {
    return `
        <aside class="sidebar">
            <h3 class="sidebar-title">Filter</h3>
            ${Categories(categories)}
            ${LikedFilter()}
        </aside>
    `;
}
