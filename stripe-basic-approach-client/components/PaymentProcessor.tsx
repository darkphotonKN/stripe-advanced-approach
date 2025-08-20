'use client';

import { useState } from 'react';
import { paymentAPI } from '@/lib/api';
import { stripePromise } from '@/lib/stripe';
import { CardElement, Elements, useStripe, useElements } from '@stripe/react-stripe-js';

interface PaymentProcessorProps {
  customerId: string | null;
}

function PaymentForm({ customerId }: { customerId: string | null }) {
  const stripe = useStripe();
  const elements = useElements();
  const [amount, setAmount] = useState('10.00');
  const [loading, setLoading] = useState(false);
  const [paymentStatus, setPaymentStatus] = useState<'idle' | 'processing' | 'succeeded' | 'failed'>('idle');
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements) return;

    setLoading(true);
    setPaymentStatus('processing');
    setError(null);

    try {
      const amountInCents = Math.round(parseFloat(amount) * 100);
      
      const { client_secret } = await paymentAPI.createPaymentIntent(
        amountInCents,
        customerId || undefined
      );

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
      <div className="text-center">
        <div className="text-green-600 text-6xl mb-4">✓</div>
        <h3 className="text-2xl font-bold text-green-600 mb-2">Payment Successful!</h3>
        <p className="text-gray-600 mb-4">Amount paid: ${amount}</p>
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
      <div>
        <label className="block text-sm font-medium mb-2">Amount ($)</label>
        <input
          type="number"
          step="0.01"
          min="0.50"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          className="w-full px-3 py-2 border rounded"
          disabled={loading}
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-2">Card Details</label>
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
      </div>

      <button
        type="submit"
        disabled={!stripe || loading}
        className={`w-full py-3 rounded font-semibold transition ${
          loading
            ? 'bg-gray-400 text-gray-200 cursor-not-allowed'
            : 'bg-green-500 text-white hover:bg-green-600'
        }`}
      >
        {loading ? 'Processing...' : `Pay $${amount}`}
      </button>

      {error && (
        <div className="text-red-500 text-center">{error}</div>
      )}

      {paymentStatus === 'failed' && (
        <button
          type="button"
          onClick={resetPayment}
          className="w-full text-gray-500 hover:text-gray-700"
        >
          Try Again
        </button>
      )}
    </form>
  );
}

export default function PaymentProcessor({ customerId }: PaymentProcessorProps) {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Phase 3: Payment Processing</h2>
      
      <p className="text-gray-600 mb-4">
        Enter payment details and amount to process an immediate payment
      </p>

      <Elements stripe={stripePromise}>
        <PaymentForm customerId={customerId} />
      </Elements>

      <div className="mt-4 p-4 bg-gray-50 rounded text-sm text-gray-600">
        <p className="font-semibold mb-1">Test Card Numbers:</p>
        <ul className="space-y-1">
          <li>✓ Success: 4242 4242 4242 4242</li>
          <li>✗ Decline: 4000 0000 0000 0002</li>
          <li>Use any future date and any 3-digit CVC</li>
        </ul>
      </div>
    </div>
  );
}