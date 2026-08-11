import { Categories } from "./Filter.js";
import { DEFAULT_CATEGORIES } from "../../services/categories.js";

export default function createPost(categories = DEFAULT_CATEGORIES) {
  return `
    <div class = "modal-overlay hidden">
        <div class =modal-container>
            <div class = "modal-header">
                <h1 class = "modal-title">Create Post</h1>
                <button class = "modal-close-btn" id = "cancelCreatePost">X</button>
            </div>
            <div class = "modal-body">
                <form class = "Create post form">
                    <input type="text" placeholder="enter Title" class = "modal-input">
                    <textarea id="message" maxlength="200" placeholder="write your message" class = "modal-textarea"></textarea>
                    ${Categories({ categories })}
                    <button class = "modal-submit-btn" type = "submit">submit</button>
                </form>
            </div>
        </div>
    </div>
    `;
}
