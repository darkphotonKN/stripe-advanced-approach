'use client';

import { useState } from 'react';
import { paymentAPI } from '@/lib/api';
import { stripePromise } from '@/lib/stripe';
import { CardElement, Elements, useStripe, useElements } from '@stripe/react-stripe-js';

interface OneTimePaymentProps {
  customerId: string | null;
  enabled: boolean;
}

function PaymentForm({ customerId }: { customerId: string }) {
  const stripe = useStripe();
  const elements = useElements();
  const [loading, setLoading] = useState(false);
  const [paymentStatus, setPaymentStatus] = useState<'idle' | 'succeeded' | 'failed'>('idle');
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements) return;

    setLoading(true);
    setError(null);

    try {
      const { client_secret } = await paymentAPI.createPaymentIntent(2000, customerId);

      const cardElement = elements.getElement(CardElement);
      if (!cardElement) return;

      const { error: stripeError, paymentIntent } = await stripe.confirmCardPayment(client_secret, {
        payment_method: {
          card: cardElement,
        },
      });

      if (stripeError) {
        setError(stripeError.message || 'Payment failed');
        setPaymentStatus('failed');
      } else if (paymentIntent?.status === 'succeeded') {
        setPaymentStatus('succeeded');
      }
    } catch (err) {
      setError('Payment processing failed');
      setPaymentStatus('failed');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const resetPayment = () => {
    setPaymentStatus('idle');
    setError(null);
  };

  if (paymentStatus === 'succeeded') {
    return (
      <div>
        <p className="text-green-600 text-2xl mb-2">✓ Payment Successful!</p>
        <p className="text-gray-600 mb-4">$20.00 charged successfully</p>
        <button
          onClick={resetPayment}
          className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600"
        >
          Make Another Payment
        </button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="border p-3 rounded">
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

      <button
        type="submit"
        disabled={!stripe || loading}
        className="bg-green-500 text-white px-6 py-3 rounded hover:bg-green-600 disabled:opacity-50 font-semibold"
      >
        {loading ? 'Processing...' : 'Pay $20.00'}
      </button>

      {error && (
        <div className="text-red-500">{error}</div>
      )}

      {paymentStatus === 'failed' && (
        <button
          type="button"
          onClick={resetPayment}
          className="text-gray-500 hover:text-gray-700"
        >
          Try Again
        </button>
      )}
    </form>
  );
}

export default function OneTimePayment({ customerId, enabled }: OneTimePaymentProps) {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Phase 4: One-Time Payment</h2>
      
      {!enabled ? (
        <p className="text-gray-500">Create a customer first to process payments</p>
      ) : !customerId ? (
        <p className="text-gray-500">Customer ID required</p>
      ) : (
        <div>
          <p className="text-gray-600 mb-4">
            Process a one-time payment of $20.00
          </p>
          
          <Elements stripe={stripePromise}>
            <PaymentForm customerId={customerId} />
          </Elements>
        </div>
      )}
    </div>
  );
}