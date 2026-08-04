

async function fetchCategories() {
    const response = await fetch("/api/categories", {
        method: "GET",
    });

    const categories = await response.json();
    return categories
    
}