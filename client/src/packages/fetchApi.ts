import axios from 'axios';
import {useAuthStore} from "../stores/auth.ts";
import {useToast} from "@/hooks/useToast.ts";
import router from "@/routes";

const fetchApi = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL + '/api/',
  timeout: 60000,
  headers: {
    'Accept': 'application/json',
  }
})

fetchApi.interceptors.request.use(
  async config => {
    const {token} = useAuthStore();

    if (config?.headers) {
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }

    return config;
  },
  error => {
    return Promise.reject(error);
  },
);

fetchApi.interceptors.response.use(
  (response) => {
    if (response?.data?.message) {
      useToast().push({title: 'Success', message: response?.data?.message, kind: 'success'})
    }
    return response
  },
  async (error) => {
    const errorResponse = error?.response;

    if (errorResponse?.data?.errors) {
      const errors = errorResponse?.data?.errors
      for (const el in errors) {
        for (const error of errors[el]) {
          useToast().push({title: 'Error', message: error, kind: 'error'})
        }
      }
    } else if (error?.response?.data?.message) {
      const err = error?.response?.data?.message

      if (error?.response.status === 403) {
        useToast().push({title: 'Error', message: 'Forbidden', kind: 'error'})
      } else {
        useToast().push({title: 'Error', message: err, kind: 'error'})
      }
    } else {
      useToast().push({title: 'Error', message: error?.message, kind: 'error'})
    }

    if (errorResponse?.status === 401) {
      await router.push('/login');
      // window.location.href = '/login';
    }

    return Promise.reject(error)
  }
);

export default fetchApi;
