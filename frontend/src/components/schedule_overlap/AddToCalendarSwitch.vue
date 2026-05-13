<template>
  <div>
    <div class="tw-mb-1 tw-text-sm tw-text-black">Add to calendar</div>
    <div class="tw-mb-2 tw-text-xs tw-text-dark-gray">
      Automatically add scheduled events to your calendar
    </div>

    <div v-if="writableAccounts.length === 0" class="tw-text-xs tw-text-dark-gray">
      Connect a {{ outlookEnabled ? 'Google or Outlook' : 'Google' }} calendar to use this feature.
    </div>

    <template v-else>
      <v-switch
        :input-value="addToCalendar"
        @change="handleToggle"
        inset
        hide-details
        class="tw-mt-0"
        label="Enabled"
      />

      <template v-if="addToCalendar">
        <!-- Account selector -->
        <v-select
          :items="writableAccounts"
          item-text="label"
          item-value="key"
          :value="selectedAccountKey"
          @input="selectAccount"
          label="Account"
          outlined
          dense
          hide-details
          class="tw-mt-3"
        />

        <!-- Sub-calendar selector — shown when the selected account has sub-calendars -->
        <v-select
          v-if="subCalendarItems.length > 0"
          :items="subCalendarItems"
          item-text="name"
          item-value="id"
          :value="defaultCalendarId || subCalendarItems[0].id"
          @input="selectCalendar"
          label="Calendar"
          outlined
          dense
          hide-details
          class="tw-mt-2"
        />
      </template>
    </template>
  </div>
</template>

<script>
import { patch } from "@/utils"
import { mapState } from "vuex"

export default {
  name: "AddToCalendarSwitch",

  props: {
    addToCalendar: { type: Boolean, default: false },
    defaultCalendarKey: { type: String, default: null },
    defaultCalendarId: { type: String, default: null },
    calendarAccounts: { type: Object, default: () => ({}) },
  },

  computed: {
    ...mapState(["outlookEnabled"]),
    writableAccounts() {
      return Object.entries(this.calendarAccounts || {})
        .filter(([, acc]) => acc.calendarType === "google" || (this.outlookEnabled && acc.calendarType === "outlook"))
        .map(([key, acc]) => ({
          key,
          label: `${acc.email} (${acc.calendarType === "google" ? "Google" : "Outlook"})`,
        }))
    },

    selectedAccountKey() {
      if (this.defaultCalendarKey) return this.defaultCalendarKey
      return this.writableAccounts[0]?.key ?? null
    },

    selectedAccount() {
      if (!this.selectedAccountKey) return null
      return this.calendarAccounts[this.selectedAccountKey] ?? null
    },

    subCalendarItems() {
      const subs = this.selectedAccount?.subCalendars
      if (!subs) return []
      return Object.entries(subs)
        .filter(([, sub]) => sub.enabled !== false)
        .map(([id, sub]) => ({ id, name: sub.name || id }))
        .sort((a, b) => a.name.localeCompare(b.name))
    },
  },

  methods: {
    handleToggle(enabled) {
      patch("/user/calendar-options", { addToCalendar: enabled })
      this.$emit("update:addToCalendar", enabled)
    },

    selectAccount(key) {
      this.$emit("update:defaultCalendarKey", key)
      // Reset calendar ID when account changes
      this.$emit("update:defaultCalendarId", null)
      patch("/user/calendar-options", { defaultCalendarKey: key, defaultCalendarId: "" })
    },

    selectCalendar(id) {
      this.$emit("update:defaultCalendarId", id)
      patch("/user/calendar-options", { defaultCalendarId: id })
    },
  },
}
</script>
