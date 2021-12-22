import { createStore, createLogger } from 'vuex'
import { Mutations } from "./mutations";

const debug = import.meta.env.DEV

interface User {
    name: string;
    token: string;
    perms: Set<string>;
}

export interface State {
    user?: User | null;
    preference: {
        theme: string | null;
        locale: string | null;
    }
    ajaxLoading: boolean;
}

function loadObject(key: string) {
    let value = null
    try {
        value = JSON.parse(localStorage.getItem(key) as string)
    } catch {
    }
    return value
}

function initState(): State {
    const user = Object.assign({}, loadObject('user'))
    const locale = navigator.language.startsWith('zh') ? 'zh' : 'en'
    return {
        user: { perms: new Set(user.perms), name: user.name, token: user.token },
        preference: Object.assign({ theme: 'light', locale: locale }, loadObject('preference')),
        ajaxLoading: false,
    }
}

export const store = createStore<State>({
    strict: debug,
    state: initState(),
    getters: {
        anonymous(state) {
            return !state.user?.token
        },
        allow(state) {
            return (perm: string) => state.user?.perms?.has('*') || state.user?.perms?.has(perm)
        },
    },
    mutations: {
        [Mutations.Logout](state) {
            localStorage.removeItem("user");
            state.user = null;
        },
        [Mutations.SetUser](state, user) {
            localStorage.setItem("user", JSON.stringify(user));
            state.user = { perms: new Set(user.perms), name: user.name, token: user.token };
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
