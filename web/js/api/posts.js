import { navigate } from "./../router/router.js";

let loading = false;
let hasMore = true;

export async function fetchPost(cursor = null, filters = {}) {
    if (loading || !hasMore) return [];

    loading = true;

    try {

        const params = new URLSearchParams();
        if (cursor) params.append("cursor", cursor);
        for (const categoryID of(filters.categories || [])) {
            params.append("category", categoryID);
        }
        if (filters.liked) params.append("liked", "true");

        const query = params.toString();
        const response = await fetch("/api/posts" + (query ? "?" + query : ""), {
            method: "GET",
        });


        if (!response.ok) {
            if (response.status == 401) {
                navigate("/login");
                return [];
            }
            // I could add a towst warning that geteting psost failed
            return [];
        }

        let posts = await response.json();

        return posts;
    } catch (error) {
        return [];
    } finally {
        loading = false;
    }
}

export function resetFetchPost() {
    
    loading = false;
    hasMore = true;
}

export function stopFetchPost() {
    hasMore = false;
}

export async function fetchCreatePost(data) {
    try {
        const response = await fetch("/api/posts", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });

        const text = await response.text();
        let body;
        try {
            body = JSON.parse(text);
        } catch {
            body = text;
        }
        return { code: response.status, body };
    } catch (error) {
        console.error(error);
        return { code: 500, body: "Error in request" };
    }
}

export async function fetchPostDetails(postID) {
    try {
        const response = await fetch("/api/posts/" + postID, {
            method: "GET",
        });
        if (response.status != 200) {
            return { code: response.status, body: await response.text() };
        }
        return { code: response.status, body: await response.json() };
    } catch (error) {
        console.error(error);
        return { code: 500, body: "Error in request" };
    }
}
