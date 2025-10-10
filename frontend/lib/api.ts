import { API_ENDPOINTS } from "./api-config";

export async function updateProfile(token: string, username: string, email: string) {
  const response = await fetch(API_ENDPOINTS.updateProfile, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ username, email })
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return response.json();
}

export async function changePassword(token: string, currentPassword: string, newPassword: string) {
  const response = await fetch(API_ENDPOINTS.updatePassword, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ 
      current_password: currentPassword, 
      new_password: newPassword 
    })
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return response.json();
}

export async function deleteAccount(token: string, password: string) {
  const response = await fetch(API_ENDPOINTS.deleteAccount, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ password })
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }

  return response.json();
}

export async function toggleLikedRecipe(token: string | null, recipeId: number, isLiked: boolean) {
  if (!isLiked) {
      const response = await fetch(API_ENDPOINTS.addLikedRecipe, {
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
      const response = await fetch(API_ENDPOINTS.removeLikedRecipe(recipeId), {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error('Failed to unlike recipe');
      }
    }
}