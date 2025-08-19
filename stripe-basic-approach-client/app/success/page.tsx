"use client";

import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { Suspense } from "react";

// callback after user checksout on stripe payment page
function SuccessContent() {
  const searchParams = useSearchParams();
  const sessionId = searchParams.get("session_id");

  return (
    <div className="min-h-screen bg-white flex items-center justify-center p-5">
      <div className="max-w-md w-full text-center">
        <div className="text-5xl text-green-600 mb-5">✅</div>
        <h1 className="text-3xl font-bold mb-4">Payment Successful!</h1>
        <p className="text-gray-600 mb-6">
          Thank you for your purchase. Your payment has been processed.
        </p>

        <div className="bg-gray-50 p-5 rounded-lg mb-6">
          <h3 className="text-xl font-semibold mb-3">Order Details</h3>
          <p className="mb-2">
            <strong>Product:</strong> Cool T-Shirt
          </p>
          <p className="mb-2">
            <strong>Amount:</strong> $20.00
          </p>
          <p>
            <strong>Session ID:</strong> {sessionId || "N/A"}
          </p>
        </div>

        <p className="text-gray-600 mb-4">
          You should receive a confirmation email shortly.
        </p>
        <Link href="/" className="text-blue-600 hover:underline">
          ← Back to Shop
        </Link>
      </div>
    </div>
  );
}

export default function Success() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <SuccessContent />
    </Suspense>
  );
}
