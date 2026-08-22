

import {fetchComments, fetchCreateComment} from "./../api/comments.js"
import { navigate } from "./../router/router.js"

export async function sendComments(postID) {
    const response = await fetchComments(postID);
    if (response.code == 200) {
        return {comments: Array.isArray(response.body) ? response.body : []};
    }

    if (response.code == 401) {
        navigate("/login");
        return {error: "You need to sign in to read the comments."};
    }
    (response.body);
    return {error: "Could not load the comments."};
}

export async function createComment(postID, content) {
    const response = await fetchCreateComment(postID, content);
    if (response.code == 200) {
        return await sendComments(postID);
    }

    if (response.code == 401) {
        navigate("/login");
        return {error: "You need to sign in to comment."};
    }
    if (response.code == 400) {
        return {error: typeof response.body == "string" && response.body ? response.body.trim() : "This comment is not valid."};
    }
    (response.body);
    return {error: "Could not post this comment. Please try again."};
}
