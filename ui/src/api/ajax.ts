import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { store } from "../store";
import { router } from "../router/router";
import { Mutations } from "@/store/mutations";
import { t, te } from '@/locales';

// export interface AjaxOptions {
// }

export interface Result<T> {
    code: number;
    info?: string;
    data?: T;
}

class Ajax {
    private ajax: AxiosInstance;

    constructor() {
        this.ajax = axios.create({
            baseURL: import.meta.env.MODE === 'development' ? '/api' : '/api',
            timeout: 10000,
            // withCredentials: true,            
        })

        this.ajax.interceptors.request.use(
            (config: any) => {
                if (store.state.user?.token) {
                    config.headers.Authorization = "Bearer " + store.state.user.token
                }
                // store.commit(Mutations.SetAjaxLoading, true);
                return config;
            },
            (error: any) => {
                return Promise.reject(error);
            }
        )

        this.ajax.interceptors.response.use(
            (response: any) => {
                // store.commit(Mutations.SetAjaxLoading, false);
                return response;
            },
            (error: any) => {
                if (this.handleError(error)) {
                    // Stop Promise chain
                    return new Promise(() => { })
                } else {
                    return Promise.reject(error)
                }
            }
        )
    }

    private handleError(error: any): boolean {
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    store.commit(Mutations.Logout);
                    if (error.config.method === "get") {
                        router.replace({
                            name: 'login',
                            query: {
                                redirect: router.currentRoute.value.fullPath
                            }
                        });
                    } else {
                        this.showError(error)
                    }
                    return true
                case 403:
                    router.replace("/403");
                    return true
                case 404:
                    router.replace("/404");
                    return true
                case 500:
                    this.showError(error)
            }
        } else {
            window.message.error(error.message, { duration: 5000 });
        }
        return false
    }

    private showError(error: any) {
        const code = error.response.data?.code || 1;
        const info = te('errors.'+code) ? t('errors.'+code) : error.response.data?.info || error.message;
        window.message.error(info, { duration: 5000 });
    }

    async get<T>(url: string, args?: any, config?: AxiosRequestConfig): Promise<Result<T>> {
        config = { ...config, params: args }
        const r = await this.ajax.get<Result<T>>(url, config);
        return r.data;
    }

    async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<Result<T>> {
        config = { ...config, headers: { 'Content-Type': 'application/json' } }
        // Object.assign(config || {}, {
        //     headers: {
        //         'Content-Type': 'application/json',
        //     },
        // })
        const r = await this.ajax.post<Result<T>>(url, data, config);
        return r.data;
    }

    async request<T>(config: AxiosRequestConfig): Promise<Result<T>> {
        const r = await this.ajax.request<Result<T>>(config);
        return r.data;
    }
}

export default new Ajax;