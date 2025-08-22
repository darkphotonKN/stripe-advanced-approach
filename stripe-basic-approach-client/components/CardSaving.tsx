"use client";

import { useState } from "react";
import { paymentMethodAPI } from "@/lib/api";
import { stripePromise } from "@/lib/stripe";
import {
  CardElement,
  Elements,
  useStripe,
  useElements,
} from "@stripe/react-stripe-js";

interface CardSavingProps {
  customerId: string | null;
  enabled: boolean;
}

function CardSetupForm({ customerId }: { customerId: string }) {
  const stripe = useStripe();
  const elements = useElements();
  const [loading, setLoading] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements || !customerId) return;

    setLoading(true);
    setError(null);

    try {
      // gets the actual client secret after intent is made on BE
      const { client_secret } = await paymentMethodAPI.saveCard(customerId);

      const cardElement = elements.getElement(CardElement);
      if (!cardElement) return;

      // use stripe elements to save the card on stripe's system directly
      const { error: stripeError } = await stripe.confirmCardSetup(
        client_secret,
        {
          payment_method: {
            card: cardElement,
          },
        },
      );

      if (stripeError) {
        setError(stripeError.message || "Failed to save card");
      } else {
        setSaved(true);
      }
    } catch (err) {
      setError("Failed to save payment method");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (saved) {
    return (
      <div>
        <p className="text-green-600 mb-2">✓ Card saved successfully</p>
        <p className="text-sm text-gray-600">
          Card saved to customer for future use
        </p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div id="card-element" className="border p-3 rounded">
        <CardElement
          options={{
            style: {
              base: {
                fontSize: "16px",
                color: "#424770",
                "::placeholder": {
                  color: "#aab7c4",
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
        {loading ? "Saving..." : "Save Card"}
      </button>

      {error && <p className="text-red-500">{error}</p>}
    </form>
  );
}

export default function CardSaving({ customerId, enabled }: CardSavingProps) {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">
        Phase 3: Card Saving (Optional)
      </h2>

      {!enabled ? (
        <p className="text-gray-500">
          Create a customer first to save payment methods
        </p>
      ) : !customerId ? (
        <p className="text-gray-500">Customer ID required</p>
      ) : (
        <div>
          <p className="text-gray-600 mb-4">
            Save a card to the customer for future use (optional)
          </p>

          <Elements stripe={stripePromise}>
            <CardSetupForm customerId={customerId} />
          </Elements>
        </div>
      )}
    </div>
  );
}
