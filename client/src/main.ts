import {createApp} from 'vue'
import router from '@/router'
import Root from '@/App.vue'

import 'bootstrap/dist/css/bootstrap-grid.css';
import 'bootstrap/dist/css/bootstrap-utilities.css';
import '@/assets/reset.css';

import '@/assets/sass/app.scss';


const app = createApp(Root)

app.config.warnHandler = function (msg, vm, trace) {
    console.log(`Warn: ${msg}\nTrace: ${trace}`);
}

app.config.errorHandler = function (msg, vm, trace) {
    console.error(`Error: ${msg}\nTrace: ${trace}`);
}

app.use(router)

app.mount('#app')
