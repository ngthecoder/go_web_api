'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import { API_ENDPOINTS } from '@/lib/api-config';

interface ShoppingItem {
  ingredient_id: number;
  name: string;
  quantity: number;
  unit: string;
  notes: string;
}

interface ShoppingListData {
  recipe_id: number;
  shopping_list: ShoppingItem[];
}

export default function ShoppingListPage() {
  const params = useParams();
  const [shoppingList, setShoppingList] = useState<ShoppingListData | null>(null);
  const [recipeName, setRecipeName] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [ownedIngredients, setOwnedIngredients] = useState<number[]>([]);
  const [ownedIngredientsDetails, setOwnedIngredientsDetails] = useState<ShoppingItem[]>([]);
  const [isUpdating, setIsUpdating] = useState(false);

  useEffect(() => {
    const fetchRecipeName = async () => {
      try {
        const response = await fetch(API_ENDPOINTS.recipeById(String(params.id)));
        if (!response.ok) throw new Error('Recipe not found');
        const data = await response.json();
        setRecipeName(data.recipe.name);
      } catch (err) {
        console.error('Error fetching recipe:', err);
      }
    };

    if (params.id) {
      fetchRecipeName();
      fetchShoppingList();
    }
  }, [params.id]);

  const fetchShoppingList = async () => {
    setIsUpdating(true);
    try {
      const ingredientIds = ownedIngredients.join(',');
      const url = ownedIngredients.length > 0
        ? `${API_ENDPOINTS.shoppingList(String(params.id))}?have_ingredients=${ingredientIds}`
        : API_ENDPOINTS.shoppingList(String(params.id));
      
      const response = await fetch(url);
      
      if (!response.ok) {
        throw new Error('Failed to fetch shopping list');
      }
      
      const data = await response.json();
      setShoppingList(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong');
      console.error('Error fetching shopping list:', err);
    } finally {
      setLoading(false);
      setIsUpdating(false);
    }
  };

  const toggleIngredient = (item: ShoppingItem) => {
    const ingredientId = item.ingredient_id;
    
    if (ownedIngredients.includes(ingredientId)) {
      setOwnedIngredients(prev => prev.filter(id => id !== ingredientId));
      setOwnedIngredientsDetails(prev => 
        prev.filter(detail => detail.ingredient_id !== ingredientId)
      );
    } else {
      setOwnedIngredients(prev => [...prev, ingredientId]);
      setOwnedIngredientsDetails(prev => [...prev, item]);
    }
  };

  const handleUpdateList = () => {
    fetchShoppingList();
  };

  const handlePrint = () => {
    window.print();
  };

  if (loading) {
    return (
      <div className="container mx-auto p-4 max-w-3xl">
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto p-4 max-w-3xl">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <h2 className="text-xl font-semibold text-red-800 mb-2">Error</h2>
          <p className="text-red-600">{error}</p>
          <Link
            href={`/recipes/${params.id}`}
            className="mt-4 inline-block px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
          >
            Back to Recipe
          </Link>
        </div>
      </div>
    );
  }

  if (!shoppingList) {
    return (
      <div className="container mx-auto p-4 max-w-3xl">
        <div className="text-center py-12">
          <p className="text-gray-600">Shopping list not found</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4 max-w-3xl">
      <nav className="mb-6 text-sm text-gray-600">
        <Link href="/" className="hover:text-blue-600">Home</Link>
        <span className="mx-2">→</span>
        <Link href={`/recipes/${params.id}`} className="hover:text-blue-600">
          {recipeName || 'Recipe'}
        </Link>
        <span className="mx-2">→</span>
        <span className="text-gray-900">Shopping List</span>
      </nav>

      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Shopping List</h1>
        {recipeName && (
          <p className="text-gray-600">For: <span className="font-medium">{recipeName}</span></p>
        )}
      </div>

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
        <div className="flex items-start">
          <svg className="w-5 h-5 text-blue-600 mr-2 mt-0.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <p className="text-sm text-gray-700">
              <strong>How to use:</strong> Check off the ingredients you already have at home, 
              then click "Update List" to see only what you need to buy.
            </p>
          </div>
        </div>
      </div>

      {shoppingList?.shopping_list === null || shoppingList?.shopping_list?.length === 0 ? (
        <div className="bg-green-50 border border-green-200 rounded-lg p-8 text-center">
          <svg className="w-16 h-16 text-green-500 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h2 className="text-2xl font-semibold text-green-800 mb-2">
            You're All Set! 🎉
          </h2>
          <p className="text-green-700">
            You have all the ingredients needed for this recipe.
          </p>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-sm border">
          <div className="p-6 border-b">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold">
                Items to Buy ({shoppingList?.shopping_list?.length})
              </h2>
              <button
                onClick={handlePrint}
                className="text-sm text-gray-600 hover:text-gray-900 flex items-center gap-1"
              >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                </svg>
                Print
              </button>
            </div>
          </div>

          <div className="divide-y">
            {shoppingList?.shopping_list?.map((item, index) => (
              <label
                key={item.ingredient_id}
                className="flex items-start p-4 hover:bg-gray-50 cursor-pointer transition-colors"
              >
                <input
                  type="checkbox"
                  checked={ownedIngredients.includes(item.ingredient_id)}
                  onChange={() => toggleIngredient(item)}
                  className="mt-1 mr-4 h-5 w-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <div className="flex-1">
                  <div className="flex items-baseline gap-2">
                    <span className="font-medium text-gray-900">{item.name}</span>
                    <span className="text-sm text-gray-600">
                      {item.quantity} {item.unit}
                    </span>
                  </div>
                  {item.notes && (
                    <p className="text-sm text-gray-500 mt-1">
                      {item.notes}
                    </p>
                  )}
                </div>
                <Link
                  href={`/ingredients/${item.ingredient_id}`}
                  className="ml-2 text-sm text-blue-600 hover:text-blue-800"
                  onClick={(e) => e.stopPropagation()}
                >
                  View
                </Link>
              </label>
            ))}
          </div>
        </div>
      )}

      {ownedIngredientsDetails.length > 0 && (
        <div className="mt-6 bg-gray-50 border border-gray-200 rounded-lg p-4">
          <h3 className="text-sm font-medium text-gray-700 mb-2">
            Ingredients you have ({ownedIngredientsDetails.length}):
          </h3>
          <div className="flex flex-wrap gap-2">
            {ownedIngredientsDetails?.map(item => (
              <span
                key={item.ingredient_id}
                className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800"
              >
                {item.name}
                <button
                  onClick={() => toggleIngredient(item)}
                  className="ml-2 hover:text-blue-600"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        </div>
      )}

      <div className="mt-6 flex flex-wrap gap-4">
        <button
          onClick={handleUpdateList}
          disabled={isUpdating}
          className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
        >
          {isUpdating ? (
            <>
              <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Updating...
            </>
          ) : (
            <>
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Update List
            </>
          )}
        </button>

        <Link
          href={`/recipes/${params.id}`}
          className="px-6 py-3 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition-colors"
        >
          Back to Recipe
        </Link>

        <button
          onClick={handlePrint}
          className="px-6 py-3 bg-white border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors flex items-center gap-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
          </svg>
          Print List
        </button>
      </div>

      <style jsx global>{`
        @media print {
          body * {
            visibility: hidden;
          }
          .container, .container * {
            visibility: visible;
          }
          .container {
            position: absolute;
            left: 0;
            top: 0;
            width: 100%;
          }
          button, nav, .bg-blue-50 {
            display: none !important;
          }
        }
      `}</style>
    </div>
  );
}
