import {createApp} from 'vue'
import router from '@/router'
import Root from '@/App.vue'

const app = createApp(Root)

app.config.warnHandler = function (msg, vm, trace) {
    console.log(`Warn: ${msg}\nTrace: ${trace}`);
}

app.use(router)

app.mount('#app')
