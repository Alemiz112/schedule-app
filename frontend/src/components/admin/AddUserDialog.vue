<template>
  <v-dialog :value="value" @input="$emit('input', $event)" max-width="400" persistent>
    <v-card>
      <v-card-title>Add user</v-card-title>
      <v-card-text>
        <v-form ref="form" v-model="formValid" lazy-validation>
          <v-text-field
            v-model="email"
            label="Email"
            type="email"
            outlined
            dense
            :rules="[(v) => !!v || 'Email is required']"
            class="tw-mb-1"
          />
          <v-text-field
            v-model="firstName"
            label="First name (optional)"
            outlined
            dense
            class="tw-mb-1"
          />
          <v-text-field
            v-model="lastName"
            label="Last name (optional)"
            outlined
            dense
            hide-details
          />
          <div v-if="errorMsg" class="tw-mt-3 tw-text-sm tw-text-red">
            {{ errorMsg }}
          </div>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="cancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="loading"
          :disabled="!formValid"
          @click="submit"
        >
          Create
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { post } from "@/utils"

export default {
  name: "AddUserDialog",

  props: {
    value: { type: Boolean, required: true },
  },

  data: () => ({
    email: "",
    firstName: "",
    lastName: "",
    formValid: true,
    loading: false,
    errorMsg: "",
  }),

  methods: {
    cancel() {
      this.reset()
      this.$emit("input", false)
    },
    reset() {
      this.email = ""
      this.firstName = ""
      this.lastName = ""
      this.errorMsg = ""
      this.$refs.form?.resetValidation()
    },
    async submit() {
      if (!this.$refs.form.validate()) return
      this.loading = true
      this.errorMsg = ""
      try {
        const user = await post("/admin/users", {
          email: this.email,
          firstName: this.firstName,
          lastName: this.lastName,
        })
        this.$emit("created", user)
        this.reset()
        this.$emit("input", false)
      } catch (err) {
        if (err.status === 409) {
          this.errorMsg = "A user with this email already exists."
        } else {
          this.errorMsg = "Something went wrong. Please try again."
        }
      } finally {
        this.loading = false
      }
    },
  },
}
</script>
