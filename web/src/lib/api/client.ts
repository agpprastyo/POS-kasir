import axios from "axios"
import {
    AuthApi,
    CancellationReasonsApi,
    CategoriesApi,
    Configuration,
    OrdersApi,
    PaymentMethodsApi,
    ProductsApi,
    UsersApi
} from "@/lib/api/generated";



const BASE_PATH = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080/api/v1'

export const axiosInstance = axios.create({
    baseURL: BASE_PATH,
})

axiosInstance.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;


        if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/auth/refresh')) {
            originalRequest._retry = true;

            try {

                await axiosInstance.post('/auth/refresh');


                return axiosInstance(originalRequest);
            } catch (refreshError) {
                console.error("Token refresh failed. Redirecting to login...", refreshError);
                if (!window.location.pathname.includes('/login')) {
                    window.location.href = '/login';
                }
                return Promise.reject(refreshError);
            }
        }

        if (error.response?.status === 401 && !window.location.pathname.includes('/login')) {

            if (!originalRequest.url?.includes('/auth/refresh')) {
                window.location.href = '/login';
            }
        }

        return Promise.reject(error)
    }
)

const config = new Configuration({
    basePath: BASE_PATH,
    baseOptions: {
        withCredentials: true,
    },
})


export const authApi = new AuthApi(config, undefined, axiosInstance)
export const usersApi = new UsersApi(config, undefined, axiosInstance)
export const cancellationReasonsApi = new CancellationReasonsApi(config, undefined, axiosInstance)
export const categoriesApi = new CategoriesApi(config, undefined, axiosInstance)
export const productsApi = new ProductsApi(config, undefined, axiosInstance)
export const ordersApi = new OrdersApi(config, undefined, axiosInstance)
export const paymentMethodsApi = new PaymentMethodsApi(config, undefined, axiosInstance)

