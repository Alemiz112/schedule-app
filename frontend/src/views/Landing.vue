<template>
  <div>
    <!-- Invite-only banner -->
    <div
      v-if="!allowRegistration"
      class="tw-flex tw-items-center tw-justify-center tw-gap-2 tw-bg-very-dark-gray tw-px-4 tw-py-2.5 tw-text-sm tw-text-white"
    >
      <v-icon small color="white">mdi-lock-outline</v-icon>
      <span>This instance is <strong>invite only</strong> — registration is currently closed.</span>
    </div>

    <!-- Navbar -->
    <div class="tw-bg-white tw-px-4">
      <div class="tw-m-auto tw-flex tw-max-w-6xl tw-items-center tw-py-5">
        <Logo type="timeful" />
        <v-spacer />
        <LandingPageHeader>
          <div v-if="authUser" class="tw-ml-2">
            <AuthUserMenu />
          </div>
          <v-btn v-else text :to="{ name: 'sign-in' }">Sign in</v-btn>
        </LandingPageHeader>
      </div>
    </div>

    <!-- Hero -->
    <div class="hero-section tw-bg-white tw-px-4 tw-py-24 sm:tw-py-36">
      <div class="tw-m-auto tw-max-w-2xl tw-text-center">
        <h1
          class="tw-mb-5 tw-text-4xl tw-font-semibold tw-tracking-tight tw-text-black sm:tw-text-5xl xl:tw-text-6xl"
        >
          Find a time that works for everyone
        </h1>
        <p class="tw-mb-10 tw-text-lg tw-text-dark-gray">
          Stop the back-and-forth. Timeful shows you exactly when your whole
          group is free.
        </p>

        <div v-if="authUser">
          <v-btn
            large
            color="primary"
            dark
            class="tw-rounded-lg tw-px-10"
            @click="openDashboard"
          >
            Open dashboard
          </v-btn>
        </div>
        <div
          v-else
          class="tw-mx-auto tw-flex tw-max-w-sm tw-flex-col tw-gap-3"
        >
          <v-btn
            large
            color="primary"
            dark
            class="tw-w-full tw-rounded-lg"
            :to="{ name: 'sign-in' }"
          >
            Get started — it's free
          </v-btn>
          <div class="tw-flex tw-gap-2">
            <v-text-field
              v-model="eventCode"
              placeholder="Have an event code?"
              outlined
              dense
              hide-details
              background-color="white"
              @keyup.enter="openEventCode"
            />
            <v-btn
              outlined
              color="primary"
              height="40"
              :disabled="!eventCode.trim()"
              @click="openEventCode"
            >
              <v-icon>mdi-arrow-right</v-icon>
            </v-btn>
          </div>
        </div>

        <!-- Calendar logos -->
        <div
          class="tw-mt-10 tw-flex tw-flex-wrap tw-items-center tw-justify-center tw-gap-3"
        >
          <span class="tw-text-sm tw-text-dark-gray">Works with</span>
          <div
            v-for="cal in calendars"
            :key="cal.label"
            class="tw-flex tw-items-center tw-gap-1.5 tw-rounded-full tw-border tw-border-light-gray-stroke tw-bg-off-white tw-px-3 tw-py-1"
          >
            <img
              :src="cal.src"
              :alt="cal.label"
              class="tw-h-4 tw-w-4 tw-object-contain"
              :style="cal.invert && darkMode ? { filter: 'invert(1)' } : {}"
            />
            <span class="tw-text-xs tw-text-very-dark-gray">{{ cal.label }}</span>
          </div>
          <div
            class="tw-flex tw-items-center tw-gap-1.5 tw-rounded-full tw-border tw-border-light-gray-stroke tw-bg-off-white tw-px-3 tw-py-1"
          >
            <v-icon small color="#6B6B6B">mdi-calendar-text</v-icon>
            <span class="tw-text-xs tw-text-very-dark-gray">ICS feeds</span>
          </div>
        </div>
      </div>
    </div>

    <!-- How it works -->
    <div class="tw-bg-green tw-px-4 tw-py-16">
      <div class="tw-m-auto tw-max-w-4xl">
        <h2
          class="tw-mb-12 tw-text-center tw-text-2xl tw-font-semibold tw-text-white sm:tw-text-3xl"
        >
          How it works
        </h2>
        <div class="tw-grid tw-grid-cols-1 tw-gap-10 sm:tw-grid-cols-3">
          <div
            v-for="(step, i) in steps"
            :key="i"
            class="tw-flex tw-flex-col tw-items-center tw-text-center"
          >
            <div
              class="tw-mb-4 tw-flex tw-h-11 tw-w-11 tw-items-center tw-justify-center tw-rounded-full tw-bg-white"
            >
              <v-icon color="green">{{ step.icon }}</v-icon>
            </div>
            <p class="tw-mb-1 tw-font-semibold tw-text-white">{{ step.title }}</p>
            <p class="tw-text-sm tw-text-pale-green">{{ step.desc }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Features -->
    <div class="tw-bg-off-white tw-px-4 tw-py-16">
      <div class="tw-m-auto tw-max-w-4xl">
        <h2
          class="tw-mb-10 tw-text-center tw-text-2xl tw-font-semibold tw-text-black sm:tw-text-3xl"
        >
          Everything your team needs
        </h2>
        <div class="tw-grid tw-grid-cols-1 tw-gap-5 sm:tw-grid-cols-2 lg:tw-grid-cols-3">
          <div
            v-for="feature in features"
            :key="feature.title"
            class="tw-rounded-xl tw-bg-white tw-p-6"
          >
            <div
              class="tw-mb-3 tw-flex tw-h-10 tw-w-10 tw-items-center tw-justify-center tw-rounded-lg tw-bg-ligher-green"
            >
              <v-icon color="green">{{ feature.icon }}</v-icon>
            </div>
            <h3 class="tw-mb-1.5 tw-font-semibold tw-text-black">{{ feature.title }}</h3>
            <p class="tw-text-sm tw-text-dark-gray">{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom CTA -->
    <div class="tw-bg-white tw-px-4 tw-py-20">
      <div class="tw-m-auto tw-max-w-lg tw-text-center">
        <h2
          class="tw-mb-3 tw-text-2xl tw-font-semibold tw-text-black sm:tw-text-3xl"
        >
          Ready to get started?
        </h2>
        <p class="tw-mb-8 tw-text-dark-gray">Free to use. No credit card required.</p>
        <v-btn
          v-if="authUser"
          large
          color="primary"
          dark
          class="tw-rounded-lg tw-px-10"
          @click="openDashboard"
        >
          Go to dashboard
        </v-btn>
        <v-btn
          v-else
          large
          color="primary"
          dark
          class="tw-rounded-lg tw-px-10"
          :to="{ name: 'sign-in' }"
        >
          Get started — it's free
        </v-btn>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hero-section {
  background-image: radial-gradient(
    ellipse 100% 80% at 50% -5%,
    rgba(0, 153, 76, 0.07) 0%,
    transparent 70%
  );
}
</style>

<script>
import LandingPageHeader from "@/components/landing/LandingPageHeader.vue"
import Logo from "@/components/Logo.vue"
import AuthUserMenu from "@/components/AuthUserMenu.vue"
import { mapState } from "vuex"

export default {
  name: "Landing",

  metaInfo: {
    title: "Timeful - Find a time to meet",
  },

  components: {
    LandingPageHeader,
    Logo,
    AuthUserMenu,
  },

  data: () => ({
    eventCode: "",
    calendars: [
      { label: "Google Calendar", src: require("@/assets/gcal_logo.png"), invert: false },
      { label: "Outlook", src: require("@/assets/outlook_logo.svg"), invert: false },
      { label: "Apple Calendar", src: require("@/assets/apple_logo.svg"), invert: true },
    ],
    steps: [
      {
        icon: "mdi-calendar-plus",
        title: "Create an event",
        desc: "Set your date range and time preferences in seconds.",
      },
      {
        icon: "mdi-share-variant",
        title: "Share the link",
        desc: "Invite everyone with a single link — no sign-up required.",
      },
      {
        icon: "mdi-check-circle-outline",
        title: "Find the perfect time",
        desc: "See exactly when your whole group is available.",
      },
    ],
    features: [
      {
        icon: "mdi-calendar-sync",
        title: "Calendar sync",
        desc: "Autofill your availability from Google, Outlook, Apple Calendar, or any ICS feed.",
      },
      {
        icon: "mdi-account-group",
        title: "No account needed",
        desc: "Anyone can respond as a guest. Only the organizer needs to sign up.",
      },
      {
        icon: "mdi-shield-check",
        title: "Privacy first",
        desc: "We never store your calendar data — it's only fetched to prefill your availability.",
      },
      {
        icon: "mdi-cellphone",
        title: "Mobile friendly",
        desc: "Works great on any device — phone, tablet, or desktop.",
      },
      {
        icon: "mdi-earth",
        title: "Timezone aware",
        desc: "All times are automatically shown in each person's local timezone.",
      },
      {
        icon: "mdi-bell-outline",
        title: "Email reminders",
        desc: "Send reminders and get notified once everyone has responded.",
      },
    ],
  }),

  computed: {
    ...mapState(["authUser", "darkMode", "allowRegistration"]),
  },

  methods: {
    openEventCode() {
      const code = this.eventCode.trim()
      if (!code) return
      this.$router.push({ name: "event", params: { eventId: code } })
    },
    openDashboard() {
      this.$router.push({ name: "home" })
    },
  },
}
</script>
