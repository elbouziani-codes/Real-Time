const CHECK_SVG = `
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="20 6 9 17 4 12"></polyline>
    </svg>
`;

export function CategoryItem(category = {}) {
    const { title = "", currentColor = "", icon = "", id = "" } = category;
    return `
        <label class="category-item ${currentColor}" >
            <input type="checkbox" class="category-item__checkbox" id="${id}">
            <span class="category-item__check">${CHECK_SVG}</span>
            <span class="category-item__icon">${icon}</span>
            <span class="category-item__text">${title}</span>
        </label>
    `;
}

export function Categories(categories = []) {
    return `<div class="filter-group">
          <div class="category-header">
            <span class="category-header__icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="3" y1="9" x2="21" y2="9"></line>
                <line x1="9" y1="21" x2="9" y2="9"></line>
              </svg>
            </span>
            <span>Categories</span>
          </div><div class="category-body">${categories.map(CategoryItem).join("")}</div>`;
}

export function LikedFilter() {
    return `<div class="filter-group">
        <div class="category-header">
            <span class="category-header__icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
                </svg>
            </span>
            <span>Liked</span>
        </div>
        <div class="category-body">
            <label class="category-item">
                <input type="checkbox" class="category-item__checkbox" id="likedFilter">
                <span class="category-item__check">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                </span>
                <span class="category-item__text">Liked posts</span>
            </label>
        </div>
    </div>`;
}
