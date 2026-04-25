<template>
  <div class="tw-select-none tw-py-4">
    <!-- ── Success ──────────────────────────────────────── -->
    <div
      v-if="submitted"
      class="tw-flex tw-flex-col tw-items-center tw-py-12 tw-text-center"
    >
      <v-icon color="primary" style="font-size: 3rem" class="tw-mb-4"
        >mdi-check-circle</v-icon
      >
      <div class="tw-mb-2 tw-text-xl tw-font-medium tw-text-black">
        Appointment requested!
      </div>
      <div class="tw-text-dark-gray">
        Your request has been submitted. The organiser will confirm your
        appointment.
      </div>
    </div>

    <!-- ── Step 2: details form ────────────────────────── -->
    <template v-else-if="selectedSlot">
      <div class="tw-mb-4 tw-flex tw-items-center tw-gap-2">
        <v-btn icon small @click="selectedSlot = null">
          <v-icon>mdi-arrow-left</v-icon>
        </v-btn>
        <div class="tw-text-xl tw-font-medium tw-text-black">Your details</div>
      </div>

      <div
        class="tw-mb-6 tw-flex tw-items-center tw-gap-2 tw-rounded-lg tw-border tw-border-solid tw-border-light-gray tw-px-4 tw-py-3"
      >
        <v-icon color="primary" small>mdi-calendar-clock</v-icon>
        <div class="tw-font-medium tw-text-black">
          {{ formatFullSlot(selectedSlot) }}
        </div>
      </div>

      <v-form ref="form" v-model="formValid" lazy-validation>
        <v-text-field
          v-model="name"
          label="Your name"
          outlined
          dense
          :rules="[(v) => !!v || 'Name is required']"
          class="tw-mb-1"
        />
        <v-text-field
          v-model="email"
          label="Email (optional)"
          outlined
          dense
          class="tw-mb-1"
        />
        <v-textarea
          v-model="notes"
          label="Notes (optional)"
          outlined
          dense
          rows="2"
          hide-details
          class="tw-mb-4"
        />
        <v-btn
          block
          color="primary"
          class="tw-bg-green tw-text-white"
          :loading="loading"
          :disabled="!formValid"
          @click="submit"
        >
          Request appointment
        </v-btn>
        <div v-if="errorMsg" class="tw-mt-2 tw-text-center tw-text-sm tw-text-red">
          {{ errorMsg }}
        </div>
      </v-form>
    </template>

    <!-- ── Step 1: time grid ───────────────────────────── -->
    <template v-else>
      <!-- Navigation -->
      <div class="tw-mb-2 tw-flex tw-items-center tw-justify-between tw-px-4">
        <v-btn
          outlined
          icon
          class="tw-border-gray"
          :disabled="!hasPrevPage"
          @click="prevPage"
        >
          <v-icon>mdi-chevron-left</v-icon>
        </v-btn>
        <div class="tw-text-lg tw-font-medium tw-capitalize">
          {{ pageLabel }}
        </div>
        <v-btn
          outlined
          icon
          class="tw-border-gray"
          :disabled="!hasNextPage"
          @click="nextPage"
        >
          <v-icon>mdi-chevron-right</v-icon>
        </v-btn>
      </div>

      <!-- Grid -->
      <div class="tw-flex tw-pl-4 tw-pr-4">
        <!-- Time labels -->
        <div class="tw-w-8 tw-flex-none sm:tw-w-12">
          <!-- header spacer -->
          <div class="tw-h-14"></div>
          <!-- labels -->
          <div
            v-for="(row, r) in timeRows"
            :key="r"
            :style="{ height: TIMESLOT_HEIGHT + 'px' }"
            class="-tw-ml-1 tw-flex tw-items-start tw-justify-end tw-pr-1 tw-text-right tw-text-xs tw-font-light tw-uppercase tw-text-dark-gray sm:tw-pr-2"
          >
            {{ row.label }}
          </div>
        </div>

        <!-- Day columns -->
        <div class="tw-flex tw-flex-1">
          <div
            v-for="(col, c) in columns"
            :key="c"
            class="tw-flex-1 tw-min-w-0"
          >
            <!-- Day header -->
            <div class="tw-flex tw-h-14 tw-flex-col tw-items-center tw-justify-center tw-bg-white">
              <div class="tw-text-[12px] tw-font-light tw-capitalize tw-text-very-dark-gray">
                {{ col.dateString }}
              </div>
              <div class="tw-text-base tw-capitalize sm:tw-text-lg">
                {{ col.dayText }}
              </div>
            </div>

            <!-- Cells -->
            <div
              v-for="(row, r) in timeRows"
              :key="r"
              class="timeslot tw-w-full tw-border-b tw-border-r tw-border-white"
              :style="{ height: TIMESLOT_HEIGHT + 'px' }"
              :class="cellClass(col, row)"
              @click="!row.isPast && !isSlotBooked(col, row) && pickSlot(col, row)"
            ></div>
          </div>
        </div>
      </div>

      <div class="tw-mt-3 tw-px-4 tw-text-xs tw-text-dark-gray">
        Click a slot to book it
      </div>
    </template>
  </div>
</template>

<script>
import { get, post, dateToDowDate } from "@/utils"
import { eventTypes } from "@/constants"
import dayjs from "dayjs"

const HOUR_HEIGHT = 60

export default {
  name: "AppointmentBookingPanel",

  props: {
    event: { type: Object, required: true },
  },

  data: () => ({
    page: 0,       // weekOffset for DOW, page index for specific dates
    bookedSlots: [], // [{startDate, endDate}] — pending or approved bookings
    selectedSlot: null,
    name: "",
    email: "",
    notes: "",
    formValid: true,
    loading: false,
    submitted: false,
    errorMsg: "",
  }),

  mounted() {
    this.fetchBookedSlots()
  },

  computed: {
    isWeekly() {
      return this.event.type === eventTypes.DOW
    },
    TIMESLOT_HEIGHT() {
      const inc = this.event.timeIncrement || 30
      return Math.floor((HOUR_HEIGHT * inc) / 60)
    },
    /** One column per event day, mapped to actual dates for the current page/week */
    columns() {
      const cols = this.event.dates.map((refDate) => {
        let actualDate
        if (this.isWeekly) {
          actualDate = dateToDowDate(
            this.event.dates,
            new Date(refDate),
            this.page,
            true
          )
        } else {
          actualDate = new Date(refDate)
        }
        return {
          dayText: dayjs(actualDate).format("ddd"),
          dateString: dayjs(actualDate).format("MMM D"),
          actualDate,
        }
      })
      cols.sort((a, b) => a.actualDate - b.actualDate)

      if (!this.isWeekly) {
        const now = new Date()
        const filtered = cols.filter((c) => c.actualDate > now)
        return filtered.slice(this.page * 7, this.page * 7 + 7)
      }
      return cols
    },
    /** One row per time slot within the availability window */
    timeRows() {
      if (this.columns.length === 0) return []
      const base = new Date(this.columns[0].actualDate)
      const now = new Date()
      const inc = this.event.timeIncrement || 30
      const windowMinutes = this.event.duration * 60
      const rows = []
      for (let min = 0; min < windowMinutes; min += inc) {
        const slotStart = new Date(base)
        slotStart.setUTCMinutes(slotStart.getUTCMinutes() + min)
        rows.push({
          minuteOffset: min,
          label: min % 60 === 0 ? dayjs(slotStart).format("h a") : "",
          isPast: slotStart <= now,
        })
      }
      return rows
    },
    hasPrevPage() {
      return this.page > 0
    },
    hasNextPage() {
      if (this.isWeekly) return true
      const now = new Date()
      const total = this.event.dates.filter((d) => new Date(d) > now).length
      return (this.page + 1) * 7 < total
    },
    pageLabel() {
      if (this.columns.length === 0) return ""
      if (this.isWeekly) {
        const first = this.columns[0].actualDate
        const last = this.columns[this.columns.length - 1].actualDate
        return `${dayjs(first).format("MMM D")} – ${dayjs(last).format("MMM D")}`
      }
      return dayjs(this.columns[0].actualDate).format("MMMM YYYY")
    },
  },

  methods: {
    async fetchBookedSlots() {
      try {
        this.bookedSlots = await get(
          `/events/${this.event._id}/appointment-requests/booked`
        )
      } catch {
        // Non-fatal — worst case a taken slot stays green until they try to submit
      }
    },

    isSlotBooked(col, row) {
      const { start } = this.slotDates(col, row)
      return this.bookedSlots.some(
        (b) => new Date(b.startDate).getTime() === start.getTime()
      )
    },

    cellClass(col, row) {
      if (row.isPast || this.isSlotBooked(col, row)) {
        return "tw-cursor-default tw-bg-off-white"
      }
      return "tw-cursor-pointer tw-bg-ligher-green hover:tw-bg-green hover:tw-opacity-70 tw-transition-colors tw-duration-75"
    },

    prevPage() {
      if (this.hasPrevPage) this.page--
    },
    nextPage() {
      if (this.hasNextPage) this.page++
    },

    slotDates(col, row) {
      const start = new Date(col.actualDate)
      start.setUTCMinutes(start.getUTCMinutes() + row.minuteOffset)
      const end = new Date(start)
      end.setUTCMinutes(end.getUTCMinutes() + (this.event.timeIncrement || 30))
      return { start, end }
    },

    pickSlot(col, row) {
      const { start, end } = this.slotDates(col, row)
      this.selectedSlot = { start, end }
      this.errorMsg = ""
    },

    formatFullSlot(slot) {
      return `${dayjs(slot.start).format("ddd, MMM D")}  ·  ${dayjs(slot.start).format("h:mm")}–${dayjs(slot.end).format("h:mm A")}`
    },

    async submit() {
      if (!this.$refs.form.validate()) return
      this.loading = true
      this.errorMsg = ""
      try {
        await post(`/events/${this.event._id}/appointment-requests`, {
          startDate: this.selectedSlot.start,
          endDate: this.selectedSlot.end,
          name: this.name,
          email: this.email,
          notes: this.notes,
        })
        this.submitted = true
      } catch {
        this.errorMsg = "Something went wrong. Please try again."
      } finally {
        this.loading = false
      }
    },
  },
}
</script>
