import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  
  const customerId = localStorage.getItem('stripeCustomerId');
  if (customerId) {
    config.headers['X-Customer-Id'] = customerId;
  }
  return config;
});

export const productAPI = {
  setupProducts: async () => {
    const response = await api.post('/setup-products');
    return response.data;
  },
};

export const customerAPI = {
  create: async () => {
    const response = await api.post('/create-customer');
    return response.data;
  },
};

export const paymentMethodAPI = {
  saveCard: async (customerId: string) => {
    const response = await api.post('/save-card', {
      customer_id: customerId,
    });
    return response.data;
  },
};

export const paymentAPI = {
  createPaymentIntent: async (amount: number, customerId: string) => {
    const response = await api.post('/create-payment-intent', {
      amount,
      customer_id: customerId,
    });
    return response.data;
  },
};

export const subscriptionAPI = {
  create: async (priceId: string, customerId: string, email: string) => {
    const response = await api.post('/create-subscription', {
      price_id: priceId,
      customer_id: customerId,
      email,
    });
    return response.data;
  },
};

export const authAPI = {
  signUp: async (email: string, password: string, name: string) => {
    const response = await api.post('/signup', {
      email,
      password,
      name,
    });
    return response.data;
  },
  signIn: async (email: string, password: string) => {
    const response = await api.post('/signin', {
      email,
      password,
    });
    return response.data;
  },
};

export default api;