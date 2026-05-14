<template>
  <div
    v-if="showGradient"
    class="tw-pointer-events-none tw-absolute tw-bottom-0 tw-left-0 tw-right-0 tw-z-20 tw-flex tw-h-16 tw-items-end tw-justify-center"
    :style="{ background: gradientStyle }"
  >
    <v-btn
      v-if="showArrow"
      fab
      x-small
      class="tw-pointer-events-auto tw-transform"
      @click="scrollToBottom"
    >
      <v-icon>mdi-chevron-down</v-icon>
    </v-btn>
  </div>
</template>

<script>
import { mapState } from "vuex"

export default {
  name: "OverflowGradient",
  props: {
    scrollContainer: {
      type: HTMLElement,
      required: true,
    },
    showArrow: {
      type: Boolean,
      default: true,
    },
  },
  computed: {
    ...mapState(["darkMode"]),
    gradientStyle() {
      const color = this.darkMode ? "44,44,44" : "255,255,255"
      return `linear-gradient(to bottom, rgba(${color},0) 0%, rgba(${color},1) 100%)`
    },
  },
  data() {
    return {
      showGradient: false,
    }
  },
  mounted() {
    this.checkScroll()
    this.scrollContainer.addEventListener("scroll", this.checkScroll)
  },
  beforeDestroy() {
    this.scrollContainer.removeEventListener("scroll", this.checkScroll)
  },
  methods: {
    /**
     * Checks if the scroll bar is scrolled to the bottom of the client
     */
    checkScroll() {
      const { scrollHeight, clientHeight, scrollTop } = this.scrollContainer
      this.showGradient =
        scrollHeight > clientHeight && scrollTop < scrollHeight - clientHeight - 1 // 1px tolerance
    },
    scrollToBottom() {
      this.scrollContainer.scrollTop = this.scrollContainer.scrollHeight
    },
  },
}
</script>
