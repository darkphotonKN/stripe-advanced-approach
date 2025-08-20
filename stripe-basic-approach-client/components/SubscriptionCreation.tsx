'use client';

import { useState } from 'react';
import { subscriptionAPI } from '@/lib/api';
import { stripePromise } from '@/lib/stripe';
import { CardElement, Elements, useStripe, useElements } from '@stripe/react-stripe-js';

interface SubscriptionCreationProps {
  customerId: string | null;
  priceId: string | null;
  enabled: boolean;
}

function SubscriptionForm({ customerId, priceId }: { customerId: string; priceId: string }) {
  const stripe = useStripe();
  const elements = useElements();
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [subscriptionStatus, setSubscriptionStatus] = useState<'idle' | 'active' | 'failed'>('idle');
  const [subscriptionId, setSubscriptionId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements || !email) return;

    setLoading(true);
    setError(null);

    try {
      const { subscription_id, client_secret } = await subscriptionAPI.create(
        priceId,
        customerId,
        email
      );

      const cardElement = elements.getElement(CardElement);
      if (!cardElement) return;

      const { error: stripeError, paymentIntent } = await stripe.confirmCardPayment(client_secret, {
        payment_method: {
          card: cardElement,
        },
      });

      if (stripeError) {
        setError(stripeError.message || 'Subscription failed');
        setSubscriptionStatus('failed');
      } else if (paymentIntent?.status === 'succeeded') {
        setSubscriptionId(subscription_id);
        setSubscriptionStatus('active');
      }
    } catch (err) {
      setError('Subscription creation failed');
      setSubscriptionStatus('failed');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const resetSubscription = () => {
    setSubscriptionStatus('idle');
    setSubscriptionId(null);
    setEmail('');
    setError(null);
  };

  if (subscriptionStatus === 'active') {
    return (
      <div>
        <p className="text-green-600 text-2xl mb-2">✓ Subscription Active!</p>
        <p className="text-gray-600 mb-2">Subscription ID: {subscriptionId}</p>
        <p className="text-gray-600 mb-4">Email: {email}</p>
        <button
          onClick={resetSubscription}
          className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600"
        >
          Create Another Subscription
        </button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-2">Email</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full px-3 py-2 border rounded"
          placeholder="customer@example.com"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-2">Card Details</label>
        <div id="subscription-card-element" className="border p-3 rounded">
          <CardElement
            options={{
              style: {
                base: {
                  fontSize: '16px',
                  color: '#424770',
                  '::placeholder': {
                    color: '#aab7c4',
                  },
                },
              },
            }}
          />
        </div>
      </div>

      <button
        type="submit"
        disabled={!stripe || loading || !email}
        className="bg-purple-500 text-white px-6 py-3 rounded hover:bg-purple-600 disabled:opacity-50 font-semibold w-full"
      >
        {loading ? 'Processing...' : 'Subscribe'}
      </button>

      {error && (
        <div className="text-red-500">{error}</div>
      )}

      {subscriptionStatus === 'failed' && (
        <button
          type="button"
          onClick={resetSubscription}
          className="text-gray-500 hover:text-gray-700 w-full"
        >
          Try Again
        </button>
      )}
    </form>
  );
}

export default function SubscriptionCreation({ customerId, priceId, enabled }: SubscriptionCreationProps) {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Phase 5: Subscription Creation</h2>
      
      {!enabled ? (
        <p className="text-gray-500">Complete previous phases first</p>
      ) : !customerId || !priceId ? (
        <p className="text-gray-500">Customer ID and Price ID required</p>
      ) : (
        <div>
          <p className="text-gray-600 mb-4">
            Create a subscription with immediate first payment
          </p>
          
          <Elements stripe={stripePromise}>
            <SubscriptionForm customerId={customerId} priceId={priceId} />
          </Elements>
        </div>
      )}
    </div>
  );
}