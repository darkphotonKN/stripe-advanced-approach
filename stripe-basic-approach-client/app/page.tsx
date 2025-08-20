'use client';

import { useState } from 'react';
import CustomerLifecycle from '@/components/CustomerLifecycle';
import PaymentMethodManager from '@/components/PaymentMethodManager';
import PaymentProcessor from '@/components/PaymentProcessor';

export default function Home() {
  const [customerId, setCustomerId] = useState<string | null>(null);
  const [activePhase, setActivePhase] = useState<number>(1);

  const handleCustomerCreated = (id: string) => {
    setCustomerId(id);
  };

  const handlePaymentMethodSaved = () => {
    console.log('Payment method saved successfully');
  };

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">
            Stripe Payment Flow Demo
          </h1>
          <p className="text-gray-600">
            Practice implementing the backend to match this frontend flow
          </p>
        </div>

        <div className="flex justify-center mb-8">
          <div className="flex space-x-4">
            <button
              onClick={() => setActivePhase(1)}
              className={`px-4 py-2 rounded ${
                activePhase === 1
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              Phase 1: Customer
            </button>
            <button
              onClick={() => setActivePhase(2)}
              className={`px-4 py-2 rounded ${
                activePhase === 2
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              } ${!customerId ? 'opacity-50 cursor-not-allowed' : ''}`}
              disabled={!customerId}
            >
              Phase 2: Payment Method
            </button>
            <button
              onClick={() => setActivePhase(3)}
              className={`px-4 py-2 rounded ${
                activePhase === 3
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              } ${!customerId ? 'opacity-50 cursor-not-allowed' : ''}`}
              disabled={!customerId}
            >
              Phase 3: Payment
            </button>
          </div>
        </div>

        <div className="space-y-6">
          {activePhase === 1 && (
            <CustomerLifecycle onCustomerCreated={handleCustomerCreated} />
          )}
          
          {activePhase === 2 && (
            <PaymentMethodManager 
              customerId={customerId}
              onPaymentMethodSaved={handlePaymentMethodSaved}
            />
          )}
          
          {activePhase === 3 && (
            <PaymentProcessor customerId={customerId} />
          )}
        </div>

        <div className="mt-8 p-4 bg-blue-50 rounded-lg border border-blue-200">
          <h3 className="font-semibold text-blue-900 mb-2">Backend Implementation Guide:</h3>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>• POST /api/customers - Create Stripe customer</li>
            <li>• POST /api/payment-methods/setup-intent - Create setup intent for saving cards</li>
            <li>• GET /api/payment-methods/:customerId - List saved payment methods</li>
            <li>• DELETE /api/payment-methods/:id - Detach payment method</li>
            <li>• POST /api/payments/create-intent - Create payment intent</li>
            <li>• POST /api/payments/confirm - Confirm payment (if needed)</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
