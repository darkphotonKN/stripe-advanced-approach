"use client";

import { useState, useEffect } from "react";
import AuthWrapper from "@/components/AuthWrapper";
import ProductSetup from "@/components/ProductSetup";
import CustomerCreation from "@/components/CustomerCreation";
import CardSaving from "@/components/CardSaving";
import OneTimePayment from "@/components/OneTimePayment";
import SubscriptionCreation from "@/components/SubscriptionCreation";

export default function Home() {
  const [stripeCustomerId, setStripeCustomerId] = useState<string | null>(null);
  const [subscriptionPriceId, setSubscriptionPriceId] = useState<string | null>(
    null,
  );
  const [productsCreated, setProductsCreated] = useState(false);
  const [customerCreated, setCustomerCreated] = useState(false);

  useEffect(() => {
    const storedCustomerId = localStorage.getItem("stripeCustomerId");
    const storedPriceId = localStorage.getItem("subscriptionPriceId");

    if (storedCustomerId) {
      setStripeCustomerId(storedCustomerId);
      setCustomerCreated(true);
    }
    if (storedPriceId) {
      setSubscriptionPriceId(storedPriceId);
      setProductsCreated(true);
    }
  }, []);

  const handleProductsCreated = (priceId: string) => {
    setSubscriptionPriceId(priceId);
    setProductsCreated(true);
  };

  const handleCustomerCreated = (customerId: string) => {
    setStripeCustomerId(customerId);
    setCustomerCreated(true);
  };

  const resetAll = () => {
    localStorage.clear();
    setStripeCustomerId(null);
    setSubscriptionPriceId(null);
    setProductsCreated(false);
    setCustomerCreated(false);
    window.location.reload();
  };

  return (
    <AuthWrapper>
      <div className="min-h-screen bg-gray-50 p-4">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-gray-800 mb-2">
              Stripe Payment Flow POC
            </h1>
            <p className="text-gray-600 mb-4">
              Following the exact POC flow: Products → Customer → Save Card →
              Payment → Subscription
            </p>
            <button
              onClick={resetAll}
              className="text-sm text-red-500 hover:text-red-700"
            >
              Reset Stripe Data
            </button>
          </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <ProductSetup onProductsCreated={handleProductsCreated} />

          <CustomerCreation
            onCustomerCreated={handleCustomerCreated}
            enabled={productsCreated}
          />

          <CardSaving customerId={stripeCustomerId} enabled={customerCreated} />

          <OneTimePayment
            customerId={stripeCustomerId}
            enabled={customerCreated}
          />

          <SubscriptionCreation
            customerId={stripeCustomerId}
            priceId={subscriptionPriceId}
            enabled={customerCreated && productsCreated}
          />
        </div>

        <div className="mt-8 p-6 bg-blue-50 rounded-lg border border-blue-200">
          <h3 className="font-bold text-blue-900 mb-4">
            Backend Implementation Requirements:
          </h3>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            <div>
              <h4 className="font-semibold text-blue-800 mb-2">
                Phase 1: Product Setup
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/setup-products</li>
                <li>• Creates Stripe products/prices</li>
                <li>• Returns subscription_price_id</li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-blue-800 mb-2">
                Phase 2: Customer Creation
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/create-customer</li>
                <li>• Creates Stripe customer</li>
                <li>• Returns customer_id</li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-blue-800 mb-2">
                Phase 3: Card Saving
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/save-card</li>
                <li>• Creates setup intent</li>
                <li>• Returns client_secret</li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-blue-800 mb-2">
                Phase 4: One-Time Payment
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/create-payment-intent</li>
                <li>• Amount: 2000 ($20.00)</li>
                <li>• Returns client_secret</li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-blue-800 mb-2">
                Phase 5: Subscription
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/create-subscription</li>
                <li>• Creates subscription + payment intent</li>
                <li>• Returns subscription_id + client_secret</li>
              </ul>
            </div>
          </div>

          <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded">
            <h4 className="font-semibold text-green-800 mb-2">Authentication Endpoints:</h4>
            <ul className="text-green-700 text-sm space-y-1">
              <li>• POST /api/signup - Create account (email, password, name)</li>
              <li>• POST /api/signin - Login (email, password)</li>
              <li>• Returns JWT token for authorization</li>
            </ul>
          </div>

          <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded">
            <p className="text-sm text-yellow-800">
              <strong>Critical Flow:</strong> Products must be created first,
              then customer, then all payment operations use the stored customer
              ID. Frontend confirms all payments.
            </p>
          </div>
        </div>

        <div className="mt-4 p-4 bg-gray-100 rounded text-sm text-gray-600">
          <p className="font-semibold mb-1">Test Card Numbers:</p>
          <ul className="space-y-1">
            <li>✓ Success: 4242 4242 4242 4242</li>
            <li>✗ Decline: 4000 0000 0000 0002</li>
            <li>Use any future date and any 3-digit CVC</li>
          </ul>
        </div>
        </div>
      </div>
    </AuthWrapper>
  );
}
