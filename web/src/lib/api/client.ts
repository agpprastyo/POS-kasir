import axios from "axios"
import {
    ActivityLogsApi,
    AuthApi,
    CancellationReasonsApi,
    CategoriesApi,
    Configuration,
    OrdersApi,
    PaymentMethodsApi,
    ProductsApi,
    PromotionsApi,
    ReportsApi,
    SettingsApi,
    UsersApi,
    PrinterApi
} from "@/lib/api/generated";



const BASE_PATH = import.meta.env.VITE_API_BASE!

if (!BASE_PATH) {
    throw new Error('VITE_API_BASE is not defined')
}

export const axiosInstance = axios.create({
    baseURL: BASE_PATH,
    withCredentials: true,
})

axiosInstance.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/auth/refresh') && !originalRequest.url?.includes('/auth/login')) {
            originalRequest._retry = true;

            try {
                await axiosInstance.post('/auth/refresh');
                return axiosInstance(originalRequest);
            } catch (refreshError) {
                console.error("Token refresh failed. Redirecting to login...", refreshError);
                if (!window.location.pathname.includes('/login')) {
                    const pathSegments = window.location.pathname.split('/').filter(Boolean);
                    const locale = (pathSegments.length > 0 && pathSegments[0].length === 2) ? pathSegments[0] : 'id';
                    window.location.href = `/${locale}/login`;
                }
                return Promise.reject(refreshError);
            }
        }

        if (error.response?.status === 401 && !window.location.pathname.includes('/login') && !originalRequest.url?.includes('/auth/login')) {

            if (!originalRequest.url?.includes('/auth/refresh')) {
                const pathSegments = window.location.pathname.split('/').filter(Boolean);
                const locale = (pathSegments.length > 0 && pathSegments[0].length === 2) ? pathSegments[0] : 'id';
                window.location.href = `/${locale}/login`;
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
export const promotionsApi = new PromotionsApi(config, undefined, axiosInstance)
export const reportsApi = new ReportsApi(config, undefined, axiosInstance)
export const activityLogsApi = new ActivityLogsApi(config, undefined, axiosInstance)
export const settingsApi = new SettingsApi(config, undefined, axiosInstance)
export const printerApi = new PrinterApi(config, undefined, axiosInstance)
