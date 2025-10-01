'use client';

import ProtectedRoute from '@/components/ProtectedRoute';
import Link from 'next/link';

export default function ChangePasswordPage() {
  return (
    <ProtectedRoute>
      <div className="container mx-auto p-4 max-w-2xl">
        <div className="bg-white rounded-lg shadow-sm border p-6">
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Change Password</h1>
            <Link 
              href="/profile"
              className="text-gray-600 hover:text-gray-900 transition-colors"
            >
              ← Back to Profile
            </Link>
          </div>
          
          <div className="text-center py-12">
            <div className="bg-green-50 rounded-lg p-8">
              <svg className="mx-auto h-16 w-16 text-green-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Change Password Feature</h3>
              <p className="text-gray-600 mb-4">
                You'll be able to update your account password here.
              </p>
              <p className="text-sm text-gray-500">
                Backend endpoint needed: <code className="bg-gray-100 px-2 py-1 rounded">PUT /api/user/password</code>
              </p>
            </div>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
