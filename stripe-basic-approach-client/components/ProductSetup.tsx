'use client';

import { useState } from 'react';
import { productAPI } from '@/lib/api';

interface ProductSetupProps {
  onProductsCreated: (priceId: string) => void;
}

export default function ProductSetup({ onProductsCreated }: ProductSetupProps) {
  const [loading, setLoading] = useState(false);
  const [productsCreated, setProductsCreated] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const setupProducts = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const data = await productAPI.setupProducts();
      localStorage.setItem('subscriptionPriceId', data.subscription_price_id);
      setProductsCreated(true);
      onProductsCreated(data.subscription_price_id);
    } catch (err) {
      setError('Failed to setup products');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Phase 1: Product Setup</h2>
      
      {!productsCreated ? (
        <div>
          <p className="text-gray-600 mb-4">
            Create sample products and prices in Stripe catalog
          </p>
          <button
            onClick={setupProducts}
            disabled={loading}
            className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600 disabled:opacity-50"
          >
            {loading ? 'Creating Products...' : 'Create Sample Products'}
          </button>
        </div>
      ) : (
        <div>
          <p className="text-green-600 mb-2">✓ Products created successfully</p>
          <p className="text-sm text-gray-600">Price IDs stored for later use</p>
        </div>
      )}
      
      {error && (
        <p className="text-red-500 mt-4">{error}</p>
      )}
    </div>
  );
}