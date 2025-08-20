'use client';

import { useState, useEffect } from 'react';
import { customerAPI } from '@/lib/api';

interface CustomerCreationProps {
  onCustomerCreated: (customerId: string) => void;
  enabled: boolean;
}

export default function CustomerCreation({ onCustomerCreated, enabled }: CustomerCreationProps) {
  const [customerId, setCustomerId] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const storedCustomerId = localStorage.getItem('stripeCustomerId');
    if (storedCustomerId) {
      setCustomerId(storedCustomerId);
      onCustomerCreated(storedCustomerId);
    }
  }, [onCustomerCreated]);

  const createCustomer = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await customerAPI.create();
      const newCustomerId = data.customer_id;
      localStorage.setItem('stripeCustomerId', newCustomerId);
      setCustomerId(newCustomerId);
      onCustomerCreated(newCustomerId);
    } catch (err) {
      setError('Failed to create customer');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const resetCustomer = () => {
    localStorage.removeItem('stripeCustomerId');
    setCustomerId(null);
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Phase 2: Customer Creation</h2>
      
      {!enabled ? (
        <p className="text-gray-500">Setup products first to enable customer creation</p>
      ) : !customerId ? (
        <div>
          <p className="text-gray-600 mb-4">
            Create a Stripe customer to enable payment operations
          </p>
          <button
            onClick={createCustomer}
            disabled={loading}
            className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600 disabled:opacity-50"
          >
            {loading ? 'Creating...' : 'Fetch Customer ID'}
          </button>
        </div>
      ) : (
        <div>
          <p className="text-green-600 mb-2">✓ Customer created successfully</p>
          <p className="text-sm text-gray-600 mb-4">Customer ID: {customerId}</p>
          <button
            onClick={resetCustomer}
            className="text-red-500 hover:text-red-700 text-sm"
          >
            Reset Customer (for testing)
          </button>
        </div>
      )}
      
      {error && (
        <p className="text-red-500 mt-4">{error}</p>
      )}
    </div>
  );
}