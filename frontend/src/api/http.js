export const base_link = import.meta.env.VITE_API_URL;


async function handleError(res, defaultMessage) {
    let message = defaultMessage

    try {
        const data = await res.json()
        if (data?.message) {
            message = data.message
        }
    } catch {
        // ignore JSON parse errors
    }

    switch (res.status) {
        case 400:
            throw new Error(message || "Некорректный запрос")
        case 401:
            throw new Error(message || "Unauthorized")
        case 403:
            throw new Error("Недостаточно прав")
        case 404:
            throw new Error("Ресурс не найден")
        case 500:
            throw new Error("Ошибка сервера. Попробуйте позже")
        default:
            throw new Error(message || `Ошибка (${res.status})`)
    }
}



// await handleError(res, "text")



export async function signUp(name, password) {
    const res = await fetch(`${base_link}/sign-up`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            name: name,
            password: password
        })
    })

    if (!res.ok) {
        await handleError(res, "Failed to sign up")
    }
    // return await res.json()
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
        await handleError(res, "Failed to log in")
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
        await handleError(res, "Failed to get all users")
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
        await handleError(res, "Failed to create new user")
    }
    // return await res.json()
}

export async function getUserByIDAdmin(userId){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users/${userId}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        await handleError(res, "Failed to get user by ID")
    }
    return await res.json()
}


export async function getUserTasksAdmin(userId){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users/${userId}/tasks`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
       await handleError(res, "Failed to get user's tasks")
    }
    return await res.json()
}


export async function createNewTaskAdmin(userId, title, description){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users/${userId}/tasks`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ title: title, description: description })
    })
    if (!res.ok) {
        await handleError(res, "Failed to create new task")
    }
    return await res.json()
}


export async function renameUserAdmin(userId, newName){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users/${userId}/rename`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ name: newName })
    })
    if (!res.ok) {
        await handleError(res, "Failed to rename user")
    }
    return await res.json()
}

export async function changeUserPasswordAdmin(userId, oldPassword, newPassword){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users/${userId}/password`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ old_password: oldPassword, new_password: newPassword })
    })
    if (!res.ok) {
        await handleError(res, "Failed to change user's password")
    }
    return await res.json()
}


export async function updateUserRoleAdmin(userId){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users/${userId}/role`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        await handleError(res, "Failed to update user's role")
    }
    return await res.json()
}


export async function deleteUserAdmin(userId){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/users/${userId}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        await handleError(res, "Failed to delete user")
    }
    return await res.json()
}

export async function getAllTasksAdmin(){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/tasks`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        await handleError(res, "Failed to get all tasks")
    }
    return await res.json()
}

export async function deleteTaskAdmin(taskId){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/tasks/${taskId}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        await handleError(res, "Failed to delete task")
    }
    return await res.json()
}


export async function switchTaskStatusAdmin(taskId){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/tasks/${taskId}/switch`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
    })
    if (!res.ok) {
        await handleError(res, "Failed to switch task status")
    }
    return await res.json()
}

export async function changeTaskTitleAdmin(taskId, newTitle){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/tasks/${taskId}/title`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ title: newTitle })
    })
    if (!res.ok) {
        await handleError(res, "Failed to change task title")
    }
    return await res.json()
}

export async function changeTaskDescriptionAdmin(taskId, newDescription){
    const token = localStorage.getItem("token")
    const res = await fetch(`${base_link}/admin/tasks/${taskId}/description`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ description: newDescription })
    })
    if (!res.ok) {
        await handleError(res, "Failed to change task description")
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
        await handleError(res, "Failed to load profile")
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
        await handleError(res, "Failed to rename profile")
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
        await handleError(res, "Failed to change password")
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
        await handleError(res, "Failed to delete profile")
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
        await handleError(res, "Failed to load your tasks")
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
        await handleError(res, "Failed to create new task")
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
        await handleError(res, "Failed to delete task")
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
        await handleError(res, "Failed to switch task status")
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
        await handleError(res, "Failed to change task title")
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
        await handleError(res, "Failed to change task description")
    }
    return await res.json()
}
