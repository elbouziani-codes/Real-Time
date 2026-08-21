export default function notFoundPage() {
    return `
        <section class="not-found">
            <div class="not-found__panel">
                <div class="not-found__code">404</div>
                <h1>Page not found</h1>
                <p>
                    The link you opened does not match any page in this app.
                    You can go back to the home page and continue from there.
                </p>
                <a class="not-found__button" href="/">Go Home</a>
            </div>
        </section>
    `;
}
