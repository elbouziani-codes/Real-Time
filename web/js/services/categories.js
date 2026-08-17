import fetchCategories from './../api/categories.js';

export let DEFAULT_CATEGORIES = [];

export async function seedCategories() {
    const response = await fetchCategories();

    if (response.code == true && Array.isArray(response.body)) {
        DEFAULT_CATEGORIES = response.body;
        return;
    }

    DEFAULT_CATEGORIES = ["error in fetch Categories"];
}
