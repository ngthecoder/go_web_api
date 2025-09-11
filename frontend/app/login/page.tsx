'use client'

import { useState } from "react";
import { AuthResponse, LoginRequest } from "@/lib/types"
import { loginUser } from "@/lib/auth"

export default function LoginPage() {

    const [isLoading, setLoading] = useState(false)
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    const [formData, setFormData] = useState<LoginRequest>({
        email: "",
        password: ""
    })

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        console.log('Form submitted!');
        
        try {
            setLoading(true)

            const response:AuthResponse = await loginUser(formData)

            localStorage.setItem('token', response.token);
            localStorage.setItem('user', JSON.stringify(response.user));
            
            setMessage("Login completed")
            setFormData({email: "", password: "" });
        }
        catch(error: any) {
            setError(error.message)
        }
        finally {
            setLoading(false)
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    return (
        <div className="h-screen flex items-center justify-center bg-gray-50 overflow-hidden">
            <div className="w-full max-w-2xl mx-auto px-8">
                <h2 className="text-center text-3xl font-extrabold text-gray-900">
                    Log into your account
                </h2>

                {message && (
                    <div className="mt-6 mb-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded-lg">
                        <div className="flex">
                            <div className="text-green-500 mr-3">✓</div>
                            <div>{message}</div>
                        </div>
                    </div>
                )}

                {error && (
                    <div className="mt-6 mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded-lg">
                        <div className="flex">
                            <div className="text-red-500 mr-3">✕</div>
                            <div>{error}</div>
                        </div>
                    </div>
                )}
        
                <form onSubmit={handleSubmit} className="space-y-6">
                    <div className="space-y-4">
                        <div>
                            <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                                Email
                            </label>
                            <input
                                id="email"
                                name="email"
                                type="email"
                                value={formData.email}
                                onChange={handleChange}
                                required
                                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                placeholder="Enter email"
                            />
                        </div>
                        
                        <div>
                            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                                Password
                            </label>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                value={formData.password}
                                onChange={handleChange}
                                required
                                className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                placeholder="Enter password"
                            />
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        className="w-full py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                    >
                        {isLoading? "Registering..." : "Register"}
                    </button>
                </form>
            </div>
        </div>
    );
}