

export default async function getProfile() {
    try {
        const response = await fetch("http://localhost:8081/api/me", {
            method: "GET",
            credentials: "include",
        });

        console.log(response.status)
        if (response.ok) {
            return await response.json();
        }

        window.location.pathname = "/auth";
        return null;

    } catch (error) {
        console.error("Failed to get profile:", error);
        return null;
    }
}



// async function postLogin(credentials) {
//     const response = await fetch("http://localhost:8080/api/login", {
//         method: "POST",
//         headers: {
//             "Content-Type": "application/json",
//         },
//         credentials: "include",
//         body: JSON.stringify(credentials),
//     });

//     return await response.json();
// }

// async function logout() {
//     await fetch("/api/logout", {
//         method: "POST",
//         credentials: "include",
//     });
// }