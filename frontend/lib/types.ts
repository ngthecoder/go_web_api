export interface User {
    id: string;
    username: string;
    email: string;
    password: string;
    created_at: Date;
    updated_at: Date;
}

export interface RegisterRequest {
    username: string;
    email: string;
    password: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface AuthResponse {
    user: User;
    token: string;
}
