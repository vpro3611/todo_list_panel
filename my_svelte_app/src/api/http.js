const base_link = "http://localhost:8080"


export async function signIn(name, password) {
    const res = await fetch(`${base_link}/sign-up`, {
        method: "POST",
        headers: { "Content-type": "application/json" },
        body: JSON.stringify({
            name: name,
            password: password
        })
    })

    if (!res.ok) {
        throw new Error(`status ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}   


export async function login(name, password) {
    const res = await fetch(`${base_link}/login`, {
        method: "POST",
        headers: {"Content-type": "application/json"},
        body: JSON.stringify({
            name: name,
            password: password
        })
    })
    if (!res.ok) {
        throw new Error(`login failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}