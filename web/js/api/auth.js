

export async function fetchMe() {
    try {
        const response = await fetch("/api/me", {
            method: "GET",
        });
        if(response.status == 200){
            return {code:response.status, body: await response.json()}
        }else{
            return {code:response.status, body: "errors"}
        }
    } catch (error) {
        console.log("Failed to get profile:", error);
    }
}



export async function fetchLogin(credentials) {

    try {

        const response = await fetch("/api/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(credentials),
        });


        const text = await response.text();


        let data;

        try {
            data = JSON.parse(text);
        } catch {
            data = text;
        }

        return {ok: response.ok, body: data};


    } catch(error){

        console.error("login request failed:", error);

        return {
            ok:false,
            data:null
        };
    }
}


export async function fetchRegister(credentials) {
    try {

        const response = await fetch("/api/register", {
        method: "POST",
        headers:{
            "Content-Type":"application/json",
        },
        credentials:"include",
        body:JSON.stringify(credentials),
        });


        const text = await response.text();



        return response.ok;
    } catch (error) {
        console.error("Register request failed:", error);
        return false;
    }
}
export async function fetchLogout() {
    try {
        const response = await fetch("/api/logout", {
            method: "POST",
            credentials: "include",
        });
        return {code: response.status, body: await response.text()};
    } catch (error) {
        console.log("Logout request failed:", error);
        return {code: 500, body: "Error in request"};
    }
}