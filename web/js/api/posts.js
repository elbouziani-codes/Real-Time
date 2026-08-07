

export async function fetchPost() {

    try{
        const response = await fetch("/api/posts" , {
        method:"GET",
    })
        let allResult = await response.json()
        return allResult
    }catch(error){
        console.log(error)
    }
}


export async function fetchCreatePost(data) {

    try{
        const response = await fetch("/api/posts" , {
        method:"POST",
        body: JSON.stringify(data),
        })
        let allResult = await response.json()
        return {code:response.ok, body:allResult}
    }catch(error){
        console.log(error)
    }
}

// "GET /api/posts