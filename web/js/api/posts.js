// The backend fixes the feed page size at 20 and gives back a plain array with
// no pagination metadata, so a page shorter than that is the only signal that
// the feed is exhausted. Anything at or below it stops the infinite scroll.
const POST_PAGE_SIZE = 20;

let loading = false;
let hasMore = true;

// fetchPost requests one page of the feed using the backend's keyset cursor:
//  GET /api/posts?cursor=<last received post id>&category=<id>...&liked=true
// cursor is the id of the last post the client already holds, or null/empty for
// the first page. Filters ride along on the query string.
export async function fetchPost(cursor = null, filters = {}) {
    if (loading || !hasMore) return {code: 1000, body: [], len: 0};

    loading = true;

    try{
        const params = new URLSearchParams();
        if (cursor) params.append("cursor", cursor);
        for (const categoryID of (filters.categories || [])) {
            params.append("category", categoryID);
        }
        if (filters.liked) params.append("liked", "true");

        const query = params.toString();
        const response = await fetch("/api/posts" + (query ? "?" + query : ""), {
            method: "GET",
        });
        let allResult = await response.json()

        if (!Array.isArray(allResult)) {
            hasMore = false;
            // A 200 with a non-array body is the backend saying "no posts match"
            // (Go encodes a nil slice as null): a legitimately empty page, so
            // report it as a success so the feed can show an empty state.
            return {code: response.status, body: [], len: 0};
        }
        if (allResult.length < POST_PAGE_SIZE) {
            hasMore = false;
        }
        return {code: response.status, body: allResult, len: allResult.length}
    }catch(error){
        console.log(error)
        return {code:500 , body:"Error in request"}
    }finally {
        loading = false;
    }
}



// resetFetchPost clears the paging guards so a new session (a filter change, a
// logout) starts fetching from the first page again instead of staying on the
// previous session's state.
export function resetFetchPost() {
    loading = false;
    hasMore = true;
}

// stopFetchPost halts pagination after a failed request (unknown cursor, server
// error...) so the infinite scroll does not keep re-requesting the same page.
export function stopFetchPost() {
    hasMore = false;
}

// Creates a post and returns the backend response. The new post id arrives as
// JSON on success; failures are plain-text bodies.
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

// "GET /api/posts

// The backend serves one post on its own resource:
//  GET /api/posts/{id} -> the same PostInfo shape a feed post carries
//                         (ID, Title, Content, Author, Categories, Likes,
//                         DisLikes, LikeInfo, CreatedAt). The session cookie
//                         decides the user, so LikeInfo already holds the
//                         reaction the signed in user has on this post.
// A non 200 answer is plain text (400 invalid id, 401 no session, 404 unknown
// post), so only a 200 is decoded as json.
export async function fetchPostDetails(postID) {

    try{
        const response = await fetch("/api/posts/" + postID , {
            method:"GET",
        })
        if (response.status != 200){
            return {code: response.status, body: await response.text()}
        }
        return {code: response.status, body: await response.json()}
    }catch(error){
        console.log(error)
        return {code:500 , body:"Error in request"}
    }
}