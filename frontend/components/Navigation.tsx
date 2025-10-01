'use client';
import Link from 'next/link';
import { useState, useRef, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';

export default function Navigation() {
  const { isAuthenticated, user, logout } = useAuth();

  const [showDropdown, setShowDropdown] = useState(false);

  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setShowDropdown(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleLogout = () => {
    setShowDropdown(false);
    logout();
  };

  return (
    <nav className="bg-white shadow-sm border-b">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <Link href="/" className="text-xl font-bold text-gray-900 hover:text-blue-600">
            Recipe API
          </Link>

          <div className="hidden md:flex space-x-6">
            <Link href="/" className="text-gray-700 hover:text-blue-600 transition-colors">
              Recipes
            </Link>
            <Link href="/ingredients" className="text-gray-700 hover:text-blue-600 transition-colors">
              Ingredients
            </Link>
            <Link href="/find-recipes" className="text-gray-700 hover:text-blue-600 transition-colors">
              Find by Ingredients
            </Link>

            {isAuthenticated && (
              <Link href="/liked-recipes" className="text-gray-700 hover:text-blue-600 transition-colors">
                My Liked Recipes
              </Link>
            )}
          </div>

          <div className="flex items-center space-x-4">
            {!isAuthenticated ? (
              <>
                <Link 
                  href="/login" 
                  className="text-gray-700 hover:text-blue-600 transition-colors"
                >
                  Login
                </Link>
                <Link 
                  href="/register" 
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  Register
                </Link>
              </>
            ) : (
              <div className="relative" ref={dropdownRef}>
                <button
                  onClick={() => setShowDropdown(!showDropdown)}
                  className="flex items-center space-x-2 text-gray-700 hover:text-blue-600 transition-colors"
                >
                  <div className="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-medium">
                    {user?.username?.charAt(0).toUpperCase() || 'U'}
                  </div>

                  <span className="hidden md:inline">{user?.username}</span>
                  
                  <svg 
                    className={`w-4 h-4 transition-transform ${showDropdown ? 'rotate-180' : ''}`}
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>

                {showDropdown && (
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50">
                    <Link 
                      href="/profile" 
                      className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                      onClick={() => setShowDropdown(false)}
                    >
                      View Profile
                    </Link>
                    <Link 
                      href="/profile/edit" 
                      className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                      onClick={() => setShowDropdown(false)}
                    >
                      Edit Profile
                    </Link>
                    <Link 
                      href="/profile/change-password" 
                      className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                      onClick={() => setShowDropdown(false)}
                    >
                      Change Password
                    </Link>
                    
                    <hr className="my-2" />

                    <button
                      onClick={handleLogout}
                      className="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100 transition-colors"
                    >
                      Logout
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        <div className="md:hidden pb-4 border-t border-gray-200 mt-4 pt-4">
          <div className="flex flex-col space-y-2">
            <Link href="/" className="text-gray-700 hover:text-blue-600 py-1">
              Recipes
            </Link>
            <Link href="/ingredients" className="text-gray-700 hover:text-blue-600 py-1">
              Ingredients  
            </Link>
            <Link href="/find-recipes" className="text-gray-700 hover:text-blue-600 py-1">
              Find by Ingredients
            </Link>
            {isAuthenticated && (
              <Link href="/liked-recipes" className="text-gray-700 hover:text-blue-600 py-1">
                My Liked Recipes
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
