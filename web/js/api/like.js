

// The backend exposes one reaction resource for likes and dislikes:
//  POST  /api/reactions        {parent_id, is_like}  -> the new reaction id
//  PATCH /api/reactions/{id}   {is_like}             -> switches the reaction,
//                                                       or deletes it when
//                                                       is_like already matches
export async function fetchReact(parentID, isLike) {
    try {
        const response = await fetch("/api/reactions", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ parent_id: parentID, is_like: isLike }),
        });
        if (response.status != 200) {
            return { code: response.status, body: await response.text() };
        }
        return { code: response.status, body: await response.json() };
    } catch {
        return { code: 500, body: "Error in request" };
    }
}


export async function fetchUpdateReact(reactionID, isLike) {
    try {
        const response = await fetch("/api/reactions/" + reactionID, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ is_like: isLike }),
        });
        return { code: response.status, body: await response.text() };
    } catch {
        return { code: 500, body: "Error in request" };
    }
}
