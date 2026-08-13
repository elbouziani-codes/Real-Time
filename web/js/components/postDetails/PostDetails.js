import UserAvatar from "./../UserAvatar.js";
import createUser from "./../../models/User.js";
import { reactionState } from "./../home/post.js";
import formatDateTime from "./../../utils/time.js";
import { escapeHTML } from "./../../utils/helpers.js";

export function authorModel(author = {}) {
    const { NickName = "", FirstName = "", LastName = "" } = author;
    const fullName = [FirstName, LastName].filter(Boolean).join(" ");
    return createUser({
        name: fullName || NickName,
        handle: NickName ? "@" + NickName : "",
        letter: (NickName || fullName).charAt(0),
    });
}

function renderBody(content) {
    return content
        .split("\n")
        .filter((line) => line.trim() != "")
        .map((line) => `<p class="post-details__text">${escapeHTML(line)}</p>`)
        .join("");
}

export default function PostDetails(post = {}) {
    const {
        ID = {},
        Title = '',
        Author = {},
        Content = '',
        Categories = [],
        Likes = 0,
        DisLikes = 0,
        CreatedAt = 0,
    } = post;

    const author = authorModel(Author);
    const { liked, disliked } = reactionState(post);
    const likeClass = liked ? " active" : "";
    const disLikeClass = disliked ? " active" : "";

    return `
        <header class="post-details__header">
            <h1 class="post-details__title">${escapeHTML(Title)}</h1>
            <div class="post-details__meta">
                <div class="post-details__author">
                    ${UserAvatar(author, { sizeClass: 'user-avatar--lg', className: 'post-details__author-avatar' })}
                    <div class="post-details__author-info user-info">
                        <span class="post-details__author-name user-info__name">${escapeHTML(author.name)}</span>
                        <span class="post-details__author-handle user-info__handle">${escapeHTML(author.handle)}</span>
                    </div>
                </div>
                <div class="post-details__date">
                    <span>📅</span>
                    <span>${formatDateTime(CreatedAt)}</span>
                </div>
            </div>
            <div class="post-details__categories">
                ${Categories.map((e) => `<span class="post-details__category">${escapeHTML(e.title)}</span>`).join("")}
            </div>
        </header>

        <div class="post-details__body">
            ${renderBody(Content)}
        </div>

        <div class="post-details__actions" data-post-id="${ID.Value ?? ''}">
            <button class="post-details__action-btn post-details__action-btn--like${likeClass}">
                <span>👍</span>
                <span class="post-details__action-count">${Likes}</span>
            </button>
            <button class="post-details__action-btn post-details__action-btn--dislike${disLikeClass}">
                <span>👎</span>
                <span class="post-details__action-count">${DisLikes}</span>
            </button>
        </div>
    `;
}
