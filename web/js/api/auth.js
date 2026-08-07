import { navigate } from "../router/router.js";


export async function getProfile(path) {
    try {
        const response = await fetch("/api/me", {
            method: "GET",
        });

        console.log(response.status)
        if (response.ok) {
            if (path == "/login" || path == "/register" ){
                navigate("/");
            }
            return await response.json();
        }
        if (path != "/login" && path != "/register" ){
            navigate("/login");
        } 
        return null;

    } catch (error) {
        console.error("Failed to get profile:", error);
        return null;
    }
}



export async function postLogin(credentials) {

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


        return {
            ok: response.ok,
            data
        };


    } catch(error){

        console.error("login request failed:", error);

        return {
            ok:false,
            data:null
        };
    }
}


export async function postRegister(credentials) {
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

        console.log("SERVER RESPONSE:", text);


        return response.ok;
    } catch (error) {
        console.error("Register request failed:", error);
        return false;
    }
}
async function logout() {
    await fetch("/api/logout", {
        method: "POST",
        credentials: "include",
    });
}