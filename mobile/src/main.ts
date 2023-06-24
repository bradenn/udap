import {createApp} from 'vue'
import Root from '@/App.vue'
import 'udap-ui/assets/general.scss';
import "./assets/app.scss"
import 'udap-ui/assets/bootstrap-grid.css';
import 'udap-ui/assets/bootstrap-utilities.css';
import 'udap-ui/assets/reset.css';

import router from '@/router'

// import updateSW from "@/registerServiceWorker";
// @ts-ignore
import {registerSW} from "virtual:pwa-register";

window.addEventListener('load', () => {
    registerSW({
        onNeedRefresh() {
            console.log("REFRESH")
        },
        onOfflineReady() {
            console.log("OFFLINE")
        },
        immediate: true
    })
})

const app = createApp(Root)
app.use(router)
app.mount('#app')
