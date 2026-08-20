export default function endPage() {
    return `
        <section class="session-ended">
            <div class="session-ended__panel">
                <div class="session-ended__icon">💬</div>
                <h1>Connection moved to another tab</h1>
                <p>
                    This session was closed because your account was opened in a newer tab.
                    To keep messages in sync, only one active connection is allowed at a time.
                </p>
                <button class="session-ended__button" type="button">Back to Home</button>
            </div>
        </section>
    `;
}
