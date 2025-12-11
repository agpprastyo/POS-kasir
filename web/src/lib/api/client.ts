import {AuthApi, Configuration, UsersApi} from "@/lib/api/generated";


const config = new Configuration({
    basePath: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8000/api/v1',
    baseOptions: {
        withCredentials: true,
    },
})

export const authApi = new AuthApi(config)
export const usersApi = new UsersApi(config)
