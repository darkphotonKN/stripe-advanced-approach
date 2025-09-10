"use client";

import { useState } from "react";
import { paymentAPI } from "@/lib/api";
import { stripePromise } from "@/lib/stripe";
import {
  CardElement,
  Elements,
  useStripe,
  useElements,
} from "@stripe/react-stripe-js";

interface OneTimePaymentProps {
  customerId: string | null;
  enabled: boolean;
}

function PaymentForm({ customerId }: { customerId: string }) {
  const stripe = useStripe();
  const elements = useElements();
  const [loading, setLoading] = useState(false);
  const [paymentStatus, setPaymentStatus] = useState<
    "idle" | "succeeded" | "failed"
  >("idle");
  const [error, setError] = useState<string | null>(null);
  const [productName, setProductName] = useState("Test Product");
  const [productDescription, setProductDescription] = useState("One-time payment test");
  const [amount, setAmount] = useState(20);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements) return;

    setLoading(true);
    setError(null);

    try {
      const { client_secret } = await paymentAPI.createPaymentIntent(
        amount * 100, // Convert to cents
        customerId,
      );

      console.log("Client secret for payment intent was:", client_secret);

      const cardElement = elements.getElement(CardElement);
      if (!cardElement) return;

      const { error: stripeError, paymentIntent } =
        await stripe.confirmCardPayment(client_secret, {
          payment_method: {
            card: cardElement,
          },
        });

      if (stripeError) {
        setError(stripeError.message || "Payment failed");
        setPaymentStatus("failed");
      } else if (paymentIntent?.status === "succeeded") {
        setPaymentStatus("succeeded");
      }
    } catch (err) {
      setError("Payment processing failed");
      setPaymentStatus("failed");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const resetPayment = () => {
    setPaymentStatus("idle");
    setError(null);
  };

  if (paymentStatus === "succeeded") {
    return (
      <div>
        <p className="text-green-600 text-2xl mb-2">✓ Payment Successful!</p>
        <p className="text-gray-600 mb-4">${amount.toFixed(2)} charged successfully for {productName}</p>
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
      <div className="space-y-3">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Product Name
          </label>
          <input
            type="text"
            value={productName}
            onChange={(e) => setProductName(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500 text-black placeholder-gray-700"
            placeholder="Enter product name"
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Description
          </label>
          <input
            type="text"
            value={productDescription}
            onChange={(e) => setProductDescription(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500 text-black placeholder-gray-700"
            placeholder="Enter description"
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Amount ($)
          </label>
          <input
            type="number"
            min="1"
            step="0.01"
            value={amount}
            onChange={(e) => setAmount(parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500 text-black placeholder-gray-700"
            placeholder="Enter amount"
          />
        </div>
      </div>

      <div className="border p-3 rounded">
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
        disabled={!stripe || loading || !productName || amount <= 0}
        className="bg-green-500 text-white px-6 py-3 rounded hover:bg-green-600 disabled:opacity-50 font-semibold"
      >
        {loading ? "Processing..." : `Pay $${amount.toFixed(2)}`}
      </button>

      {error && <div className="text-red-500">{error}</div>}

      {paymentStatus === "failed" && (
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

export default function OneTimePayment({
  customerId,
  enabled,
}: OneTimePaymentProps) {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Step 4: One-Time Payment</h2>

      {!enabled ? (
        <p className="text-gray-500">
          Create a customer first to process payments
        </p>
      ) : !customerId ? (
        <p className="text-gray-500">Customer ID required</p>
      ) : (
        <div>
          <p className="text-gray-600 mb-4">
            Process a one-time payment with custom product details
          </p>

          <Elements stripe={stripePromise}>
            <PaymentForm customerId={customerId} />
          </Elements>
        </div>
      )}
    </div>
  );
}
