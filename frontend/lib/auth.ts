import { AuthResponse, LoginRequest, RegisterRequest } from "./types"

export async function registerUser(data: RegisterRequest):Promise<AuthResponse> {
    const response = await fetch("http://localhost:8000/api/auth/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })

    if (!response.ok) {
        const error = await response.text()
        throw new Error(error);
    }

    const responseJSON: AuthResponse = await response.json()

    return responseJSON
}

export async function loginUser(data: LoginRequest):Promise<AuthResponse> {
    const response = await fetch("http://localhost:8000/api/auth/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })

    if (!response.ok) {
        const error = await response.text()
        throw new Error(error);
    }

    const responseJSON: AuthResponse = await response.json()

    return responseJSON
}