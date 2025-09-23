'use client';

import { useState, useEffect } from 'react';
import SignUp from './SignUp';
import SignIn from './SignIn';
import { customerAPI } from '@/lib/api';

interface AuthWrapperProps {
  children: React.ReactNode;
}

export default function AuthWrapper({ children }: AuthWrapperProps) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [authMode, setAuthMode] = useState<'signin' | 'signup'>('signin');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('authToken');
    setIsAuthenticated(!!token);
    setLoading(false);
  }, []);

  const checkAndSetExistingCustomer = async () => {
    try {
      const data = await customerAPI.getExisting();
      if (data.exists && data.stripe_customer_id) {
        localStorage.setItem('stripeCustomerId', data.stripe_customer_id);
      }
    } catch (err) {
      console.log('No existing Stripe customer found');
    }
  };

  const handleAuthSuccess = async (token: string) => {
    setIsAuthenticated(true);
    // Check for existing Stripe customer after successful login
    await checkAndSetExistingCustomer();
  };

  const handleSignOut = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('stripeCustomerId');
    localStorage.removeItem('subscriptionPriceId');
    setIsAuthenticated(false);
    window.location.reload();
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          {authMode === 'signin' ? (
            <SignIn
              onSuccess={handleAuthSuccess}
              onSwitchToSignUp={() => setAuthMode('signup')}
            />
          ) : (
            <SignUp
              onSuccess={handleAuthSuccess}
              onSwitchToSignIn={() => setAuthMode('signin')}
            />
          )}
        </div>
      </div>
    );
  }

  return (
    <div>
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <h1 className="text-xl font-semibold text-gray-900">
              Stripe Payment Flow
            </h1>
            <button
              onClick={handleSignOut}
              className="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium"
            >
              Sign Out
            </button>
          </div>
        </div>
      </nav>
      {children}
    </div>
  );
}