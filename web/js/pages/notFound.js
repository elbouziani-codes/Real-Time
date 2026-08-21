export default function ErrorPage({code = "404", title = "Page not found", content = "The link you opened does not match any page in this app.You can go back to the home page and continue from there."}) {
    return `
        <section class="not-found">
            <div class="not-found__panel">
                <div class="not-found__code">${code}</div>
                <h1>${title}</h1>
                <p>${content}</p>
                <a class="not-found__button" href="/">Go Home</a>
            </div>
        </section>
    `;
}
