
import { sendPostDetails, reactToPost, CURRENT_POST } from "./../services/posts.js"
import { sendComments, createComment } from "./../services/comments.js"
import { reactionState } from "./../components/home/post.js"
import PostDetails from "./../components/postDetails/PostDetails.js"
import Comments, { CommentsList } from "./../components/postDetails/Comments.js"
import { navigate } from "./../router/router.js"

const LIKE = ".post-details__action-btn--like";
const DISLIKE = ".post-details__action-btn--dislike";

function postID() {
    return new URLSearchParams(window.location.search).get("id") ?? "";
}

export default async function PostDetailsListener() {
    const main = document.querySelector(".post-details__main");
    if (!main) return;
    document.querySelector(".post-details")?.addEventListener("click", backListener);
    main.addEventListener("click", reactListener);
    main.addEventListener("submit", commentListener);

    const id = postID();
    const {post, error} = await sendPostDetails(id);
    if (!document.contains(main)) return;
    if (!post) {
        main.innerHTML = `<p class="error-state">${error}</p>`;
        return;
    }

    const commentsSkeleton = `
        <div class="skeleton-comments" role="status" aria-label="Loading comments">
            <div class="skeleton"></div>
            <div class="skeleton"></div>
            <div class="skeleton"></div>
        </div>`;
    main.innerHTML = PostDetails(post) + `<div class="post-details__comments-slot">${commentsSkeleton}</div>`;

    const comments = await sendComments(id);
    if (!document.contains(main)) return;
    main.querySelector(".post-details__comments-slot").innerHTML = Comments(comments.comments ?? []);
    if (comments.error) {
        setStatus(main, comments.error);
    }
}

function backListener(e) {
    if (!e.target.closest(".post-details__back")) return;
    navigate("/");
}
async function reactListener(e) {
    const button = e.target.closest(`${LIKE}, ${DISLIKE}`);
    if (!button) return;

    const actions = button.closest(".post-details__actions");
    const id = actions?.dataset.postId;
    if (!id) return;

    const isLike = button.matches(LIKE);
    const post = await reactToPost(id, isLike);
    if (!post || !document.contains(actions)) return;

    const likeBtn = actions.querySelector(LIKE);
    const disLikeBtn = actions.querySelector(DISLIKE);
    const {liked, disliked} = reactionState(post);
    likeBtn.querySelector(".post-details__action-count").textContent = post.Likes;
    disLikeBtn.querySelector(".post-details__action-count").textContent = post.DisLikes;
    likeBtn.classList.toggle("active", liked);
    disLikeBtn.classList.toggle("active", disliked);
}

async function commentListener(e) {
    const form = e.target.closest(".post-details__comment-form");
    if (!form) return;
    e.preventDefault();

    const main = form.closest(".post-details__main");
    const input = form.querySelector(".post-details__comment-input");
    const submit = form.querySelector(".post-details__comment-submit");
    const id = CURRENT_POST?.ID;
    const content = input.value.trim();

    if (!id) return;
    if (content == "") {
        setStatus(main, "Write something before posting.");
        return;
    }

    submit.disabled = true;
    setStatus(main, "");
    const {comments, error} = await createComment(id, content);
    if (!document.contains(main)) return;
    submit.disabled = false;

    if (error) {
        setStatus(main, error);
        return;
    }

    input.value = "";
    const list = main.querySelector(".post-details__comments-list");
    const count = main.querySelector(".post-details__comments-count");
    if (list) list.innerHTML = CommentsList(comments);
    if (count) count.textContent = `(${comments.length})`;
}

function setStatus(main, message) {
    const status = main?.querySelector(".post-details__comment-status");
    if (status) status.textContent = message;
}
