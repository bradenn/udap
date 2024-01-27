import 'udap-ui/assets/general.scss';

import {createApp} from 'vue'
import Root from '@/App.vue'

import "./assets/app.scss"
import 'udap-ui/assets/bootstrap-grid.css';
import 'udap-ui/assets/bootstrap-utilities.css';
import 'udap-ui/assets/reset.css';

import router from '@/router'


const app = createApp(Root)
app.use(router)
app.mount('#app')
