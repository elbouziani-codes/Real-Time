import { navigate } from "../router/router.js";

export default function endListener() {
    document.querySelector(".session-ended__button")?.addEventListener("click", () => {
        navigate("/");
    });
}
