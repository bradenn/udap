<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import PaneMenuToggle from "@/components/pane/PaneMenuToggle.vue";
import type {Macro, Remote, Zone} from "@/types";
import Tasks from "@/components/task/Tasks.vue";
import Task from "@/components/task/Task.vue";
import Keyboard from "@/components/Keyboard.vue";
import macroService from "@/services/macroService";

const router = useRouter()

const remote = inject("remote") as Remote

interface InputField {
  name: string,
  description: string,
  value: any,
  type: string,
  values?: any[]
}

const state = reactive({
  zone: {
    value: "",
    preview: "",
    zones: [] as Zone[]
  },
  name: {
    value: "",
    preview: ""
  },
  key: {
    value: "",
    preview: "",
    keys: [] as string[]
  },
  value: {
    value: "",
    preview: ""
  },
  field: 0,
})

onMounted(() => {
  updateRemote()
})

watchEffect(() => {
  updateRemote()
  return remote
})

function sortZones(a: Zone, b: Zone): number {
  if (a.pinned && !b.pinned) {
    return -1;
  } else if (!a.pinned && b.pinned) {
    return 1;
  } else if (a.pinned && b.pinned && a.name > b.name) {
    return 2;
  } else if (a.pinned && b.pinned && a.name < b.name) {
    return -2;
  }
  return 0
}

function updateRemote() {
  let zones = remote.zones.filter(z => !z.deleted).sort(sortZones);
  if (!zones) return;
  state.zone.zones = zones
}

function selectZone(id: string) {

}

function update(value: string) {
  if (value === "{bksp}") {
    state.name.value = state.name.value.substring(0, state.name.value.length - 1)
  } else if (value === "{enter}") {
    state.field++;
  } else if (value === "{space}") {
    state.name.value += " "
  } else {
    state.name.value += value
  }
}

function createSubroutine() {
  macroService.createMacro({
    zone: state.zone.value,
    name: state.name.value,
    type: state.key.value,
    value: state.value.value,
  } as Macro).then(res => {
    console.log(res)
  }).catch(err => {
    console.log(err)
  })
}

function next(value: number) {
  state.field = value
}

</script>

<template>
  <div class="generic-grid">
    <Tasks title="New Subroutine">
      <Task :active="state.field === 0"
            :index="0"
            :next="next"
            :preview="state.zone.preview"
            :value="state.zone.value"
            description="Which zone should this macro act on?"
            title="Zone">
        <PaneMenuToggle v-for="item in state.zone.zones" :active="state.zone.value === item.id"
                        :fn="() => selectZone(item.id)"
                        :subtext="item.entities.map(e => e.name).join(', ')"
                        :title="item.name"></PaneMenuToggle>
      </Task>
      <div class="element">
        <div class="subplot justify-content-center label-w500 label-c1" style="height: 1.8rem"
             @click="createSubroutine">Create Subroutine</div>
      </div>
    </Tasks>


    <div v-if="state.field === 2" class="kb-input">
      <Keyboard :input="update" class="keyboard-location" keySet="" keyboardClass="simple-keyboard"></Keyboard>
    </div>
  </div>
</template>

<style scoped>


.kb-input {
  position: absolute;
  bottom: 4rem;
}

.generic-grid > div {
  width: 18rem;
}

.generic-grid {
  display: flex;
  justify-content: center;
  /*width: 100%;*/
  /*height: 100%;*/
  /*display: grid;*/
  grid-column-gap: 0.25rem;
  /*grid-row-gap: 0.25rem;*/
  /*grid-template-rows: repeat(1, 1fr);*/
  /*grid-template-columns: repeat(3, 1fr);*/
}
</style>