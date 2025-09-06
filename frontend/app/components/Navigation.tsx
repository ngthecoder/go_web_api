import Link from 'next/link';

export default function Navigation() {
  return (
    <nav className="bg-white shadow-sm border-b mb-6">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <Link href="/" className="text-xl font-bold">
            Recipe API
          </Link>
          
          <div className="flex space-x-6">
            <Link href="/" className="text-gray-700 hover:text-blue-600">
              レシピ一覧
            </Link>
            <Link href="/find-recipes" className="text-gray-700 hover:text-blue-600">
              食材検索
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
}