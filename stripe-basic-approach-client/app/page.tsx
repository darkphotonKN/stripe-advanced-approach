"use client";

import { useState } from "react";

export default function Home() {
  const [isLoading, setIsLoading] = useState(false);
  const [status, setStatus] = useState("");

  const handleCheckout = async () => {
    setIsLoading(true);
    setStatus("");

    try {
      const response = await fetch("/api/create-checkout-session", {
        method: "POST",
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }

      const { url } = await response.json();
      // redirect to the generated checkout url
      window.location.href = url;
    } catch (error) {
      console.error("Error:", error);
      setStatus("Error: " + (error as Error).message);
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-white flex items-center justify-center p-5">
      <div className="max-w-md w-full border border-gray-300 p-8 text-center rounded-lg">
        <h1 className="text-3xl font-bold mb-4">🦄 Cool T-Shirt</h1>
        <p className="text-gray-600 mb-4">
          The most amazing t-shirt you'll ever own!
        </p>
        <h2 className="text-2xl font-bold mb-6">$20.00</h2>
        <button
          onClick={handleCheckout}
          disabled={isLoading}
          className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white border-0 py-4 px-8 text-base rounded cursor-pointer w-full"
        >
          {isLoading ? "Creating checkout..." : "Buy Now"}
        </button>
        {status && <p className="mt-4 text-red-500">{status}</p>}
      </div>
    </div>
  );
}
