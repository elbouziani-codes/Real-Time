

export default async function fetchCategories() {
    try {
        const response = await fetch("/api/categories", {
            method: "GET",
        });

        if (!response.ok) {
            return { code: response.status, body: await response.text() };
        }
        return { code: response.status, body: await response.json() };
    } catch (error) {
        return { code: 500, body: "Error in request" };
    }
}
