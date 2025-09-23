"use client";

import { useState } from "react";
import { authAPI } from "@/lib/api";

interface SignInProps {
  onSuccess: (token: string) => void;
  onSwitchToSignUp: () => void;
}

export default function SignIn({ onSuccess, onSwitchToSignUp }: SignInProps) {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const data = await authAPI.signIn(formData.email, formData.password);
      localStorage.setItem("authToken", data.access_token);
      onSuccess(data.access_token);
    } catch (err: any) {
      setError(err.response?.data?.error || "Sign in failed");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
      <h2 className="text-2xl font-bold mb-6 text-center">Sign In</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">Email</label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-2">Password</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Signing In..." : "Sign In"}
        </button>
      </form>

      {error && <p className="text-red-500 mt-4 text-center">{error}</p>}

      <p className="text-center mt-6 text-sm text-gray-600">
        Don't have an account?{" "}
        <button
          onClick={onSwitchToSignUp}
          className="text-blue-500 hover:text-blue-700 font-medium"
        >
          Sign Up
        </button>
      </p>
    </div>
  );
}
