import {createApp} from 'vue'
import App from './App.vue'
import router from './router/index.js'

import 'bootstrap/scss/bootstrap.scss'
import 'sweetalert2/src/sweetalert2.scss'

createApp(App).//
    use(router).//
    mount('#app')