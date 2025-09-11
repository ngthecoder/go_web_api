"use client"

import { User } from "@/lib/types";
import { useEffect, useState } from "react";

export default function ProfilePage() {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const userData = localStorage.getItem('user');
        const token = localStorage.getItem('token');

        if (userData && token) {
            setUser(JSON.parse(userData));
        }
        setIsLoading(false);
    }, []);

    const handleLogout = () => {
        localStorage.removeItem('user');
        localStorage.removeItem('token');
        window.location.href = '/login';
    };

    if (isLoading) {
        return <div className="flex justify-center items-center min-h-screen">Loading...</div>;
    }

    if (!user) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-center">
                    <p className="mb-4">Please log in to view your profile</p>
                    <a href="/login" className="text-blue-600 hover:underline">Go to Login</a>
                </div>
            </div>
        );
    }

    return (
        <div className="py-8">
            <div className="max-w-4xl mx-auto px-4">
                <div className="bg-white rounded-lg p-6">
                    <div className="flex justify-between items-start mb-6">
                        <h1 className="text-3xl font-bold text-gray-900">Profile</h1>
                        <button 
                            onClick={handleLogout}
                            className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                        >
                            Logout
                        </button>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Username</label>
                            <p className="mt-1 text-lg text-gray-900">{user.username}</p>
                        </div>
                        
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Email</label>
                            <p className="mt-1 text-lg text-gray-900">{user.email}</p>
                        </div>
                        
                        <div>
                            <label className="block text-sm font-medium text-gray-700">User ID</label>
                            <p className="mt-1 text-sm text-gray-600 font-mono">{user.id}</p>
                        </div>
                        
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Member Since</label>
                            <p className="mt-1 text-lg text-gray-900">
                                {new Date(user.created_at).toLocaleDateString()}
                            </p>
                        </div>
                    </div>

                    <div className="mt-8 pt-6 border-t">
                        <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
                        <div className="space-y-2">
                            <button className="block w-full text-left px-4 py-2 bg-blue-50 hover:bg-blue-100 rounded">
                                View My Liked Recipes (Coming Soon)
                            </button>
                            <button className="block w-full text-left px-4 py-2 bg-green-50 hover:bg-green-100 rounded">
                                My Shopping Lists (Coming Soon)
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}