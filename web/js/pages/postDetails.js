import { Navbar } from "../components/navbar.js";
import { me } from "../services/me.js";

export default function pagePost() {
    return `
        <section class="post-details">
            ${Navbar(me)}
            <div class="post-details__layout post-details__layout--single">
                <article class="post-details__content">
                    <button class="post-details__back">← Back to Home</button>
                    <div class="post-details__main">
                        <div class="skeleton-post" role="status" aria-label="Loading post">
                            <div class="skeleton skeleton--title"></div>
                            <div class="skeleton skeleton--meta"></div>
                            <div class="skeleton skeleton--line"></div>
                            <div class="skeleton skeleton--line"></div>
                            <div class="skeleton skeleton--short"></div>
                        </div>
                    </div>
                </article>
            </div>
        </section>
    `;
}
