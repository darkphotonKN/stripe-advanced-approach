"use client";

import { useState, useEffect } from "react";
import { productAPI, paymentAPI } from "@/lib/api";
import { stripePromise } from "@/lib/stripe";
import {
  CardElement,
  Elements,
  useStripe,
  useElements,
} from "@stripe/react-stripe-js";

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  price_id: string;
  type?: "one-time" | "subscription";
}

interface BuyProductProps {
  customerId: string | null;
  enabled: boolean;
  onProductPurchased?: (productName: string) => void;
}

function ProductPurchaseForm({
  customerId,
  products,
  onProductPurchased,
}: {
  customerId: string;
  products: Product[];
  onProductPurchased?: (productName: string) => void;
}) {
  const stripe = useStripe();
  const elements = useElements();
  const [loading, setLoading] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [paymentStatus, setPaymentStatus] = useState<
    "idle" | "succeeded" | "failed"
  >("idle");
  const [error, setError] = useState<string | null>(null);

  const handlePurchase = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements || !selectedProduct) return;

    setLoading(true);
    setError(null);

    try {
      let client_secret;
      
      if (selectedProduct.type === "subscription") {
        // Create subscription for recurring products
        const subscriptionResponse = await paymentAPI.subscribeToProduct(
          selectedProduct.id,
          customerId
        );
        client_secret = subscriptionResponse.client_secret;
      } else {
        // Create payment intent for one-time products
        const purchaseResponse = await paymentAPI.purchaseProduct(
          selectedProduct.id,
          customerId,
        );
        client_secret = purchaseResponse.client_secret;
      }

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
        if (onProductPurchased && selectedProduct) {
          onProductPurchased(selectedProduct.name);
        }
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
    setSelectedProduct(null);
  };

  if (paymentStatus === "succeeded") {
    return (
      <div>
        <p className="text-green-600 text-2xl mb-2">
          ✓ {selectedProduct?.type === 'subscription' ? 'Subscription Active!' : 'Purchase Successful!'}
        </p>
        <p className="text-gray-600 mb-4">
          {selectedProduct?.type === 'subscription' 
            ? `You have successfully subscribed to ${selectedProduct?.name} for $${(selectedProduct?.price || 0) / 100}/month`
            : `You have successfully purchased ${selectedProduct?.name} for $${(selectedProduct?.price || 0) / 100}`
          }
        </p>
        <button
          onClick={resetPayment}
          className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600"
        >
          {selectedProduct?.type === 'subscription' ? 'Subscribe to Another Product' : 'Buy Another Product'}
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Product Selection */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Select a Product
        </label>
        <div className="space-y-2">
          {products.map((product) => (
            <div
              key={product.id}
              className={`border rounded-lg p-4 cursor-pointer transition-colors ${selectedProduct?.id === product.id
                  ? "border-purple-500 bg-purple-50"
                  : "border-gray-300 hover:border-gray-400"
                }`}
              onClick={() => setSelectedProduct(product)}
            >
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <h4 className="font-semibold text-black">{product.name}</h4>
                    {product.type === "subscription" ? (
                      <span className="text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded-full font-medium">
                        Subscription
                      </span>
                    ) : (
                      <span className="text-xs px-2 py-1 bg-green-100 text-green-700 rounded-full font-medium">
                        One-time
                      </span>
                    )}
                  </div>
                  <p className="text-gray-600 text-sm">{product.description}</p>
                </div>
                <span className="font-bold text-purple-600 ml-4">
                  ${(product.price / 100).toFixed(2)}
                  {product.type === "subscription" && (
                    <span className="text-xs text-gray-500">/mo</span>
                  )}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Payment Form */}
      {selectedProduct && (
        <form onSubmit={handlePurchase} className="space-y-4">
          <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
            <p className="text-sm text-gray-700 mb-2">Selected Product:</p>
            <div className="flex items-center gap-2 mb-1">
              <p className="font-semibold text-black">{selectedProduct.name}</p>
              {selectedProduct.type === "subscription" ? (
                <span className="text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded-full font-medium">
                  Subscription
                </span>
              ) : (
                <span className="text-xs px-2 py-1 bg-green-100 text-green-700 rounded-full font-medium">
                  One-time
                </span>
              )}
            </div>
            <p className="text-purple-600 font-bold">
              Total: ${(selectedProduct.price / 100).toFixed(2)}
              {selectedProduct.type === "subscription" && (
                <span className="text-xs text-gray-500">/mo</span>
              )}
            </p>
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
            disabled={!stripe || loading || !selectedProduct}
            className="w-full bg-purple-500 text-white px-6 py-3 rounded hover:bg-purple-600 disabled:opacity-50 font-semibold"
          >
            {loading
              ? "Processing..."
              : selectedProduct.type === "subscription"
                ? `Subscribe for $${(selectedProduct.price / 100).toFixed(2)}/mo`
                : `Purchase for $${(selectedProduct.price / 100).toFixed(2)}`}
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
      )}
    </div>
  );
}

export default function BuyProduct({ customerId, enabled, onProductPurchased }: BuyProductProps) {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (enabled) {
      fetchProducts();
    }
  }, [enabled]);

  const fetchProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await productAPI.getProducts();
      setProducts(data.products || []);
    } catch (err) {
      setError("Failed to load products");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-4">Step 3: Purchase Products</h2>

      {loading ? (
        <p className="text-gray-500">Loading products...</p>
      ) : error ? (
        <div>
          <p className="text-red-500">{error}</p>
          <button
            onClick={fetchProducts}
            className="mt-2 text-blue-500 hover:text-blue-700"
          >
            Retry
          </button>
        </div>
      ) : products.length === 0 ? (
        <p className="text-gray-500">
          No products available. Please create products first in Step 1.
        </p>
      ) : (
        <div>
          <p className="text-gray-600 mb-4">
            Select and purchase one of the pre-created products
          </p>

          <Elements stripe={stripePromise}>
            {customerId ? (
              <ProductPurchaseForm
                customerId={customerId}
                products={products}
                onProductPurchased={onProductPurchased}
              />
            ) : (
              <div className="text-gray-500">Loading customer information...</div>
            )}
          </Elements>
        </div>
      )}
    </div>
  );
}
