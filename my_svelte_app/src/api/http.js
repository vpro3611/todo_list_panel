const base_link = "http://localhost:8080"


export async function signIn(name, password) {
    const res = await fetch(`${base_link}/sign-up`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
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
        headers: {"Content-Type": "application/json"},
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

// ADMIN FUNCTIONS
export async function getAllUsersAdmin(){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        throw new Error(`getAllUsers failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}

export async function createNewUserAdmin(name, password){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify( {
            name: name,
            password: password
        })
    })
    if (!res.ok) {
        throw new Error(`createNewUser failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}


// SUBJUCTIVE TO ENLARGEMENT
// { ADMIN FUNCTIONS END }


export async function getMe() {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        throw new Error(`getMe failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}


export async function renameMe(newName) {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/rename`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ name: newName })
    })
    if (!res.ok) {
        throw new Error(`renameMe failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}

export async function changeMyPassword(oldPassword, newPassword) {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/password`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
    })
    if (!res.ok) {
        throw new Error(`changeMyPassword failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}

export async function deleteMe() {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        throw new Error(`deleteMe failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}

export async function getMyTasks() {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/tasks`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        throw new Error(`getMyTasks failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}


export async function createNewTask(title, description) {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/tasks`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ title: title, description: description })
    })
    if (!res.ok) {
        throw new Error(`createNewTask failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}


export async function deleteTask(taskId) {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/tasks/${taskId}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        throw new Error(`deleteTask failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}

export async function switchTaskStatus(taskId) {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/tasks/${taskId}/switch`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        throw new Error(`switchTaskStatus failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}

export async function changeTaskTitle(taskId, newTitle) {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/tasks/${taskId}/title`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ title: newTitle })
    })
    if (!res.ok) {
        throw new Error(`changeTaskTitle failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}


export async function changeTaskDescription(taskId, newDescription) {
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/me/tasks/${taskId}/description`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ description: newDescription })
    })
    if (!res.ok) {
        throw new Error(`changeTaskDescription failed ${res.status} -> ${res.statusText}`)
    }
    return await res.json()
}