"use client";

import { useState, useEffect } from "react";
import { subscriptionAPI } from "@/lib/api";

interface SubscriptionStatus {
  has_access: boolean;
  status: string;
  cancel_at_period_end: boolean;
}

interface SubscribeToSiteProps {
  enabled: boolean;
  onStatusUpdate?: (status: SubscriptionStatus) => void;
}

export default function SubscribeToSite({ enabled, onStatusUpdate }: SubscribeToSiteProps) {
  const [activeTab, setActiveTab] = useState<"subscribe" | "content">("subscribe");
  const [loading, setLoading] = useState(false);
  const [subscriptionStatus, setSubscriptionStatus] = useState<SubscriptionStatus | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    if (enabled) {
      fetchSubscriptionStatus();
    }
  }, [enabled]);

  const fetchSubscriptionStatus = async () => {
    try {
      const status = await subscriptionAPI.getStatus();
      setSubscriptionStatus(status);
      if (onStatusUpdate) {
        onStatusUpdate(status);
      }
    } catch (err: any) {
      console.error("Failed to fetch subscription status:", err);
    }
  };

  const handleSubscribe = async () => {
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      // Call subscribe endpoint - no parameters needed
      const response = await subscriptionAPI.subscribe();

      // Handle successful subscription
      setSuccess("Successfully subscribed to Pro plan!");

      // Refresh subscription status to show updated access
      await fetchSubscriptionStatus();
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create subscription");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (!enabled) {
    return (
      <div className="bg-white p-6 rounded-lg shadow-md opacity-50">
        <h2 className="text-2xl font-bold mb-4">Step 2: Subscribe to Pro Plan</h2>
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Step 2: Subscribe to Pro Plan</h2>

      {/* Tabs */}
      <div className="flex border-b mb-6">
        <button
          onClick={() => setActiveTab("subscribe")}
          className={`px-6 py-3 font-medium transition-colors ${
            activeTab === "subscribe"
              ? "border-b-2 border-blue-500 text-blue-600"
              : "text-gray-600 hover:text-gray-800"
          }`}
        >
          Subscribe
        </button>
        <button
          onClick={() => setActiveTab("content")}
          className={`px-6 py-3 font-medium transition-colors ${
            activeTab === "content"
              ? "border-b-2 border-blue-500 text-blue-600"
              : "text-gray-600 hover:text-gray-800"
          }`}
        >
          Subscriber Only Content
        </button>
      </div>

      {/* Subscribe Tab */}
      {activeTab === "subscribe" && (
        <div className="flex flex-col items-center justify-center py-12">
          {subscriptionStatus?.has_access ? (
            <div className="text-center">
              <div className="mb-6">
                <svg
                  className="w-20 h-20 mx-auto text-green-500"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
              <h3 className="text-2xl font-bold text-gray-800 mb-2">
                You&apos;re Already Subscribed!
              </h3>
              <p className="text-gray-600 mb-4">
                Status: <span className="font-semibold capitalize">{subscriptionStatus.status}</span>
              </p>
              {subscriptionStatus.cancel_at_period_end ? (
                <p className="text-yellow-600">
                  Your subscription will expire at the end of the period
                </p>
              ) : (
                <p className="text-green-600">
                  Your subscription will renew automatically
                </p>
              )}
            </div>
          ) : (
            <div className="text-center">
              <h3 className="text-3xl font-bold text-gray-800 mb-8">
                Subscribe to Our Site
              </h3>
              <p className="text-gray-600 mb-8 max-w-md">
                Get unlimited access to premium features, exclusive content, and priority support
              </p>

              <button
                onClick={handleSubscribe}
                disabled={loading}
                className="bg-gradient-to-r from-blue-500 to-purple-600 text-white px-12 py-6 rounded-lg text-2xl font-bold hover:from-blue-600 hover:to-purple-700 disabled:opacity-50 transition-all transform hover:scale-105 shadow-lg"
              >
                {loading ? "Processing..." : "Subscribe Now"}
              </button>

              <div className="mt-8 text-sm text-gray-500">
                <p>Secure payment powered by Stripe</p>
              </div>
            </div>
          )}

          {error && (
            <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-600">
              {error}
            </div>
          )}

          {success && (
            <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg text-green-600">
              {success}
            </div>
          )}
        </div>
      )}

      {/* Subscriber Content Tab */}
      {activeTab === "content" && (
        <div className="py-4">
          {subscriptionStatus?.has_access ? (
            <div className="space-y-6">
              <div className="bg-gradient-to-r from-purple-50 to-blue-50 p-6 rounded-lg border border-purple-200">
                <div className="flex items-center gap-3 mb-4">
                  <svg
                    className="w-8 h-8 text-purple-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                    />
                  </svg>
                  <h3 className="text-2xl font-bold text-gray-800">
                    Premium Content Unlocked
                  </h3>
                </div>
                <p className="text-gray-700 mb-4">
                  Welcome to the exclusive subscriber area! Here&apos;s what you get:
                </p>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="bg-white p-5 rounded-lg border border-gray-200 shadow-sm">
                  <h4 className="font-bold text-lg text-gray-800 mb-2">
                    Advanced Analytics Dashboard
                  </h4>
                  <p className="text-gray-600 text-sm">
                    View detailed metrics, insights, and performance data for your account with
                    real-time updates and custom reports.
                  </p>
                  <div className="mt-3 text-blue-600 text-sm font-medium">
                    View Dashboard →
                  </div>
                </div>

                <div className="bg-white p-5 rounded-lg border border-gray-200 shadow-sm">
                  <h4 className="font-bold text-lg text-gray-800 mb-2">
                    Priority Support
                  </h4>
                  <p className="text-gray-600 text-sm">
                    Get 24/7 priority customer support with dedicated account managers and faster
                    response times.
                  </p>
                  <div className="mt-3 text-blue-600 text-sm font-medium">
                    Contact Support →
                  </div>
                </div>

                <div className="bg-white p-5 rounded-lg border border-gray-200 shadow-sm">
                  <h4 className="font-bold text-lg text-gray-800 mb-2">
                    Exclusive Resources
                  </h4>
                  <p className="text-gray-600 text-sm">
                    Access to premium templates, tools, guides, and educational content not
                    available to free users.
                  </p>
                  <div className="mt-3 text-blue-600 text-sm font-medium">
                    Browse Resources →
                  </div>
                </div>

                <div className="bg-white p-5 rounded-lg border border-gray-200 shadow-sm">
                  <h4 className="font-bold text-lg text-gray-800 mb-2">
                    Early Access Features
                  </h4>
                  <p className="text-gray-600 text-sm">
                    Be the first to try new features and updates before they&apos;re released to the
                    general public.
                  </p>
                  <div className="mt-3 text-blue-600 text-sm font-medium">
                    View Beta Features →
                  </div>
                </div>
              </div>

              <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
                <p className="text-sm text-yellow-800">
                  <strong>Note:</strong> This is demo subscriber content. In a real application,
                  this would contain actual premium features and content.
                </p>
              </div>
            </div>
          ) : (
            <div className="text-center py-12">
              <div className="mb-6">
                <svg
                  className="w-20 h-20 mx-auto text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                  />
                </svg>
              </div>
              <h3 className="text-2xl font-bold text-gray-800 mb-4">
                Subscription Required
              </h3>
              <p className="text-gray-600 mb-6">
                Subscribe to access exclusive content and premium features
              </p>
              <button
                onClick={() => setActiveTab("subscribe")}
                className="bg-blue-500 text-white px-8 py-3 rounded-lg hover:bg-blue-600 font-semibold"
              >
                Go to Subscribe Tab
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
