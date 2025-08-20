'use client';

import { useState } from 'react';
import { authAPI } from '@/lib/api';

interface SignUpProps {
  onSuccess: (token: string) => void;
  onSwitchToSignIn: () => void;
}

export default function SignUp({ onSuccess, onSwitchToSignIn }: SignUpProps) {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    name: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const data = await authAPI.signUp(formData.email, formData.password, formData.name);
      localStorage.setItem('authToken', data.token);
      onSuccess(data.token);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Sign up failed');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
      <h2 className="text-2xl font-bold mb-6 text-center">Sign Up</h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">Name</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        
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
            minLength={6}
          />
          <p className="text-xs text-gray-500 mt-1">Minimum 6 characters</p>
        </div>
        
        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? 'Creating Account...' : 'Sign Up'}
        </button>
      </form>
      
      {error && (
        <p className="text-red-500 mt-4 text-center">{error}</p>
      )}
      
      <p className="text-center mt-6 text-sm text-gray-600">
        Already have an account?{' '}
        <button
          onClick={onSwitchToSignIn}
          className="text-blue-500 hover:text-blue-700 font-medium"
        >
          Sign In
        </button>
      </p>
    </div>
  );
}