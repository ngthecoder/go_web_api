'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

interface LikeButtonProps {
  recipeId: number;
  recipeName: string;
  initialLiked?: boolean;
  size?: 'small' | 'medium' | 'large';
  showLabel?: boolean;
  onLikeChange?: (isLiked: boolean) => void;
}

export default function LikeButton({ 
  recipeId, 
  recipeName,
  initialLiked = false,
  size = 'medium',
  showLabel = false,
  onLikeChange
}: LikeButtonProps) {
  const { isAuthenticated, token } = useAuth();
  const router = useRouter();
  const [isLiked, setIsLiked] = useState(initialLiked);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setIsLiked(initialLiked);
  }, [initialLiked]);

  const sizeClasses = {
    small: 'w-6 h-6',
    medium: 'w-8 h-8',
    large: 'w-10 h-10'
  };

  const handleClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      localStorage.setItem('redirectAfterLogin', window.location.pathname);
      router.push('/login?prompt=like-recipe');
      return;
    }

    const newLikedState = !isLiked;
    setIsLiked(newLikedState);
    setIsLoading(true);
    setError(null);

    try {
      if (!isLiked) {
        const response = await fetch('http://localhost:8000/api/user/liked-recipes/add', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({ recipe_id: recipeId })
        });

        if (!response.ok) {
          throw new Error('Failed to like recipe');
        }
      } else {
        const response = await fetch(`http://localhost:8000/api/user/liked-recipes/${recipeId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (!response.ok) {
          throw new Error('Failed to unlike recipe');
        }
      }

      if (onLikeChange) {
        onLikeChange(newLikedState);
      }
    } catch (err) {
      setIsLiked(isLiked);
      setError(err instanceof Error ? err.message : 'Something went wrong');
      console.error('Error toggling like:', err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex items-center gap-2">
      <button
        onClick={handleClick}
        disabled={isLoading}
        className={`${sizeClasses[size]} flex items-center justify-center rounded-full transition-all ${
          isLiked 
            ? 'text-red-500 hover:text-red-600' 
            : 'text-gray-400 hover:text-red-500'
        } ${isLoading ? 'opacity-50 cursor-not-allowed' : 'hover:scale-110'}`}
        title={isAuthenticated 
          ? (isLiked ? `Unlike ${recipeName}` : `Like ${recipeName}`)
          : 'Login to like recipes'
        }
      >
        {isLoading ? (
          <svg className="animate-spin" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        ) : (
          <svg 
            fill={isLiked ? 'currentColor' : 'none'} 
            viewBox="0 0 24 24" 
            stroke="currentColor" 
            strokeWidth={2}
          >
            <path 
              strokeLinecap="round" 
              strokeLinejoin="round" 
              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" 
            />
          </svg>
        )}
      </button>
      
      {showLabel && (
        <span className="text-sm text-gray-600">
          {isLiked ? 'Liked' : 'Like'}
        </span>
      )}

      {error && (
        <span className="text-xs text-red-600">{error}</span>
      )}
    </div>
  );
}
