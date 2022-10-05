<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {config} from "@/config"
import type {Controller} from "@/types"
import {PreferenceTypes} from "@/types";
import axios from "axios"
import {Preference} from "@/preferences"
import {onMounted, reactive} from "vue"
import PaneInputBox from "@/components/pane/PaneInputBox.vue";
// The local state of the component
let state = reactive({
  selected: getController(),
  auto: getController(),
  controllers: config.controllers
});

// When the view is mounted, test the default controllers to see which are live
onMounted(() => {
  // Access each controller
  for (let controller of state.controllers) {
    // Test its connectivity and update it's internal state to reflect it
    testController(controller)
  }
})

// Verify a controller is up and running, update the controller if it is
function testController(controller: Controller) {

  // Send a get request to the heartbeat endpoint of the provided controller app
  axios.get(`https://${controller.address}/status`, {
    httpsAgent: {
      rejectUnauthorized: false
    }
  }).then(res => {
    // Set the controller's status to reflect the successful request
    controller.status = true
    // Select the current controller, this will ensure a live controller will be selected by default (if possible)
    suggestController(controller.address)
  }).catch(err => {
    // Set the controller to reflect its down status
    controller.status = false
  })

}

// Recommend a working app controller to the user, typically the production node
function suggestController(address: string) {
  // Update the state to reflect the suggested address
  state.auto = address
}

// Set the preferred controller
function setController(address: string) {
  // Store the preferred controller in localStorage
  new Preference(PreferenceTypes.Controller).set(address);
  // Update the local state to reflect the selection
  state.selected = address
}

// Get the controller stored in localStorage
function getController() {
  return new Preference(PreferenceTypes.Controller).get()
}

function next() {
  window.location.href = "/#/setup/authentication"
}

</script>

<template>

  <div class="d-flex justify-content-center" style="margin-top: 6.25%">
    <PaneInputBox :apply="() => next()" :begin="true" style="width: 26rem !important;"
                  title="Authentication">
      <div class="label-sm label-o5 label-w600 lh-sm px-2">Controller</div>

      <div class="label-o3 label-c1 lh-1 px-2">Please select a control nexus from the list below.</div>
      <div class="d-flex flex-column gap-2 mt-2 px-2 pb-2">
        <div v-for="controller in state.controllers" :key="controller.address"
             :class="`${state.selected === controller.address?'border border-fog':'border border-transparent'}`"
             class="subplot d-flex justify-content-between px-3 py-2" @click="setController(controller.address)">
          <div class="d-flex flex-column justify-content-start">
            <div class="label-md label-o5 label-w500">{{ controller.name }}</div>
            <div class="label-o3 label-w500 label-sm">
                  <span v-if="controller.status" class="text-success"><i
                      class="fa-solid fa-circle-check fa-fw"></i></span>
              <span v-else class="text-warning"><i class="fa-solid fa-triangle-exclamation fa-fw"></i></span>
              {{ controller.address }}
            </div>
          </div>
          <div v-if="state.selected === controller.address" class="d-flex align-items-end flex-column">
            <div class="label-o4 label-w600">Selected</div>
            <div v-if="!controller.status" class="label-o4 label-w300 label-sm text-danger">Unresponsive</div>
          </div>
          <div v-else-if="state.auto === controller.address" class="label-o2 label-w600">Suggested</div>
        </div>

      </div>
    </PaneInputBox>
  </div>

</template>


<style scoped>
.border-fog {
  border-color: rgba(255, 255, 255, 0.25) !important;
}

.border-transparent {
  border-color: transparent !important;
}
</style>