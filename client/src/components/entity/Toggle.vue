<script lang="ts" setup>
import {onMounted, reactive} from "vue";
import type {Attribute} from "@/types";

// Define the read-only props passed in from the Attribute Component
const props = defineProps<{
  attribute: Attribute,
  commit: (value: any) => void,
  small: boolean,
}>()

// Define the local state for the toggle
const state = reactive<{
  latest: Date,
  local: Attribute,
  waiting: boolean,
  active: boolean,
}>({
  latest: new Date(),
  local: props.attribute,
  waiting: false,
  active: false,
})

// Define local styles for the toggle
const styles = {
  on: {
    min: 0,
    max: 100,
    step: 5,
    unit: '%',
    icon: '􀇯',
    class: 'slider-dim',
  },
}

// Update the local state when view mounted
onMounted(() => {
  updateState()
})

// Send the changes immediately when selected
function commitChanges(_: MouseEvent) {
  state.local.request = state.local.request === "false" ? "true" : "false"
  props.commit(state.local)
  state.waiting = true
  state.latest = new Date()
  updateState()
}

function updateState() {
  state.active = state.local.request === "true"
}

</script>

<template>
  <div v-if="small">
    <div class="h-bar gap label-sm label-w600 text-uppercase label-o4 px-2">
      <div @click="commitChanges">{{ state.active ? "ON" : "OFF" }}</div>
    </div>
  </div>
  <div v-else v-if="attribute" class="element" v-on:click.stop>

    <div class="h-bar justify-content-start align-items-center align-content-center">
      <div class="label-xxs label-o2 label-w600">{{ styles.on.icon }}</div>
      <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ props.attribute.key }}</div>
      <div class="fill"></div>

      <div class="h-bar gap label-sm label-w600 text-uppercase label-o4 px-2" @click="commitChanges">
        <div>{{ state.active ? "ON" : "OFF" }}</div>
      </div>
    </div>


  </div>
</template>

<style scoped>

</style>
