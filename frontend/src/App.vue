<template>
  <v-app>
    <AutoSnackbar color="error" :text="error" />
    <AutoSnackbar color="tw-bg-blue" :text="info" />
    <SignInNotSupportedDialog v-model="webviewDialog" />
    <SignInDialog
      v-model="signInDialog"
      @signIn="_signIn"
      @emailSignIn="_emailSignIn"
    />
    <NewDialog
      v-model="newDialogOptions.show"
      :type="newDialogOptions.openNewGroup ? 'group' : 'event'"
      :contactsPayload="newDialogOptions.contactsPayload"
      :no-tabs="newDialogOptions.eventOnly"
      :folder-id="newDialogOptions.folderId"
    />
    <div
      v-if="showHeader"
      class="tw-fixed tw-z-40 tw-h-14 tw-w-screen tw-bg-white sm:tw-h-16"
      dark
    >
      <div
        class="tw-relative tw-m-auto tw-flex tw-h-full tw-max-w-6xl tw-items-center tw-justify-center tw-px-4"
      >
        <router-link :to="{ name: 'home' }">
          <Logo type="timeful" />
        </router-link>
        <v-spacer />

        <v-btn
          v-if="($route.name === 'home' || $route.name === 'event') && !isPhone"
          color="primary"
          class="tw-mx-2 tw-rounded-md"
          :style="{
            boxShadow: '0px 2px 8px 0px #00994C80 !important',
          }"
          @click="() => _createNew()"
        >
          + Create new
        </v-btn>
        <v-btn icon @click="setDarkMode(!darkMode)" class="tw-mr-1">
          <v-icon>{{ darkMode ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
        </v-btn>
        <div v-if="authUser" class="sm:tw-ml-4">
          <AuthUserMenu />
        </div>
        <v-btn v-else id="top-right-sign-in-btn" text @click="signIn">
          Sign in
        </v-btn>
      </div>
    </div>

    <v-main>
      <div class="tw-flex tw-h-screen tw-flex-col">
        <div
          class="tw-relative tw-flex-1 tw-overscroll-auto"
          :class="routerViewClass"
        >
          <router-view v-if="loaded" :key="$route.fullPath" />
        </div>
      </div>
    </v-main>
  </v-app>
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=DM+Sans&display=swap");

html {
  overflow-y: auto !important;
  /* overscroll-behavior: none; */
  scroll-behavior: smooth;
}

* {
  font-family: "DM Sans", sans-serif;
  /* touch-action: manipulation !important; */
}

.v-messages__message {
  font-size: theme("fontSize.xs");
  line-height: 1.25;
}
.v-input--selection-controls {
  margin-top: 0px !important;
  padding-top: 0px !important;
}

/** Buttons */
.v-btn {
  letter-spacing: unset !important;
  text-transform: unset !important;
}
.v-btn:not(.v-btn--round, .v-btn-toggle > .v-btn).v-size--default {
  height: 38px !important;
  border-radius: theme("borderRadius.md") !important;
}

.v-btn.v-btn--is-elevated {
  -webkit-box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.15) !important;
  -moz-box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.15) !important;
  box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.15) !important;
  border: 1px solid theme("colors.light-gray-stroke");
}

.v-btn.v-btn--is-elevated.tw-bg-white {
  -webkit-box-shadow: 0px 1px 4px 0px rgba(0, 0, 0, 0.25) !important;
  -moz-box-shadow: 0px 1px 4px 0px rgba(0, 0, 0, 0.25) !important;
  box-shadow: 0px 1px 4px 0px rgba(0, 0, 0, 0.25) !important;
  border: 1px solid theme("colors.off-white");
}

.v-btn.v-btn--is-elevated.primary,
.v-btn.v-btn--is-elevated.tw-bg-green,
.v-btn.v-btn--is-elevated.tw-bg-white.tw-text-green {
  -webkit-box-shadow: 0px 2px 8px 0px #00994c80 !important;
  -moz-box-shadow: 0px 2px 8px 0px #00994c80 !important;
  box-shadow: 0px 2px 8px 0px #00994c80 !important;
  border: 1px solid theme("colors.light-green") !important;
}

.v-btn.v-btn--is-elevated.tw-bg-very-dark-gray {
  -webkit-box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.25) !important;
  -moz-box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.25) !important;
  box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.25) !important;
  border: 1px solid theme("colors.dark-gray") !important;
}

.v-btn.v-btn--is-elevated.tw-bg-blue,
.v-btn.v-btn--is-elevated.tw-bg-white.tw-text-blue {
  -webkit-box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.25) !important;
  -moz-box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.25) !important;
  box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.25) !important;
  border: 1px solid theme("colors.light-blue") !important;
}

/** Drop shadows */
.v-text-field.v-text-field--solo:not(.v-text-field--solo-flat)
  > .v-input__control
  > .v-input__slot {
  filter: drop-shadow(0 0.5px 2px rgba(0, 0, 0, 0.1)) !important;
  box-shadow: inset 0 -1px 0 0 rgba(0, 0, 0, 0.1) !important;
  border-radius: theme("borderRadius.md") !important;
  border: 1px solid #4f4f4f1f !important;
}
.v-menu__content {
  box-shadow: 0px 5px 5px -1px rgba(0, 0, 0, 0.1),
    0px 8px 10px 0.5px rgba(0, 0, 0, 0.07), 0px 3px 14px 1px rgba(0, 0, 0, 0.06) !important;
}
.overlay-avail-shadow-green {
  box-shadow: 0px 3px 6px 0px #1c7d454d !important;
}
.overlay-avail-shadow-yellow {
  box-shadow: 0px 2px 8px 0px #e5a8004d !important;
}

/** Switch  */
.v-input--switch--inset .v-input--selection-controls__input {
  margin-right: 0 !important;
  transform: scale(80%) !important;
}
.v-input--switch__track.primary--text {
  border: 2px theme("colors.light-green") solid !important;
}
.v-input--switch__track {
  border: 2px theme("colors.gray") solid !important;
  background-color: theme("colors.gray") !important;
  box-shadow: 0px 0.74px 4.46px 0px rgba(0, 0, 0, 0.1) !important;
}
.v-input--is-label-active .v-input--switch__track {
  background-color: currentColor !important;
  box-shadow: 0px 1.5px 4.5px 0px rgba(0, 0, 0, 0.2) !important;
}
.v-input--switch--inset .v-input--switch__track,
.v-input--switch--inset .v-input--selection-controls__input {
  opacity: 1 !important;
}
.v-input--switch__thumb {
  background-color: white !important;
}
.v-text-field__details {
  padding: 0 !important;
}

/** Error color */
.error--text .v-input__slot {
  outline: red solid;
  border-radius: 3px;
}

/** Dark mode overrides */
html.dark {
  color-scheme: dark;
}

html.dark .tw-bg-white {
  background-color: #2c2c2c !important;
}

html.dark .tw-bg-off-white {
  background-color: #383838 !important;
}

html.dark .tw-bg-light-gray {
  background-color: #404040 !important;
}

html.dark .tw-text-black {
  color: #e8e8e8 !important;
}

html.dark .tw-text-dark-gray {
  color: #c0c0c0 !important;
}

html.dark .tw-text-very-dark-gray {
  color: #d8d8d8 !important;
}

html.dark .tw-border-light-gray-stroke {
  border-color: #4a4a4a !important;
}

html.dark .tw-border-light-gray {
  border-color: #4a4a4a !important;
}

html.dark .tw-border-gray {
  border-color: #555555 !important;
}

html.dark .tw-drop-shadow {
  filter: drop-shadow(0 1px 3px rgba(0, 0, 0, 0.5)) !important;
}

html.dark .tw-drop-shadow-md {
  filter: drop-shadow(0 2px 6px rgba(0, 0, 0, 0.6)) !important;
}

html.dark .v-btn.v-btn--is-elevated {
  border-color: #4a4a4a !important;
}

html.dark .v-btn.v-btn--is-elevated.tw-bg-white {
  background-color: #383838 !important;
  border-color: #4a4a4a !important;
}

html.dark .tw-text-dark-green {
  color: #4fc87a !important;
}

html.dark .tw-text-green {
  color: #4fc87a !important;
}

/* The Dashboard/Tools panels use sm:tw-bg-[#f3f3f366] — a semi-transparent
   light gray that blends with the dark Vuetify app bg into a muddy gray.
   Override it to a proper dark surface. */
html.dark .sm\:tw-bg-\[#f3f3f366\] {
  background-color: #2c2c2c !important;
}

html.dark .animate-boba {
  display: none !important;
}
</style>

<script>
import { mapMutations, mapState, mapActions, mapGetters } from "vuex"
import {
  get,
  getLocation,
  isPhone,
  post,
  signInGoogle,
  signInOutlook,
} from "@/utils"
import {
  authTypes,
  calendarTypes,
  eventTypes,
} from "@/constants"
import AutoSnackbar from "@/components/AutoSnackbar"
import AuthUserMenu from "@/components/AuthUserMenu.vue"
import SignInNotSupportedDialog from "@/components/SignInNotSupportedDialog.vue"
import Logo from "@/components/Logo.vue"
import isWebview from "is-ua-webview"
import NewDialog from "./components/NewDialog.vue"
import SignInDialog from "@/components/SignInDialog.vue"

export default {
  name: "App",

  metaInfo: {
    htmlAttrs: {
      lang: "en-US",
    },
  },

  components: {
    AutoSnackbar,
    AuthUserMenu,
    SignInNotSupportedDialog,
    NewDialog,
    Logo,
    SignInDialog,
  },

  data: () => ({
    mounted: false,
    loaded: false,
    scrollY: 0,
    webviewDialog: false,
    signInDialog: false,
  }),

  computed: {
    ...mapState([
      "authUser",
      "error",
      "info",
      "newDialogOptions",
      "darkMode",
    ]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
    showHeader() {
      return (
        this.$route.name !== "landing" &&
        this.$route.name !== "auth" &&
        this.$route.name !== "sign-in" &&
        this.$route.name !== "sign-up" &&
        this.$route.name !== "privacy-policy"
      )
    },
    showFeedbackBtn() {
      return !this.isPhone || this.$route.name === "home"
    },
    routerViewClass() {
      let c = ""
      if (this.showHeader) {
        if (this.isPhone) {
          c += "tw-pt-12 "
        } else {
          c += "tw-pt-14 "
        }
      }
      return c
    },
  },

  methods: {
    ...mapMutations([
      "setAuthUser",
      "setSignUpFormEnabled",
      "setPricingPageConversion",
      "setDarkMode",
      "setAllowRegistration",
    ]),
    _applyDarkMode(dark) {
      if (dark) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
      this.$vuetify.theme.dark = dark
    },
    ...mapActions([
      "getEvents",
      "openNewDialog",
    ]),
    handleScroll(e) {
      this.scrollY = window.scrollY
    },
    _createNew(eventOnly = false) {
      this.openNewDialog({ eventOnly })
    },
    signIn() {
      if (
        this.$route.name === "event" ||
        this.$route.name === "group" ||
        this.$route.name === "signUp"
      ) {
        if (isWebview(navigator.userAgent)) {
          this.webviewDialog = true
          return
        }
        this.$router.push({ name: "sign-in" })
        // this.signInDialog = true
      } else {
        this.$router.push({ name: "sign-in" })
      }
    },
    _signIn(calendarType) {
      if (
        this.$route.name === "event" ||
        this.$route.name === "group" ||
        this.$route.name === "signUp"
      ) {
        let state
        if (this.$route.name === "event") {
          state = {
            eventId: this.$route.params.eventId,
            type: authTypes.EVENT_SIGN_IN,
          }
        } else if (this.$route.name === "group") {
          state = {
            groupId: this.$route.params.groupId,
            type: authTypes.GROUP_SIGN_IN,
          }
        }
        if (calendarType === calendarTypes.GOOGLE) {
          signInGoogle({
            state,
            selectAccount: true,
          })
        } else if (calendarType === calendarTypes.OUTLOOK) {
          signInOutlook({
            state,
            selectAccount: true,
          })
        }
      }
    },
    _emailSignIn(user) {
      this.setAuthUser(user)
      if (this.$route.name === "landing") {
        this.$router.push({ name: "home" })
      }
    },
  },

  async created() {
    this._applyDarkMode(this.darkMode)

    // Follow system preference when the user hasn't set one explicitly
    this._systemDarkModeQuery = window.matchMedia('(prefers-color-scheme: dark)')
    this._systemDarkModeHandler = (e) => {
      if (localStorage.getItem('darkMode') === null) {
        this.setDarkMode(e.matches)
      }
    }
    this._systemDarkModeQuery.addEventListener('change', this._systemDarkModeHandler)

    get("/status").then(({ allowRegistration }) => {
      this.setAllowRegistration(allowRegistration)
    }).catch(() => {})

    await get("/user/profile")
      .then((authUser) => {
        this.setAuthUser(authUser)
      })
      .catch(() => {
        this.setAuthUser(null)
      })
      .finally(() => {
        this.loaded = true
      })

    // Event listeners
    window.addEventListener("scroll", this.handleScroll)

    this.getEvents()
  },

  mounted() {
    this.mounted = true
    this.scrollY = window.scrollY
  },

  beforeDestroy() {
    window.removeEventListener("scroll", this.handleScroll)
    this._systemDarkModeQuery?.removeEventListener('change', this._systemDarkModeHandler)
  },

  watch: {
    darkMode(val) {
      this._applyDarkMode(val)
    },
    $route: {
      immediate: true,
      async handler() {
        const originalHref = window.location.href
        // Check for poster query parameter
        if (this.$route.query.p) {
          let location = null
          try {
            location = await getLocation()
          } catch (e) {
            // User probably has adblocker
          }

          post("/analytics/scanned-poster", {
            url: originalHref,
            location,
          })
        }
      },
    },
  },
}
</script>
