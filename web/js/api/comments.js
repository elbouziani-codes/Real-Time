

// The backend serves the comments of a post under the post's own id:
//  GET  /api/comments/{post_id} -> the post's comments, or null when it has none
//  POST /api/comments           {parent_id, content} -> the new comment id
export async function fetchComments(postID) {
    try {
        const response = await fetch("/api/comments/" + postID, {
            method: "GET",
        });
        if (response.status != 200) {
            return { code: response.status, body: await response.text() };
        }
        return { code: response.status, body: await response.json() };
    } catch {
        return { code: 500, body: "Error in request" };
    }
}

export async function fetchCreateComment(postID, content) {
    try {
        const response = await fetch("/api/comments", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ parent_id: postID, content: content }),
        });
        if (response.status != 200) {
            return { code: response.status, body: await response.text() };
        }
        return { code: response.status, body: await response.json() };
    } catch {
        return { code: 500, body: "Error in request" };
    }
}
