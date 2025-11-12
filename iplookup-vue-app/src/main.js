import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

// 引入Font Awesome图标库
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faSearch, faLocationArrow, faTrash, faCopy, faExclamationCircle } from '@fortawesome/free-solid-svg-icons'

// 添加需要的图标
library.add(faSearch, faLocationArrow, faTrash, faCopy, faExclamationCircle)

const app = createApp(App)
app.component('font-awesome-icon', FontAwesomeIcon)
app.mount('#app')