import { Navbar } from "../components/navbar.js";
import { me } from "../services/me.js";

/**
 * pagePost — /postDetails?id=POST_ID
 * Top-level composition: Navbar + the post details shell.
 *
 * Only the shell is rendered here. The post itself is loaded by
 * PostDetailsListener, so the page paints a loading state first and swaps in the
 * post, or an error state, once the request answered.
 */

export default function pagePost() {
    return `
        <section class="post-details">
            ${Navbar(me)}
            <div class="post-details__layout post-details__layout--single">
                <article class="post-details__content">
                    <button class="post-details__back">← Back to Home</button>
                    <div class="post-details__main">
                        <p class="loading">Loading post…</p>
                    </div>
                </article>
            </div>
        </section>
    `;
}
