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


export const axiosInstance = axios.create({})

axiosInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            if (!window.location.pathname.includes('/login')) {
                console.error("Session expired. Redirecting to login...")
                window.location.href = '/login'
            }
        }
        return Promise.reject(error)
    }
)

const config = new Configuration({
    basePath: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8000/api/v1',
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

