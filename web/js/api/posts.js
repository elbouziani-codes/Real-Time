

export default async function fetchPost() {

    try{
        const response = await fetch("/api/posts" , {
        method:"GET",
    })
        let allResult = await response.json()
        console.log(allResult)
        return allResult
    }catch(error){
        console.log(error)
    }
}

// "GET /api/posts