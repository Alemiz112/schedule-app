import Vue from "vue"
import Vuex from "vuex"
import { get as apiGet } from "@/utils"
import {
  createFolder as createFolderService,
  deleteFolder as deleteFolderService,
  setEventFolder as setEventFolderService,
  updateFolder as updateFolderService,
} from "../utils/services/FolderService"
import { archiveEvent as archiveEventService } from "../utils/services/EventService"

Vue.use(Vuex)

const showErrorAction = ({ commit }, error) => {
  commit("setError", "")
  setTimeout(() => commit("setError", error), 0)
}

const showInfoAction = ({ commit }, info) => {
  commit("setInfo", "")
  setTimeout(() => commit("setInfo", info), 0)
}

const refreshAuthUserAction = async ({ commit }) => {
  const authUser = await apiGet("/user/profile")
  commit("setAuthUser", authUser)
}

const openNewDialogAction = ({ commit }, { eventOnly = false, folderId = null }) => {
  commit("setNewDialogOptions", {
    show: true,
    contactsPayload: {},
    openNewGroup: false,
    eventOnly,
    folderId,
  })
}

const getEventsAction = async ({ commit, dispatch, state }) => {
  if (!state.authUser) return null
  try {
    commit("setFolders", await apiGet("/user/folders"))
    commit("setEvents", await apiGet("/user/events"))
  } catch (err) {
    dispatch("showError", "There was a problem fetching events!")
    console.error(err)
  }
  return null
}

const archiveEventAction = async ({ dispatch, state }, { eventId, archive }) => {
  try {
    await archiveEventService(eventId, archive)
    const event = state.events.find((e) => e._id === eventId)
    if (event) {
      event.isArchived = archive
    }
  } catch (err) {
    dispatch("showError", "There was a problem archiving the event!")
    console.error(err)
  }
}

const createFolderAction = async ({ commit, dispatch }, { name, color }) => {
  try {
    const folder = await createFolderService(name, color)
    commit("addFolder", {
      _id: folder.id,
      name,
      color,
      eventIds: [],
    })
  } catch (err) {
    dispatch("showError", "There was a problem creating the folder!")
    console.error(err)
  }
}

const updateFolderAction = async ({ commit, dispatch }, { folderId, name, color }) => {
  try {
    await updateFolderService(folderId, name, color)
    commit("updateFolder", { folderId, name, color })
  } catch (err) {
    dispatch("showError", "There was a problem updating the folder!")
    console.error(err)
  }
}

const deleteFolderAction = async ({ commit, dispatch }, folderId) => {
  try {
    await deleteFolderService(folderId)
    commit("removeFolder", folderId)
  } catch (err) {
    dispatch("showError", "There was a problem deleting the folder!")
    console.error(err)
  }
}

const setEventFolderAction = async ({ commit, dispatch }, { eventId, folderId }) => {
  try {
    commit("removeEventFromFolder", eventId)
    commit("addEventToFolder", { eventId, folderId })
    await setEventFolderService(eventId, folderId)
  } catch (err) {
    dispatch("showError", "There was a problem moving the event!")
    console.error(err)
  }
}

export default new Vuex.Store({
  state: {
    error: "",
    info: "",

    authUser: null,

    events: [],
    folders: [],

    outlookEnabled: !!process.env.VUE_APP_MICROSOFT_CLIENT_ID,

    featureFlagsLoaded: false,

    // Feature flags
    groupsEnabled: true,
    signUpFormEnabled: false,
    daysOnlyEnabled: true,
    overlayAvailabilitiesEnabled: true,

    // Experiments
    pricingPageConversion: "control",

    // New dialog
    newDialogOptions: {
      show: false,
      contactsPayload: {},
      openNewGroup: false,
      eventOnly: false,
      folderId: null,
    },
  },
  mutations: {
    setError(state, error) {
      state.error = error
    },
    setInfo(state, info) {
      state.info = info
    },

    setAuthUser(state, authUser) {
      state.authUser = authUser
    },

    setEvents(state, events) {
      state.events = events
    },
    setFolders(state, folders) {
      state.folders = folders
    },

    setFeatureFlagsLoaded(state, loaded) {
      state.featureFlagsLoaded = loaded
    },
    setGroupsEnabled(state, enabled) {
      state.groupsEnabled = enabled
    },
    setSignUpFormEnabled(state, enabled) {
      state.signUpFormEnabled = enabled
    },
    setDaysOnlyEnabled(state, enabled) {
      state.daysOnlyEnabled = enabled
    },
    setOverlayAvailabilitiesEnabled(state, enabled) {
      state.overlayAvailabilitiesEnabled = enabled
    },
    setPricingPageConversion(state, conversion) {
      state.pricingPageConversion = conversion
    },
    addFolder(state, folder) {
      state.folders.push(folder)
    },
    updateFolder(state, { folderId, name, color }) {
      const folder = state.folders.find((f) => f._id === folderId)
      if (folder) {
        folder.name = name
        folder.color = color
      }
    },
    removeFolder(state, folderId) {
      state.folders = state.folders.filter((f) => f._id !== folderId)
    },
    removeEventFromFolder(state, eventId) {
      state.folders.forEach((folder) => {
        folder.eventIds = folder.eventIds.filter((id) => id !== eventId)
      })
    },
    addEventToFolder(state, { eventId, folderId }) {
      const folder = state.folders.find((f) => f._id === folderId)
      if (folder) {
        folder.eventIds.push(eventId)
      }
    },

    setNewDialogOptions(
      state,
      {
        show = false,
        contactsPayload = {},
        openNewGroup = false,
        eventOnly = true,
        folderId = null,
      }
    ) {
      state.newDialogOptions = {
        show,
        contactsPayload,
        openNewGroup,
        eventOnly,
        folderId,
      }
    },
  },
  actions: {
    showError: showErrorAction,
    showInfo: showInfoAction,
    refreshAuthUser: refreshAuthUserAction,
    openNewDialog: openNewDialogAction,
    getEvents: getEventsAction,
    archiveEvent: archiveEventAction,
    createFolder: createFolderAction,
    updateFolder: updateFolderAction,
    deleteFolder: deleteFolderAction,
    setEventFolder: setEventFolderAction,
  },
  modules: {},
})
