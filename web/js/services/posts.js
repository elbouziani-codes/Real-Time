import {fetchPost, fetchPostDetails, resetFetchPost} from "./../api/posts.js"
import {fetchReact, fetchUpdateReact} from "./../api/like.js"
import {config, ZERO_UUID} from "./../config/config.js"
import PostCard  from "./../components/home/post.js"
import { navigate } from "./../router/router.js"


export let DEFAULT_POST = [];

// Post ids with a reaction request in flight, so a double click does not send
// the same reaction twice.
const reacting = new Set();

export async function sendAllPost(Scroll = false) {
    let response = await fetchPost(config.offsetPost);
    
    console.log(response)
    if (response.code == 200){
        if (!Array.isArray(response.body)) {
            if (DEFAULT_POST.length === 0) {
                DEFAULT_POST = ["not found post"];
            }
            return ;
        }
        DEFAULT_POST.push(...response.body);
        config.offsetPost += response.len
        if (Scroll){
            const postsHtml = response.body.map((e) => {return PostCard(e) }).join('');
            document.querySelector(".feed").innerHTML += postsHtml
        }
    }else{
        DEFAULT_POST = ["error in fetch Categories"]
    }
}


export function resetPosts() {
    DEFAULT_POST = [];
    config.offsetPost = 0;
    resetFetchPost();
}


export async function reactToPost(postID, isLike) {
    const post = findPost(postID);
    if (!post || reacting.has(postID)) return null;

    reacting.add(postID);
    try {
        const reactionID = post.LikeInfo?.ID?.Value;
        if (reactionID == ZERO_UUID) {
            const response = await fetchReact(postID, isLike);
            if (response.code != 200) return failedReaction(response);
            post.LikeInfo = {ID: {Value: response.body.Value}, IsLike: isLike};
            countReaction(post, isLike, 1);
            return post;
        }

        const response = await fetchUpdateReact(reactionID, isLike);
        if (response.code != 200) return failedReaction(response);
        if (post.LikeInfo.IsLike == isLike) {
            post.LikeInfo = {ID: {Value: ZERO_UUID}, IsLike: false};
            countReaction(post, isLike, -1);
        } else {
            post.LikeInfo.IsLike = isLike;
            countReaction(post, isLike, 1);
            countReaction(post, !isLike, -1);
        }
        return post;
    } finally {
        reacting.delete(postID);
    }
}

function countReaction(post, isLike, delta) {
    if (isLike) {
        post.Likes = post.Likes + delta;
    } else {
        post.DisLikes = post.DisLikes + delta
    }
}

function failedReaction(response) {
    console.log(response.body);
    if (response.code == 401) {
        navigate("/login");
    }
    return null;
}


export let CURRENT_POST = null;

function findPost(postID) {
    const post = DEFAULT_POST.find((e) => typeof e != "string" && e.ID?.Value == postID);
    if (post) return post;
    if (CURRENT_POST && CURRENT_POST.ID?.Value == postID) return CURRENT_POST;
    return null;
}

// sendPostDetails loads one post through GET /api/posts/{id} and caches it as
// CURRENT_POST. Returns {post} on success, or {error} with a message ready to
// show in the page's error state.
export async function sendPostDetails(postID) {
    CURRENT_POST = null;
    if (!postID) return {error: "This post link is missing a post id."};

    const response = await fetchPostDetails(postID);
    if (response.code == 200 && response.body && typeof response.body == "object") {
        CURRENT_POST = response.body;
        return {post: CURRENT_POST};
    }

    if (response.code == 401) {
        navigate("/login");
        return {error: "You need to sign in to read this post."};
    }
    if (response.code == 400) return {error: "This post id is not valid."};
    if (response.code == 404) return {error: "This post does not exist."};
    console.log(response.body);
    return {error: "Could not load this post. Please try again."};
}