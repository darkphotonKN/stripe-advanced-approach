'use client';

import { useState } from 'react';
import { paymentMethodAPI } from '@/lib/api';
import { stripePromise } from '@/lib/stripe';
import { CardElement, Elements, useStripe, useElements } from '@stripe/react-stripe-js';

interface PaymentMethodManagerProps {
  customerId: string | null;
  onPaymentMethodSaved: () => void;
}

function CardSetupForm({ customerId, onSuccess }: { customerId: string; onSuccess: () => void }) {
  const stripe = useStripe();
  const elements = useElements();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements || !customerId) return;

    setLoading(true);
    setError(null);

    try {
      const { client_secret } = await paymentMethodAPI.createSetupIntent(customerId);
      
      const cardElement = elements.getElement(CardElement);
      if (!cardElement) return;

      const { error: stripeError } = await stripe.confirmCardSetup(client_secret, {
        payment_method: {
          card: cardElement,
        },
      });

      if (stripeError) {
        setError(stripeError.message || 'Failed to save card');
      } else {
        onSuccess();
      }
    } catch (err) {
      setError('Failed to save payment method');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

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
        className="bg-green-500 text-white px-6 py-2 rounded hover:bg-green-600 disabled:opacity-50"
      >
        {loading ? 'Saving...' : 'Save Card for Future Use'}
      </button>
      
      {error && <p className="text-red-500">{error}</p>}
    </form>
  );
}

export default function PaymentMethodManager({ customerId, onPaymentMethodSaved }: PaymentMethodManagerProps) {
  const [showCardForm, setShowCardForm] = useState(false);
  const [savedCards, setSavedCards] = useState<any[]>([]);
  const [loadingCards, setLoadingCards] = useState(false);

  const loadSavedCards = async () => {
    if (!customerId) return;
    
    setLoadingCards(true);
    try {
      const data = await paymentMethodAPI.list(customerId);
      setSavedCards(data.payment_methods || []);
    } catch (err) {
      console.error('Failed to load payment methods:', err);
    } finally {
      setLoadingCards(false);
    }
  };

  const handleCardSaved = () => {
    setShowCardForm(false);
    loadSavedCards();
    onPaymentMethodSaved();
  };

  const removeCard = async (paymentMethodId: string) => {
    try {
      await paymentMethodAPI.detach(paymentMethodId);
      setSavedCards(savedCards.filter(card => card.id !== paymentMethodId));
    } catch (err) {
      console.error('Failed to remove card:', err);
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Phase 2: Payment Method (Optional)</h2>
      
      {!customerId ? (
        <p className="text-gray-500">Create a customer first to save payment methods</p>
      ) : (
        <div>
          <p className="text-gray-600 mb-4">
            Save a card for future use (optional - you can also pay with a new card each time)
          </p>
          
          {!showCardForm ? (
            <div>
              <button
                onClick={() => {
                  setShowCardForm(true);
                  loadSavedCards();
                }}
                className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600"
              >
                Add Payment Method
              </button>
              
              {savedCards.length > 0 && (
                <div className="mt-4">
                  <h3 className="font-semibold mb-2">Saved Cards:</h3>
                  {savedCards.map((card) => (
                    <div key={card.id} className="flex justify-between items-center p-2 border rounded mb-2">
                      <span>•••• {card.card?.last4}</span>
                      <button
                        onClick={() => removeCard(card.id)}
                        className="text-red-500 hover:text-red-700"
                      >
                        Remove
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ) : (
            <div>
              <Elements stripe={stripePromise}>
                <CardSetupForm customerId={customerId} onSuccess={handleCardSaved} />
              </Elements>
              <button
                onClick={() => setShowCardForm(false)}
                className="mt-2 text-gray-500 hover:text-gray-700"
              >
                Cancel
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}