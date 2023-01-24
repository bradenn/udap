import {createApp} from 'vue'
import router from '@/router'
import Root from '@/App.vue'
// @ts-ignore
import MathJax, {initMathJax, renderByMathjax} from 'mathjax-vue3'
import '@/assets/bootstrap-grid.css';
import '@/assets/bootstrap-utilities.css';
import '@/assets/reset.css';

import '@/assets/sass/app.scss';

function onMathJaxReady() {
    const el = document.getElementById('mathjax')
    renderByMathjax(el)
}

initMathJax({url: 'https://cdnjs.cloudflare.com/ajax/libs/mathjax/3.2.2/es5/tex-svg-full.min.js'}, onMathJaxReady)

const app = createApp(Root)
// @ts-ignore
app.use(MathJax)

// app.config.warnHandler = function (msg, vm, trace) {
//     console.log(`Warn: ${msg}\nTrace: ${trace}`);
// }
//
// app.config.errorHandler = function (msg, vm, trace) {
//     console.error(`Error: ${msg}\nTrace: ${trace}`);
// }

app.use(router)

app.mount('#app')
