import {createApp} from 'vue'
import Root from '@/App.vue'
import "./assets/app.scss"
import './assets/bootstrap-grid.css';
import './assets/bootstrap-utilities.css';
import './assets/reset.css';
import router from '@/router'
import './registerServiceWorker';

const app = createApp(Root)
app.use(router)
app.mount('#app')
