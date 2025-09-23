'use client';

import { useState, useEffect } from 'react';
import { customerAPI } from '@/lib/api';

interface CustomerCreationProps {
  onCustomerCreated: (customerId: string) => void;
}

export default function CustomerCreation({ onCustomerCreated }: CustomerCreationProps) {
  const [customerId, setCustomerId] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [checkingExisting, setCheckingExisting] = useState(false);

  useEffect(() => {
    const storedCustomerId = localStorage.getItem('stripeCustomerId');
    if (storedCustomerId) {
      setCustomerId(storedCustomerId);
      onCustomerCreated(storedCustomerId);
    } else {
      // Check for existing Stripe customer when component mounts
      checkExistingCustomer();
    }
  }, []);

  const checkExistingCustomer = async () => {
    setCheckingExisting(true);
    try {
      const data = await customerAPI.getExisting();
      if (data.exists && data.stripe_customer_id) {
        localStorage.setItem('stripeCustomerId', data.stripe_customer_id);
        setCustomerId(data.stripe_customer_id);
        onCustomerCreated(data.stripe_customer_id);
      }
    } catch (err) {
      // No existing customer or error checking, that's fine
      console.log('No existing Stripe customer found');
    } finally {
      setCheckingExisting(false);
    }
  };

  const createCustomer = async () => {
    setLoading(true);
    setError(null);
    try {
      // First check if customer already exists in DB
      const existingData = await customerAPI.getExisting();
      if (existingData.exists && existingData.stripe_customer_id) {
        // Use existing customer
        localStorage.setItem('stripeCustomerId', existingData.stripe_customer_id);
        setCustomerId(existingData.stripe_customer_id);
        onCustomerCreated(existingData.stripe_customer_id);
      } else {
        // Create new customer
        const data = await customerAPI.create();
        const newCustomerId = data.customer_id;
        localStorage.setItem('stripeCustomerId', newCustomerId);
        setCustomerId(newCustomerId);
        onCustomerCreated(newCustomerId);
      }
    } catch (err) {
      setError('Failed to create/fetch customer');
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
      <h2 className="text-2xl font-bold mb-4">Step 2: Customer Creation</h2>
      
      {checkingExisting ? (
        <p className="text-gray-600">Checking for existing customer...</p>
      ) : !customerId ? (
        <div>
          <p className="text-gray-600 mb-4">
            Fetch or create a Stripe customer to enable payment operations
          </p>
          <button
            onClick={createCustomer}
            disabled={loading}
            className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600 disabled:opacity-50"
          >
            {loading ? 'Processing...' : 'Fetch/Create Customer'}
          </button>
        </div>
      ) : (
        <div>
          <p className="text-green-600 mb-2">✓ Customer ready</p>
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