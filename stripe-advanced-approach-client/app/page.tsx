"use client";

import { useState, useEffect } from "react";
import AuthWrapper from "@/components/AuthWrapper";
import ProductSetup from "@/components/ProductSetup";
import CardSaving from "@/components/CardSaving";
import OneTimePayment from "@/components/OneTimePayment";
import BuyProduct from "@/components/BuyProduct";
import SubscribeToSite from "@/components/SubscribeToSite";
import { customerAPI } from "@/lib/api";

export default function Home() {
  const [stripeCustomerId, setStripeCustomerId] = useState<string | null>(null);
  const [subscriptionPriceId, setSubscriptionPriceId] = useState<string | null>(
    null,
  );
  const [productsCreated, setProductsCreated] = useState(false);
  const [customerCreated, setCustomerCreated] = useState(false);
  const [purchasedProducts, setPurchasedProducts] = useState<string[]>([]);
  const [subscriptionStatus, setSubscriptionStatus] = useState<any>(null);

  useEffect(() => {
    const fetchCustomerId = async () => {
      try {
        const data = await customerAPI.getExisting();
        if (data.exists && data.stripe_customer_id) {
          setStripeCustomerId(data.stripe_customer_id);
          setCustomerCreated(true);
          localStorage.setItem("stripeCustomerId", data.stripe_customer_id);
        }
      } catch (err) {
        console.log("No existing Stripe customer found");
      }
    };

    fetchCustomerId();

    const storedPriceId = localStorage.getItem("subscriptionPriceId");
    if (storedPriceId) {
      setSubscriptionPriceId(storedPriceId);
      setProductsCreated(true);
    }
  }, []);

  const handleProductsCreated = (priceId: string) => {
    setSubscriptionPriceId(priceId);
    setProductsCreated(true);
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
              Stripe Advanced Payment System
            </h1>
            <div className="mt-4 text-gray-600">
              <div className="flex items-center justify-center gap-4 mb-2">
                <span className="font-semibold text-blue-600">Admin Flow:</span>
                <span>Create Products & Subscriptions</span>
              </div>
              <div className="flex items-center justify-center gap-4">
                <span className="font-semibold text-green-600">User Flow:</span>
                <span>Save Card → Subscribe or Purchase Products</span>
              </div>
            </div>
            <button
              onClick={resetAll}
              className="mt-4 text-sm text-red-500 hover:text-red-700"
            >
              Reset Stripe Data
            </button>
          </div>

          {/* User Profile Section */}
          {stripeCustomerId && (
            <div className="mb-6 grid grid-cols-1 lg:grid-cols-3 gap-4">
              {/* Customer ID */}
              <div className="p-4 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg text-white">
                <div className="flex items-center space-x-3">
                  <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <div className="min-w-0">
                    <p className="text-sm font-medium opacity-90">Customer ID</p>
                    <p className="text-lg font-bold font-mono truncate">{stripeCustomerId}</p>
                  </div>
                </div>
              </div>

              {/* Subscription Status */}
              <div className="p-4 bg-gradient-to-r from-green-500 to-teal-600 rounded-lg text-white">
                <div className="flex items-center space-x-3">
                  <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <div>
                    <p className="text-sm font-medium opacity-90">Subscription Status</p>
                    <p className="text-lg font-bold">
                      {subscriptionStatus?.has_access ? 'Pro Active' : 'No Subscription'}
                    </p>
                  </div>
                </div>
              </div>

              {/* Purchased Products */}
              <div className="p-4 bg-gradient-to-r from-purple-500 to-pink-600 rounded-lg text-white">
                <div className="flex items-center space-x-3">
                  <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
                  </svg>
                  <div>
                    <p className="text-sm font-medium opacity-90">Products Purchased</p>
                    <p className="text-lg font-bold">{purchasedProducts.length} Items</p>
                  </div>
                </div>
              </div>
            </div>
          )}

        {/* Admin Flow Section */}
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-4 flex items-center">
            <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-lg mr-3">Admin Flow</span>
            Product & Subscription Setup
          </h2>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <ProductSetup onProductsCreated={handleProductsCreated} />
          </div>
        </div>

        {/* User Flow Section */}
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-4 flex items-center">
            <span className="bg-green-100 text-green-800 px-3 py-1 rounded-lg mr-3">User Flow</span>
            Payment & Subscription Options
          </h2>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <CardSaving customerId={stripeCustomerId} enabled={true} />

            <SubscribeToSite
              enabled={true}
              onStatusUpdate={setSubscriptionStatus}
            />

            <BuyProduct
              customerId={stripeCustomerId}
              enabled={true}
              onProductPurchased={(productName: string) => {
                setPurchasedProducts(prev => [...prev, productName]);
              }}
            />
          </div>
        </div>

        <div className="mt-8 p-6 bg-blue-50 rounded-lg border border-blue-200">
          <h3 className="font-bold text-blue-900 mb-4">
            Backend API Endpoints:
          </h3>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
            <div>
              <h4 className="font-semibold text-blue-800 mb-3 flex items-center">
                <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded mr-2">Admin</span>
                Product Management
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/payment/setup-products</li>
                <li>• POST /api/payment/setup-subscription</li>
                <li>• Creates Stripe products/prices</li>
                <li>• Returns price_id for future use</li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-blue-800 mb-3 flex items-center">
                <span className="bg-green-100 text-green-800 px-2 py-1 rounded mr-2">User</span>
                Payment Methods
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/payment/save-card</li>
                <li>• Creates setup intent</li>
                <li>• Returns client_secret for Stripe Elements</li>
                <li>• Saves card to customer for future use</li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-blue-800 mb-3 flex items-center">
                <span className="bg-green-100 text-green-800 px-2 py-1 rounded mr-2">User</span>
                Subscription Flow
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• POST /api/payment/subscription/subscribe</li>
                <li>• GET /api/payment/subscription/status</li>
                <li>• Subscribes to pre-configured pro plan</li>
                <li>• Returns subscription status & access</li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-blue-800 mb-3 flex items-center">
                <span className="bg-green-100 text-green-800 px-2 py-1 rounded mr-2">User</span>
                Product Purchases
              </h4>
              <ul className="text-blue-700 space-y-1">
                <li>• GET /api/payment/products</li>
                <li>• POST /api/payment/purchase-product</li>
                <li>• One-time purchases with saved/new cards</li>
                <li>• Returns client_secret for payment confirmation</li>
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
              <strong>Key Features:</strong> Customers auto-created during signup. Two-flow system:
              Admins create products, users choose between subscription (pro plan) or individual
              product purchases. All payments use Stripe Elements for security.
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
