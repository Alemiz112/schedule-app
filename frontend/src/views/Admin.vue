<template>
  <div class="tw-mx-auto tw-max-w-3xl tw-p-4 sm:tw-p-8">
    <div class="tw-mb-8 tw-text-2xl tw-font-medium tw-text-black">Admin Panel</div>

    <div class="tw-flex tw-flex-col tw-gap-10">
      <!-- ── Instance Settings ────────────────────────── -->
      <section>
        <div class="tw-mb-3 tw-text-lg tw-font-medium tw-text-black">
          Instance Settings
        </div>
        <v-card outlined class="tw-rounded-lg tw-px-4 tw-py-2">
          <v-switch
            v-model="allowRegistration"
            :loading="savingSettings"
            @change="saveSettings"
            hide-details
            color="primary"
          >
            <template v-slot:label>
              <span class="tw-text-sm tw-text-black">Allow new user registration</span>
            </template>
          </v-switch>
        </v-card>
      </section>

      <!-- ── Users ───────────────────────────────────── -->
      <section>
        <div class="tw-mb-3 tw-flex tw-items-center tw-justify-between">
          <div class="tw-text-lg tw-font-medium tw-text-black">Users</div>
          <v-btn small color="primary" class="tw-bg-green tw-text-white" @click="showAddUserDialog = true">
            <v-icon small class="tw-mr-1">mdi-plus</v-icon>
            Add user
          </v-btn>
        </div>

        <v-data-table
          :headers="headers"
          :items="usersWithName"
          :loading="loadingUsers"
          hide-default-footer
          disable-pagination
          class="tw-rounded-lg tw-border tw-border-solid tw-border-light-gray"
        >
          <template v-slot:item.role="{ item }">
            <v-select
              :value="item.role || 'user'"
              :items="roleOptions"
              dense
              hide-details
              class="tw-w-28 tw-text-sm"
              :disabled="item._id === authUser._id"
              @change="(role) => setRole(item, role)"
            />
          </template>
          <template v-slot:item.actions="{ item }">
            <v-btn
              icon
              small
              color="error"
              :disabled="item._id === authUser._id"
              @click="confirmDelete(item)"
            >
              <v-icon small>mdi-delete</v-icon>
            </v-btn>
          </template>
        </v-data-table>
      </section>
    </div>

    <!-- Add user dialog -->
    <AddUserDialog v-model="showAddUserDialog" @created="onUserCreated" />

    <!-- Delete confirmation dialog -->
    <v-dialog v-model="deleteDialog" max-width="380">
      <v-card>
        <v-card-title>Delete user?</v-card-title>
        <v-card-text v-if="userToDelete">
          Are you sure you want to delete
          <strong>{{ userToDelete.email }}</strong>? This cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="deleteDialog = false">Cancel</v-btn>
          <v-btn text color="error" :loading="deleting" @click="doDelete">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { get, patch, _delete } from "@/utils"
import { mapState } from "vuex"
import AddUserDialog from "@/components/admin/AddUserDialog.vue"

export default {
  name: "Admin",

  components: { AddUserDialog },

  data: () => ({
    allowRegistration: true,
    savingSettings: false,

    users: [],
    loadingUsers: false,

    showAddUserDialog: false,

    deleteDialog: false,
    userToDelete: null,
    deleting: false,

    headers: [
      { text: "Name", value: "displayName", sortable: false },
      { text: "Email", value: "email" },
      { text: "Role", value: "role", sortable: false },
      { text: "", value: "actions", sortable: false, align: "end" },
    ],
    roleOptions: [
      { text: "User", value: "user" },
      { text: "Admin", value: "admin" },
    ],
  }),

  computed: {
    ...mapState(["authUser"]),
    usersWithName() {
      return this.users.map((u) => ({
        ...u,
        displayName: [u.firstName, u.lastName].filter(Boolean).join(" ") || "—",
      }))
    },
  },

  created() {
    if (!this.authUser || this.authUser.role !== "admin") {
      this.$router.replace({ name: "home" })
      return
    }
    this.fetchSettings()
    this.fetchUsers()
  },

  methods: {
    async fetchSettings() {
      try {
        const s = await get("/admin/settings")
        this.allowRegistration = s.allowRegistration
      } catch {
        // silently ignore
      }
    },
    async saveSettings() {
      this.savingSettings = true
      try {
        await patch("/admin/settings", { allowRegistration: this.allowRegistration })
      } catch {
        // revert on failure
        this.allowRegistration = !this.allowRegistration
      } finally {
        this.savingSettings = false
      }
    },

    async fetchUsers() {
      this.loadingUsers = true
      try {
        this.users = await get("/admin/users")
      } finally {
        this.loadingUsers = false
      }
    },

    async setRole(user, role) {
      try {
        await patch(`/admin/users/${user._id}/role`, { role })
        const idx = this.users.findIndex((u) => u._id === user._id)
        if (idx !== -1) this.$set(this.users, idx, { ...this.users[idx], role })
      } catch {
        // revert — refetch to get accurate state
        this.fetchUsers()
      }
    },

    confirmDelete(user) {
      this.userToDelete = user
      this.deleteDialog = true
    },
    async doDelete() {
      if (!this.userToDelete) return
      this.deleting = true
      try {
        await _delete(`/admin/users/${this.userToDelete._id}`)
        this.users = this.users.filter((u) => u._id !== this.userToDelete._id)
        this.deleteDialog = false
        this.userToDelete = null
      } catch {
        // keep dialog open on failure
      } finally {
        this.deleting = false
      }
    },

    onUserCreated(newUser) {
      this.users.push(newUser)
    },
  },
}
</script>
