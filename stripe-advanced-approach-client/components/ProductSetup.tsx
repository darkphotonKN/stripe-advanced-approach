'use client';

import { useState } from 'react';
import { productAPI } from '@/lib/api';

interface ProductSetupProps {
  onProductsCreated: (priceId: string) => void;
}

type ProductType = 'one-time' | 'recurring';

export default function ProductSetup({ onProductsCreated }: ProductSetupProps) {
  const [loading, setLoading] = useState(false);
  const [productsCreated, setProductsCreated] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [productType, setProductType] = useState<ProductType>('one-time');
  const [productName, setProductName] = useState('Example Product');
  const [productDescription, setProductDescription] = useState('New Product');
  const [productPrice, setProductPrice] = useState(10);

  const setupProducts = async () => {
    setLoading(true);
    setError(null);
    
    try {
      let data;
      if (productType === 'recurring') {
        data = await productAPI.setupSubscription(
          productName,
          productDescription,
          productPrice * 100 // Convert to cents
        );
      } else {
        data = await productAPI.setupProducts(
          productName,
          productDescription,
          productPrice * 100 // Convert to cents
        );
      }
      
      localStorage.setItem('subscriptionPriceId', data.price_id || data.subscription_price_id);
      setProductsCreated(true);
      onProductsCreated(data.price_id || data.subscription_price_id);
    } catch (err) {
      setError('Failed to setup products');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Create Products & Subscriptions</h2>
      
      {!productsCreated ? (
        <div>
          <p className="text-gray-600 mb-4">
            Create sample products and prices in Stripe catalog
          </p>
          
          <div className="space-y-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Product Type
              </label>
              <select
                value={productType}
                onChange={(e) => setProductType(e.target.value as ProductType)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
              >
                <option value="one-time">One-Time Purchase</option>
                <option value="recurring">Recurring Subscription</option>
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Product Name
              </label>
              <input
                type="text"
                value={productName}
                onChange={(e) => setProductName(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black placeholder-gray-700"
                placeholder="Enter product name"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Product Description
              </label>
              <input
                type="text"
                value={productDescription}
                onChange={(e) => setProductDescription(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black placeholder-gray-700"
                placeholder="Enter product description"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {productType === 'recurring' ? 'Monthly Price ($)' : 'Price ($)'}
              </label>
              <input
                type="number"
                min="1"
                step="0.01"
                value={productPrice}
                onChange={(e) => setProductPrice(parseFloat(e.target.value))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black placeholder-gray-700"
                placeholder={productType === 'recurring' ? 'Enter monthly price' : 'Enter price'}
              />
            </div>
          </div>
          
          <button
            onClick={setupProducts}
            disabled={loading || !productName || !productDescription || productPrice <= 0}
            className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600 disabled:opacity-50"
          >
            {loading ? 'Creating Product...' : `Create ${productType === 'recurring' ? 'Subscription' : 'One-Time'} Product`}
          </button>
        </div>
      ) : (
        <div>
          <p className="text-green-600 mb-2">✓ Product created successfully</p>
          <p className="text-sm text-gray-600">
            {productType === 'recurring' ? 'Subscription' : 'One-time'} product price ID stored for later use
          </p>
        </div>
      )}
      
      {error && (
        <p className="text-red-500 mt-4">{error}</p>
      )}
    </div>
  );
}