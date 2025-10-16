'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { API_ENDPOINTS } from '@/lib/api-config';
import StatCard from '@/components/StatCard';

interface CategoryStats {
  ingredient_categories: Record<string, number>;
  recipe_categories: Record<string, number>;
}

interface Stats {
  total_ingredients: number;
  total_recipes: number;
  avg_prep_time: number;
  avg_cook_time: number;
  difficulty_distribution: Record<string, number>;
}

export default function StatsPage() {
  const [categoryStats, setCategoryStats] = useState<CategoryStats | null>(null);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      setLoading(true);

      const [categoriesResponse, statsResponse] = await Promise.all([
        fetch(API_ENDPOINTS.categories),
        fetch(API_ENDPOINTS.stats)
      ]);

      if (!categoriesResponse.ok || !statsResponse.ok) {
        throw new Error('Failed to fetch statistics');
      }

      const categoriesData = await categoriesResponse.json();
      const statsData = await statsResponse.json();

      setCategoryStats(categoriesData);
      setStats(statsData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load statistics');
      console.error('Error fetching stats:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto p-4 max-w-7xl">
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto p-4 max-w-7xl">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <h2 className="text-xl font-semibold text-red-800 mb-2">Error</h2>
          <p className="text-red-600">{error}</p>
          <button
            onClick={fetchStats}
            className="mt-4 px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  if (!categoryStats || !stats) {
    return (
      <div className="container mx-auto p-4 max-w-7xl">
        <div className="text-center py-12">
          <p className="text-gray-600">No statistics available</p>
        </div>
      </div>
    );
  }

  const totalRecipes = Object.values(stats.difficulty_distribution).reduce((a, b) => a + b, 0);
  const difficultyPercentages = Object.entries(stats.difficulty_distribution).map(([key, value]) => ({
    name: key,
    count: value,
    percentage: totalRecipes > 0 ? ((value / totalRecipes) * 100).toFixed(1) : 0
  }));

  const maxIngredientCount = Math.max(...Object.values(categoryStats.ingredient_categories));
  const maxRecipeCount = Math.max(...Object.values(categoryStats.recipe_categories));

  return (
    <div className="container mx-auto p-4">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Statistics Dashboard</h1>
        <p className="text-gray-600">Overview of recipes and ingredients in the system</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Total Recipes"
          value={stats.total_recipes}
          icon="📖"
          color="bg-blue-500"
        />
        <StatCard
          title="Total Ingredients"
          value={stats.total_ingredients}
          icon="🥗"
          color="bg-green-500"
        />
        <StatCard
          title="Avg Prep Time"
          value={`${stats.avg_prep_time.toFixed(1)} min`}
          icon="⏱️"
          color="bg-yellow-500"
        />
        <StatCard
          title="Avg Cook Time"
          value={`${stats.avg_cook_time.toFixed(1)} min`}
          icon="🍳"
          color="bg-orange-500"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow-sm border p-6">
          <h2 className="text-xl font-semibold mb-4">Recipes by Category</h2>
          <div className="space-y-4">
            {Object.entries(categoryStats.recipe_categories)
              .sort(([, a], [, b]) => b - a)
              .map(([category, count]) => (
                <div key={category}>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-medium text-gray-700 capitalize">
                      {category}
                    </span>
                    <span className="text-sm text-gray-600">{count}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2.5">
                    <div
                      className="bg-blue-600 h-2.5 rounded-full transition-all"
                      style={{ width: `${(count / maxRecipeCount) * 100}%` }}
                    ></div>
                  </div>
                </div>
              ))}
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border p-6">
          <h2 className="text-xl font-semibold mb-4">Ingredients by Category</h2>
          <div className="space-y-4">
            {Object.entries(categoryStats.ingredient_categories)
              .sort(([, a], [, b]) => b - a)
              .map(([category, count]) => (
                <div key={category}>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-medium text-gray-700 capitalize">
                      {category}
                    </span>
                    <span className="text-sm text-gray-600">{count}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2.5">
                    <div
                      className="bg-green-600 h-2.5 rounded-full transition-all"
                      style={{ width: `${(count / maxIngredientCount) * 100}%` }}
                    ></div>
                  </div>
                </div>
              ))}
          </div>
        </div>
      </div>
      
      <div className="bg-white rounded-lg shadow-sm border p-6 mb-8">
        <h2 className="text-xl font-semibold mb-6">Recipe Difficulty Distribution</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {difficultyPercentages?.map((item) => (
            <div
              key={item.name}
              className="relative overflow-hidden rounded-lg border-2 border-gray-200 p-6 hover:border-blue-400 transition-colors"
            >
              <div className="text-center">
                <div className="text-4xl mb-2">
                  {item.name === 'easy' && '😊'}
                  {item.name === 'medium' && '🤔'}
                  {item.name === 'hard' && '😰'}
                </div>
                <h3 className="text-lg font-semibold capitalize mb-2">{item.name}</h3>
                <div className="text-3xl font-bold text-blue-600 mb-1">{item.count}</div>
                <div className="text-sm text-gray-500">{item.percentage}% of recipes</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg border border-blue-200 p-6">
        <h2 className="text-xl font-semibold mb-4">Explore More</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Link
            href="/"
            className="flex items-center gap-3 p-4 bg-white rounded-lg hover:shadow-md transition-shadow"
          >
            <span className="text-2xl">📖</span>
            <div>
              <div className="font-medium">Browse Recipes</div>
              <div className="text-sm text-gray-600">Explore all recipes</div>
            </div>
          </Link>
          <Link
            href="/ingredients"
            className="flex items-center gap-3 p-4 bg-white rounded-lg hover:shadow-md transition-shadow"
          >
            <span className="text-2xl">🥗</span>
            <div>
              <div className="font-medium">View Ingredients</div>
              <div className="text-sm text-gray-600">See all ingredients</div>
            </div>
          </Link>
          <Link
            href="/find-recipes"
            className="flex items-center gap-3 p-4 bg-white rounded-lg hover:shadow-md transition-shadow"
          >
            <span className="text-2xl">🔍</span>
            <div>
              <div className="font-medium">Find Recipes</div>
              <div className="text-sm text-gray-600">Search by ingredients</div>
            </div>
          </Link>
        </div>
      </div>
    </div>
  );
}
