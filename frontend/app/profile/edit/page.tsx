'use client';

import ProtectedRoute from '@/components/ProtectedRoute';
import Link from 'next/link';

export default function EditProfilePage() {
  return (
    <ProtectedRoute>
      <div className="container mx-auto p-4 max-w-2xl">
        <div className="bg-white rounded-lg shadow-sm border p-6">
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Edit Profile</h1>
            <Link 
              href="/profile"
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              ← Back to Profile
            </Link>
          </div>
          
          <div className="text-center py-12">
            <div className="bg-blue-50 rounded-lg p-8">
              <svg className="mx-auto h-16 w-16 text-blue-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Edit Profile Feature</h3>
              <p className="text-gray-600 mb-4">
                You'll be able to update your username and email here.
              </p>
              <p className="text-sm text-gray-500">
                Backend endpoint needed: <code className="bg-gray-100 px-2 py-1 rounded">PUT /api/user/profile</code>
              </p>
            </div>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
