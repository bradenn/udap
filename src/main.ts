import {createApp} from 'vue'
import Root from '@/App.vue'
import 'udap-ui/assets/general.scss';
import "./assets/app.scss"
import 'udap-ui/assets/bootstrap-grid.css';
import 'udap-ui/assets/bootstrap-utilities.css';
import 'udap-ui/assets/reset.css';

import router from '@/router'


if (typeof window !== 'undefined')
    import('./pwa')


const app = createApp(Root)
app.use(router)
app.mount('#app')
