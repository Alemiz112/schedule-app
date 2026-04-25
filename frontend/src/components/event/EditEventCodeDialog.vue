<template>
  <v-dialog :value="value" @input="$emit('input', $event)" max-width="440" persistent>
    <v-card>
      <v-card-title>Custom event slug</v-card-title>
      <v-card-text>
        <!-- Current link preview -->
        <div class="tw-mb-4 tw-rounded-lg tw-bg-off-white tw-px-3 tw-py-2 tw-text-sm">
          <div class="tw-mb-0.5 tw-text-xs tw-text-dark-gray">Shared link</div>
          <span class="tw-font-mono tw-break-all tw-text-black">
            {{ origin }}/e/<span class="tw-font-semibold tw-text-green">{{ previewCode }}</span>
          </span>
        </div>

        <v-text-field
          v-model="slug"
          label="Custom slug (optional)"
          placeholder="e.g. team-meeting"
          outlined
          dense
          clearable
          :error-messages="errorMsg"
          :success-messages="successMsg"
          hint="3–30 characters · letters, numbers and hyphens · cannot start or end with a hyphen"
          persistent-hint
          @input="onInput"
          @keyup.enter="save"
        />

        <div
          v-if="event.customSlug && !slug"
          class="tw-mt-3 tw-text-xs tw-text-dark-gray"
        >
          Saving with an empty slug will remove the custom slug.
          The auto-generated code
          <span class="tw-font-mono tw-font-medium tw-text-black">{{ event.shortId }}</span>
          will still work.
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="cancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="saving"
          :disabled="!!errorMsg || saving"
          @click="save"
        >
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { put } from "@/utils"

const VALID_RE = /^[a-z0-9][a-z0-9-]{1,28}[a-z0-9]$/

export default {
  name: "EditEventCodeDialog",

  props: {
    value: { type: Boolean, required: true },
    event: { type: Object, required: true },
  },

  data() {
    return {
      slug: this.event.customSlug ?? "",
      errorMsg: "",
      successMsg: "",
      saving: false,
      debounceTimer: null,
    }
  },

  computed: {
    origin() {
      return window.location.origin
    },
    previewCode() {
      const s = this.slug?.trim().toLowerCase()
      if (s) return s
      return this.event.customSlug || this.event.shortId || this.event._id
    },
  },

  watch: {
    value(open) {
      if (open) {
        this.slug = this.event.customSlug ?? ""
        this.errorMsg = ""
        this.successMsg = ""
      }
    },
  },

  methods: {
    onInput() {
      this.errorMsg = ""
      this.successMsg = ""
      clearTimeout(this.debounceTimer)

      const normalized = (this.slug ?? "").trim().toLowerCase()
      if (!normalized) return  // empty = clear slug, that's valid

      if (!VALID_RE.test(normalized)) {
        this.errorMsg =
          "3–30 characters: letters, numbers, and hyphens only. Cannot start or end with a hyphen."
        return
      }

      this.debounceTimer = setTimeout(() => {
        this.successMsg = "Looks good!"
      }, 300)
    },

    cancel() {
      this.$emit("input", false)
    },

    async save() {
      if (this.errorMsg) return
      const normalized = (this.slug ?? "").trim().toLowerCase()

      this.saving = true
      this.errorMsg = ""
      this.successMsg = ""
      try {
        const result = await put(
          `/events/${this.event._id}/custom-slug`,
          { slug: normalized }
        )
        this.$emit("updated", result.customSlug ?? null)
        this.$emit("input", false)
      } catch (err) {
        const code = err?.parsed?.error
        if (code === "short-id-taken") {
          this.errorMsg = "This slug is already taken. Please choose another."
        } else if (code === "short-id-invalid") {
          this.errorMsg =
            "Invalid slug. Use 3–30 characters: letters, numbers, and hyphens only."
        } else {
          this.errorMsg = "Something went wrong. Please try again."
        }
      } finally {
        this.saving = false
      }
    },
  },
}
</script>
