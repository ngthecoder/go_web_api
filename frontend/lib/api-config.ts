export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export const API_ENDPOINTS = {
  register: `${API_URL}/api/auth/register`,
  login: `${API_URL}/api/auth/login`,

  recipes: `${API_URL}/api/recipes`,
  recipeById: (id: string | number) => `${API_URL}/api/recipes/${id}`,
  recipesByIngredients: `${API_URL}/api/recipes/find-by-ingredients`,
  shoppingList: (id: string | number) => `${API_URL}/api/recipes/shopping-list/${id}`,

  ingredients: `${API_URL}/api/ingredients`,
  ingredientById: (id: string | number) => `${API_URL}/api/ingredients/${id}`,

  userProfile: `${API_URL}/api/user/profile`,
  updateProfile: `${API_URL}/api/user/profile/update`,
  updatePassword: `${API_URL}/api/user/password`,
  deleteAccount: `${API_URL}/api/user/account`,
  likedRecipes: `${API_URL}/api/user/liked-recipes`,
  addLikedRecipe: `${API_URL}/api/user/liked-recipes/add`,
  removeLikedRecipe: (id: string | number) => `${API_URL}/api/user/liked-recipes/${id}`,

  categories: `${API_URL}/api/categories`,
  stats: `${API_URL}/api/stats`,
};