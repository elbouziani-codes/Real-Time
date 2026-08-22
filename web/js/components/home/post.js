import { ZERO_UUID } from "./../../config/config.js";
import formatDateTime from "./../../utils/time.js";
import { escapeHTML } from "./../../utils/helpers.js";

export function reactionState(post = {}) {
    const { LikeInfo = {} } = post;
    const reacted = Boolean(LikeInfo.ID) && LikeInfo.ID != ZERO_UUID;
    return {
        liked: reacted && LikeInfo.IsLike === true,
        disliked: reacted && LikeInfo.IsLike === false,
    };
}

function truncateText(text = "", maxLength = 80) {
    const value = String(text).trim();
    if (value.length <= maxLength) {
        return value;
    }
    return value.slice(0, maxLength).trimEnd() + "...";
}

export default function PostCard(post = {}) {
    if (typeof post == "string") {
        return `
        <article class="post-card post-card--state">
            <h3>${escapeHTML(post)}</h3>
        </article>
        `;
    }
    const {
        ID = {},
        Title = '',
        Author = {},
        Content = '',
        Categories = [],
        Likes = 0,
        DisLikes = 0,
        CreatedAt = '',
    } = post;
    const createdAt = formatDateTime(CreatedAt);
    const { liked, disliked } = reactionState(post);
    const likeClass = liked ? " active" : "";
    const disLikeClass = disliked ? " active" : "";
    const previewTitle = truncateText(Title, 80);
    const previewContent = truncateText(Content, 80);

    return `
        <article class="post-card" data-post-id="${ID ?? ''}">
            <div class="post-header">
                <div>
                    <h3>${escapeHTML(previewTitle)}</h3>
                    <span>by <span class="user-info__name">${escapeHTML(Author.NickName ?? '')}</span></span>
                </div>
                <span class="post-time">${createdAt}</span>
            </div>
            <p>${escapeHTML(previewContent)}</p>
            <div class="post-categories"> ${Categories.map((e) => `<span class="post-category">${escapeHTML(e.title ?? '')}</span>`).join("")} </div>
            <div class="post-actions">
                <button class="like-btn${likeClass}">👍 ${Likes}</button>
                <button class="comment-btn${disLikeClass}">👎 ${DisLikes}</button>
            </div>
        </article>
    `;
}
