import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { Grid } from '@element-plus/icons-vue'

const app = createApp(App)

app.use(router)
   .use(store)
   .use(ElementPlus)
   .component('Grid', Grid)
   .mount('#app')
