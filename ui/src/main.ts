import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router/router'
import { store } from './store'
import i18n from './locales'

const app = createApp(App).use(router).use(store).use(i18n);
app.mount('#app');
