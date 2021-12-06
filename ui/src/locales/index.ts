import { createI18n } from 'vue-i18n'
import en from "./en";
import zh from "./zh";
import { store } from "@/store";

const i18n = createI18n({
    legacy: false, // you must set `false`, to use Composition API
    locale: store.state.preference?.locale || 'en',
    fallbackLocale: 'en',
    messages: {
        en,
        zh,
    }
})

export const t = i18n.global.t
export const te = i18n.global.te
export default i18n