import { createApp } from 'vue';
import { createPinia } from 'pinia'; 

import App from './App.vue';
import router from './router';

// Import UnoCSS entry styles
import '@unocss/reset/tailwind.css' // Import UnoCSS Reset (Tailwind compatibility)
import 'virtual:uno.css'

import './style.css';

const app = createApp(App);
const pinia = createPinia(); 

app.use(pinia); 
app.use(router); 

app.mount('#app');
