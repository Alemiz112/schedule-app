import Vue from "vue"
import Vuetify from "vuetify/lib"
import tailwind from "../../tailwind.config"

Vue.use(Vuetify)

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: tailwind.theme.colors.green,
        error: tailwind.theme.colors.red,
      },
      dark: {
        primary: tailwind.theme.colors.green,
        error: tailwind.theme.colors.red,
        background: '#1a1a1a',
        surface: '#2c2c2c',
      },
    },
  },
  breakpoint: {
    thresholds: {
      xs: 640,
      sm: 768,
      md: 1024,
      lg: 1280,
    },
    scrollBarWidth: 0,
  },
})
