import {createApp} from 'vue'
import router from '@/router'
import Root from '@/App.vue'

import '@/assets/reset.css';
import '@/assets/bootstrap-grid.css';
import '@/assets/bootstrap-utilities.css';
import '@/assets/sass/element.css';
import '@/assets/sass/app.css';


const app = createApp(Root)

app.config.warnHandler = function (msg, vm, trace) {
    console.log(`Warn: ${msg}\nTrace: ${trace}`);
}

app.use(router)

app.mount('#app')
