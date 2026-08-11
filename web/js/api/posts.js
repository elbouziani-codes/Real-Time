let loading = false;
let hasMore = true;

export async function fetchPost(offset) {
    if (loading || !hasMore) return{code: 1000, body: [], len: 0};

    loading = true;

    try{
        const response = await fetch("/api/posts?offset="+offset , {
        method:"GET",
    })
        let allResult = await response.json()

        if (!Array.isArray(allResult)) {
            hasMore = false;
            return {code: 1000, body: [], len: 0};
        }
        return {code: response.status, body: allResult, len: allResult.length}
    }catch(error){
        console.log(error)
        return {code:500 , body:"Error in request"}
    }finally {
        loading = false;
    }
}



// resetFetchPost clears the paging guards so a new session starts fetching from
// the first page again instead of staying on the previous session's state.
export function resetFetchPost() {
    loading = false;
    hasMore = true;
}

export async function fetchCreatePost(data) {

    try{
        const response = await fetch("/api/posts" , {
        method:"POST",
        body: JSON.stringify(data),
        })
        //let allResult = await response.json()
        return {code:response.status, body:"allResult"}
    }catch(error){
        console.log(error)
        return {code:500 , body:"Error in request"}
    }
}

// "GET /api/posts

// The backend serves one post on its own resource:
//  GET /api/posts/{id} -> the same PostInfo shape a feed post carries
//                         (ID, Title, Content, Author, Categories, Likes,
//                         DisLikes, LikeInfo, CreatedAt). The session cookie
//                         decides the user, so LikeInfo already holds the
//                         reaction the signed in user has on this post.
// A non 200 answer is plain text (400 invalid id, 401 no session, 404 unknown
// post), so only a 200 is decoded as json.
export async function fetchPostDetails(postID) {

    try{
        const response = await fetch("/api/posts/" + postID , {
            method:"GET",
        })
        if (response.status != 200){
            return {code: response.status, body: await response.text()}
        }
        return {code: response.status, body: await response.json()}
    }catch(error){
        console.log(error)
        return {code:500 , body:"Error in request"}
    }
}