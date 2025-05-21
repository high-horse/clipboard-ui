import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

import { Quasar, Notify, Dialog, Dark } from 'quasar'

// Import icon libraries
import '@quasar/extras/material-icons/material-icons.css'

// Import Quasar css
import 'quasar/src/css/index.sass'

// Import types
import type { QuasarPluginOptions } from 'quasar'

const myQuasarConfig: QuasarPluginOptions = {
  plugins: {
    Notify,
    Dialog
  },
  // Remove the config object if you're not using any specific configurations
  // or only include valid QuasarUIConfiguration properties
}

const app = createApp(App)
app.use(Quasar, myQuasarConfig)

// Set dark mode if needed
Dark.set(true) // or false to disable

app.mount('#app')