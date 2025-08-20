import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const customerId = localStorage.getItem('customerId');
  if (customerId) {
    config.headers['X-Customer-Id'] = customerId;
  }
  return config;
});

export const customerAPI = {
  create: async () => {
    const response = await api.post('/customers');
    return response.data;
  },
  get: async (customerId: string) => {
    const response = await api.get(`/customers/${customerId}`);
    return response.data;
  },
};

export const paymentMethodAPI = {
  createSetupIntent: async (customerId: string) => {
    const response = await api.post('/payment-methods/setup-intent', {
      customer_id: customerId,
    });
    return response.data;
  },
  list: async (customerId: string) => {
    const response = await api.get(`/payment-methods/${customerId}`);
    return response.data;
  },
  detach: async (paymentMethodId: string) => {
    const response = await api.delete(`/payment-methods/${paymentMethodId}`);
    return response.data;
  },
};

export const paymentAPI = {
  createPaymentIntent: async (amount: number, customerId?: string, paymentMethodId?: string) => {
    const response = await api.post('/payments/create-intent', {
      amount,
      customer_id: customerId,
      payment_method_id: paymentMethodId,
    });
    return response.data;
  },
  confirmPayment: async (paymentIntentId: string) => {
    const response = await api.post('/payments/confirm', {
      payment_intent_id: paymentIntentId,
    });
    return response.data;
  },
};

export default api;