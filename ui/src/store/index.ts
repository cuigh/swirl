import { createStore, createLogger } from 'vuex'
import { Mutations } from "./mutations";

const debug = import.meta.env.DEV

export interface State {
    name: string | null;
    token: string | null;
    preference: {
        theme: string | null;
        locale: string | null;
    }
    ajaxLoading: boolean;
}

function initState(): State {
    const state: any = {
        name: localStorage.getItem("name"),
        token: localStorage.getItem("token"),
        ajaxLoading: false,
    }

    const locale = navigator.language.startsWith('zh') ? 'zh' : 'en'
    state.preference = { theme: 'light', locale: locale }
    try {
        state.preference = Object.assign(state.preference, JSON.parse(localStorage.getItem("preference") as string))
    } catch {
    }
    
    return state
}

export const store = createStore<State>({
    strict: debug,
    state: initState(),
    getters: {
        anonymous(state) {
            return !state.token
        }
    },
    mutations: {
        [Mutations.Login](state, user) {
            localStorage.setItem("name", user.name);
            localStorage.setItem("token", user.token);
            state.name = user.name;
            state.token = user.token;
        },
        [Mutations.Logout](state) {
            localStorage.removeItem("name");
            localStorage.removeItem("token");
            state.name = null;
            state.token = null;
        },
        [Mutations.SetToken](state, token) {
            localStorage.setItem("token", token);
            state.token = token;
        },
        [Mutations.SetPreference](state, preference) {
            localStorage.setItem("preference", JSON.stringify(preference));
            state.preference = preference;
        },
        [Mutations.SetAjaxLoading](state, loading) {
            state.ajaxLoading = loading;
        },
    },
    plugins: debug ? [createLogger()] : [],
})
