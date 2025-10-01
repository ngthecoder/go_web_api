'use client';

import ProtectedRoute from '@/components/ProtectedRoute';

export default function LikedRecipesPage() {
  return (
    <ProtectedRoute>
      <div className="container mx-auto p-4">
        <h1 className="text-3xl font-bold mb-6">My Liked Recipes</h1>
        <div className="bg-white rounded-lg shadow-sm border p-6">
          <p className="text-gray-600 text-center">
            You'll be able to see all your favorite recipes here!
          </p>
        </div>
      </div>
    </ProtectedRoute>
  );
}
