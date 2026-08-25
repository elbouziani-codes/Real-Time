import { fetchPost, fetchPostDetails, resetFetchPost, stopFetchPost } from "./../api/posts.js";
import { fetchReact, fetchUpdateReact } from "./../api/like.js";
import { config, ZERO_UUID } from "./../config/config.js";
import PostCard from "./../components/home/post.js";
import { navigate } from "./../router/router.js";

export let DEFAULT_POST = [];

const reacting = new Set();

let feedRequestId = 0;

export async function sendAllPost() {
    const requestId = feedRequestId;
    let posts = await fetchPost(config.postsCursor, config.postsFilters);

    if (requestId != feedRequestId) return;

    if (posts.length > 0) {
        
        if (DEFAULT_POST.length === 1 && typeof DEFAULT_POST[0] == "string") {
            DEFAULT_POST = []
        }
        const fresh = appendPosts(posts);
            const postsHtml = fresh.map((e) => PostCard(e)).join("");
            if (document.querySelector(".feed")) {
                document.querySelector(".feed").innerHTML += postsHtml;
            }
    } else {
        if (DEFAULT_POST.length === 0) {
            DEFAULT_POST[0] = "no more posts"
            if (document.querySelector(".feed")) {
                document.querySelector(".feed").innerHTML = PostCard("no more posts");
            }
        }
        stopFetchPost();
    }
}

function appendPosts(newPosts) {
    const known = new Set(DEFAULT_POST.map((e) => (e.ID)));
    const fresh = [];
    for (const post of newPosts) {
        if (post && !known.has(post.ID)) {
            known.add(post.ID);
            DEFAULT_POST.push(post);
            fresh.push(post);
        }
    }
    const last = newPosts[newPosts.length - 1];
    if (last && last.ID) {
        config.postsCursor = last.ID;
    }
    return fresh;
}

export function resetPosts() {
    DEFAULT_POST = [];
    config.postsCursor = null;
    resetFetchPost();
    feedRequestId++;
}

export function applyPostFilters(filters) {
    config.postsFilters = filters;
    document.querySelector(".feed").innerHTML = "";
    resetPosts();
}

export async function reactToPost(postID, isLike) {
    const post = findPost(postID);
    if (!post || reacting.has(postID)) return null;

    reacting.add(postID);
    try {
        const reactionID = post.LikeInfo ?.ID;
        if (reactionID == ZERO_UUID) {
            const response = await fetchReact(postID, isLike);
            if (response.code != 200) return failedReaction(response);
            post.LikeInfo = { ID: response.body, IsLike: isLike };
            countReaction(post, isLike, 1);
            return post;
        }

        const response = await fetchUpdateReact(reactionID, isLike);
        if (response.code != 200) return failedReaction(response);
        if (post.LikeInfo.IsLike == isLike) {
            post.LikeInfo = { ID: ZERO_UUID, IsLike: false };
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
        post.DisLikes = post.DisLikes + delta;
    }
}

function failedReaction(response) {
    if (response.code == 401) {
        navigate("/login");
    }
    return null;
}

export let CURRENT_POST = null;

function findPost(postID) {
    const post = DEFAULT_POST.find((e) => e.ID == postID);
    if (post) return post;
    if (CURRENT_POST && CURRENT_POST.ID == postID) return CURRENT_POST;
    return null;
}

export async function sendPostDetails(postID) {
    CURRENT_POST = null;
    if (!postID) return { error: "This post link is missing a post id." };

    const response = await fetchPostDetails(postID);
    if (response.code == 200 && response.body && typeof response.body == "object") {
        CURRENT_POST = response.body;
        return { post: CURRENT_POST };
    }

    if (response.code == 401) {
        navigate("/login");
        return { error: "You need to sign in to read this post." };
    }
    if (response.code == 400) return { error: "This post id is not valid." };
    if (response.code == 404) return { error: "This post does not exist." };
    return { error: "Could not load this post. Please try again." };
}