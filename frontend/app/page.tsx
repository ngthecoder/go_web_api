'use client';

import { useState } from 'react';

export default function HomePage() {
  const [name, setName] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  
  const callAPI = async () => {
    setLoading(true);
    try {
      const url = `http://localhost:8000/api/hello?name=${name}`;
      const response = await fetch(url);
      const data = await response.json();
      
      setMessage(data.message);
    } catch (error) {
      setMessage('APIの呼び出しに失敗');
    }
    setLoading(false);
  };

  return (
    <div className="min-h-screen bg-blue-500 text-white">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto bg-white text-black rounded-lg p-6">
          <h2 className="text-2xl font-bold mb-4">APIテスト</h2>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-bold mb-2">
                名前を入力してください
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="名前"
                className="w-full px-3 py-2 border rounded focus:outline-none focus:border-blue-500"
              />
            </div>
            
            <button
              onClick={callAPI}
              disabled={loading}
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:bg-gray-400"
            >
              {loading ? 'APIを呼び出し中...' : 'APIを呼び出す'}
            </button>
          </div>

          {message && (
            <div className="mt-6 p-4 bg-gray-100 rounded">
              <strong>APIからのレスポンス：</strong> {message}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}