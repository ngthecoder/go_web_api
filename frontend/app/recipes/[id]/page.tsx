'use client';
import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';

interface IngredientWithQuantity {
  ingredient_id: number;
  name: string;
  quantity: number;
  unit: string;
  notes: string;
}

interface Recipe {
  id: number;
  name: string;
  category: string;
  prep_time_minutes: number;
  cook_time_minutes: number;
  servings: number;
  difficulty: string;
  instructions: string;
  description: string;
}

interface RecipeDetail {
  recipe: Recipe;
  ingredients: IngredientWithQuantity[];
}

export default function RecipeDetailPage() {
  const params = useParams();
  const [recipeDetail, setRecipeDetail] = useState<RecipeDetail | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (params.id) {
      fetch(`http://localhost:8000/api/recipes/${params.id}`)
        .then(res => res.json())
        .then(data => {
          setRecipeDetail(data);
          setLoading(false);
        })
        .catch(err => {
          console.error('API Error:', err);
          setLoading(false);
        });
    }
  }, [params.id]);

  if (loading) return <div className="container mx-auto p-4">Loading...</div>;
  if (!recipeDetail) return <div className="container mx-auto p-4">Recipe not found</div>;

  const { recipe, ingredients } = recipeDetail;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-4">{recipe.name}</h1>
      <p className="text-gray-600 mb-6">{recipe.description}</p>
      
      <div className="grid md:grid-cols-2 gap-8">
        <div>
          <h2 className="text-2xl font-semibold mb-4">Recipe Information</h2>
          <ul className="space-y-2">
            <li>Category: {recipe.category}</li>
            <li>Prep Time: {recipe.prep_time_minutes} minutes</li>
            <li>Cook Time: {recipe.cook_time_minutes} minutes</li>
            <li>Servings: {recipe.servings}</li>
            <li>Difficulty: {recipe.difficulty}</li>
          </ul>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-4">Ingredients</h2>
          <ul className="space-y-2">
            {ingredients.map(ing => (
              <li key={ing.ingredient_id} className="flex justify-between">
                <span>{ing.name}</span>
                <span>{ing.quantity} {ing.unit} {ing.notes}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <div className="mt-8">
        <h2 className="text-2xl font-semibold mb-4">Instructions</h2>
        <div className="bg-gray-50 p-4 rounded-lg">
          <pre className="whitespace-pre-wrap">{recipe.instructions}</pre>
        </div>
      </div>
    </div>
  );
}