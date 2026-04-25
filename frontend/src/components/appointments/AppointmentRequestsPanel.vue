<template>
  <div class="tw-max-w-lg">
    <div class="tw-mb-4 tw-flex tw-items-center tw-justify-between">
      <div class="tw-text-xl tw-font-medium tw-text-black">
        Appointment requests
      </div>
      <v-btn icon small :loading="refreshing" @click="fetchRequests">
        <v-icon small>mdi-refresh</v-icon>
      </v-btn>
    </div>

    <!-- Status filter chips -->
    <div class="tw-mb-4 tw-flex tw-gap-2">
      <v-chip
        v-for="f in filters"
        :key="f.value"
        :color="statusFilter === f.value ? 'primary' : ''"
        :text-color="statusFilter === f.value ? 'white' : ''"
        small
        @click="statusFilter = f.value"
      >
        {{ f.label }}
        <span v-if="f.value === 'pending' && pendingCount > 0" class="tw-ml-1">
          ({{ pendingCount }})
        </span>
      </v-chip>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="tw-flex tw-justify-center tw-py-8">
      <v-progress-circular indeterminate color="primary" />
    </div>

    <!-- Empty state -->
    <div
      v-else-if="filteredRequests.length === 0"
      class="tw-py-8 tw-text-center tw-text-dark-gray"
    >
      {{
        statusFilter === "all"
          ? "No appointment requests yet. Share the event link to start receiving bookings."
          : `No ${statusFilter} requests.`
      }}
    </div>

    <!-- Request cards -->
    <div v-else class="tw-flex tw-flex-col tw-gap-3">
      <v-card
        v-for="req in filteredRequests"
        :key="req._id"
        outlined
        class="tw-rounded-lg"
      >
        <v-card-text class="tw-pb-2">
          <!-- Time slot -->
          <div class="tw-mb-2 tw-flex tw-items-center tw-gap-1.5 tw-text-black">
            <v-icon small color="primary">mdi-calendar-clock</v-icon>
            <span class="tw-font-medium">{{ formatSlot(req) }}</span>
          </div>

          <!-- Guest name + email -->
          <div class="tw-flex tw-items-center tw-gap-1.5 tw-text-very-dark-gray">
            <v-icon small>mdi-account</v-icon>
            <span>{{ req.name }}</span>
            <span v-if="req.email" class="tw-text-dark-gray">·</span>
            <span v-if="req.email" class="tw-text-dark-gray">{{ req.email }}</span>
          </div>

          <!-- Notes -->
          <div
            v-if="req.notes"
            class="tw-mt-2 tw-rounded tw-bg-light-gray tw-px-3 tw-py-2 tw-text-sm tw-text-very-dark-gray"
          >
            {{ req.notes }}
          </div>

          <!-- Status badge for non-pending -->
          <v-chip
            v-if="req.status !== 'pending'"
            :color="req.status === 'approved' ? 'success' : 'error'"
            text-color="white"
            x-small
            class="tw-mt-2"
          >
            {{ req.status === "approved" ? "Approved" : "Rejected" }}
          </v-chip>
        </v-card-text>

        <!-- Actions for pending requests -->
        <v-card-actions v-if="req.status === 'pending'" class="tw-pt-0">
          <v-btn
            text
            small
            color="error"
            :loading="loadingId === req._id + '-reject'"
            @click="reject(req)"
          >
            Decline
          </v-btn>
          <v-spacer />
          <v-btn
            small
            color="primary"
            class="tw-bg-green tw-text-white"
            :loading="loadingId === req._id + '-approve'"
            @click="approve(req)"
          >
            Approve &amp; schedule
          </v-btn>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>

<script>
import { get, post } from "@/utils"
import { mapState } from "vuex"
import dayjs from "dayjs"

export default {
  name: "AppointmentRequestsPanel",

  props: {
    event: { type: Object, required: true },
  },

  data: () => ({
    requests: [],
    loading: true,
    refreshing: false,
    loadingId: "",
    statusFilter: "pending",
    filters: [
      { label: "Pending", value: "pending" },
      { label: "Approved", value: "approved" },
      { label: "Declined", value: "rejected" },
      { label: "All", value: "all" },
    ],
  }),

  computed: {
    ...mapState(["authUser"]),
    filteredRequests() {
      if (this.statusFilter === "all") return this.requests
      return this.requests.filter((r) => r.status === this.statusFilter)
    },
    pendingCount() {
      return this.requests.filter((r) => r.status === "pending").length
    },
  },

  mounted() {
    this.fetchRequests()
  },

  methods: {
    async fetchRequests(showRefreshing = false) {
      if (showRefreshing) this.refreshing = true
      else this.loading = true
      try {
        this.requests = await get(
          `/events/${this.event._id}/appointment-requests`
        )
      } catch {
        // Silently fail — user will see empty state
      } finally {
        this.loading = false
        this.refreshing = false
      }
    },

    formatSlot(req) {
      const start = dayjs(req.startDate)
      const end = dayjs(req.endDate)
      return `${start.format("ddd, MMM D")}  ·  ${start.format("h:mm")}–${end.format("h:mm A")}`
    },

    async approve(req) {
      this.loadingId = req._id + "-approve"
      try {
        const response = await post(
          `/events/${this.event._id}/appointment-requests/${req._id}/approve`
        )
        this.updateRequest(response)
        if (!response.calendarEventCreated) {
          this.openCalendarLink(response)
        }
      } catch {
        // Keep the request as-is on error
      } finally {
        this.loadingId = ""
      }
    },

    async reject(req) {
      this.loadingId = req._id + "-reject"
      try {
        const updated = await post(
          `/events/${this.event._id}/appointment-requests/${req._id}/reject`
        )
        this.updateRequest(updated)
      } catch {
        // Keep the request as-is on error
      } finally {
        this.loadingId = ""
      }
    },

    updateRequest(updated) {
      const idx = this.requests.findIndex((r) => r._id === updated._id)
      if (idx !== -1) this.$set(this.requests, idx, updated)
    },

    openCalendarLink(req) {
      const start = new Date(req.startDate)
        .toISOString()
        .replace(/([-:]|\.000)/g, "")
      const end = new Date(req.endDate)
        .toISOString()
        .replace(/([-:]|\.000)/g, "")
      const eventId = this.event.shortId ?? this.event._id
      const title = `${this.event.name} with ${req.name}`
      const details = `Booked via Timeful: https://timeful.app/e/${eventId}`
      const url =
        `https://calendar.google.com/calendar/render?action=TEMPLATE` +
        `&text=${encodeURIComponent(title)}` +
        `&dates=${start}/${end}` +
        `&details=${encodeURIComponent(details)}` +
        (req.email ? `&add=${encodeURIComponent(req.email)}` : "")
      window.open(url, "_blank")
    },
  },
}
</script>
