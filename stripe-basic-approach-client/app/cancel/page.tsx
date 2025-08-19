import Link from 'next/link';

export default function Cancel() {
  return (
    <div className="min-h-screen bg-white flex items-center justify-center p-5">
      <div className="max-w-md w-full text-center">
        <div className="text-5xl text-red-600 mb-5">❌</div>
        <h1 className="text-3xl font-bold mb-4">Payment Cancelled</h1>
        <p className="text-gray-600 mb-6">Your payment was cancelled. No charges were made.</p>
        <Link href="/" className="text-blue-600 hover:underline">
          ← Try Again
        </Link>
      </div>
    </div>
  );
}