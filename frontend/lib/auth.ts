import { AuthResponse, LoginRequest, RegisterRequest } from "./types"
import { API_ENDPOINTS } from "./api-config"

export async function registerUser(data: RegisterRequest):Promise<AuthResponse> {
    try {
        const response = await fetch(API_ENDPOINTS.register, {
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
    } catch(error) {
        if (error instanceof Error) {
            throw error
        } else {
            throw new Error("Network error occurred")
        }
    }
}

export async function loginUser(data: LoginRequest):Promise<AuthResponse> {
    try {
        const response = await fetch(API_ENDPOINTS.login, {
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
    } catch(error) {
        if (error instanceof Error) {
            throw error
        } else {
            throw new Error("Network error occurred")
        }
    }
}