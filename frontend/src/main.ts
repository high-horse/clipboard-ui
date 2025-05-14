import { createApp } from 'vue'
import App from './App.vue'
import './style.css'
import { Quasar,QBtn,Notify, QInput, Dialog } from 'quasar'

// Import icon libraries (optional)
import '@quasar/extras/material-icons/material-icons.css'

// Import Quasar css
import 'quasar/src/css/index.sass'

createApp(App)
  .use(Quasar, {
    components: { QBtn, QInput },
    plugins: {Notify, Dialog}, // import Quasar plugins and add here if needed
    config: {
      dark: true // force dark mode
    }// optional if you want to use default config
  })
  .mount('#app')