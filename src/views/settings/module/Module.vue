<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "udap-ui/core";
import type {Module} from "udap-ui/types";
import {Timing} from "udap-ui/types";
import {onMounted, reactive, watchEffect} from "vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import ElementPair from "udap-ui/components/ElementPair.vue";
import Element from "udap-ui/components/Element.vue";
import Time from "udap-ui/components/Time.vue";

import type {RemoteTimings} from "udap-ui/timings";

interface ModuleTiming {
  loaded: boolean,
  run: Timing
  update: Timing
}

interface ConfigEntry {
  key: string,
  value: string,
}

let stub: Timing = {
  pointer: "",
  name: "",
  start: new Date().toString(),
  startNano: new Date().valueOf() * 1000,
  stop: new Date().toString(),
  stopNano: new Date().valueOf() * 1000,
  delta: 0,
  frequency: 0,
  complete: true,
  depth: 0,
  id: "id"
}
const router = core.router();
const remote = core.remote();
const timings = core.timings() as RemoteTimings
const state = reactive({
  moduleId: "" as string,
  module: {} as Module,
  config: [] as ConfigEntry[],
  moduleTimings: {
    run: null,
    update: null,
    loaded: false
  } as ModuleTiming,
  loaded: false
})


onMounted(() => {
  poll()

})

watchEffect(() => {
  poll()
  return remote.modules
})


watchEffect(() => {
  updateModuleTimings()
  return timings.timings
})

function updateModuleTimings() {
  let moduleRun = `module.${state.module.uuid}.run`
  let moduleUpdate = `module.${state.module.uuid}.update`
  state.moduleTimings.run = timings.timings.find(t => t.pointer == moduleRun)
  state.moduleTimings.update = timings.timings.find(t => t.pointer == moduleUpdate)
  state.moduleTimings.loaded = true
}


function poll() {
  if (!router.currentRoute.value) return
  state.moduleId = router.currentRoute.value.params["moduleId"] as string
  let module = remote.modules.find(m => m.id == state.moduleId)
  if (!module) return;
  state.module = module
  // state.moduleTimings = timings.getModuleTimings(state.module.uuid) as ModuleTiming
  updateModuleTimings()
  if (state.module.config.length > 1) {
    state.config = []
    let ncfg: ConfigEntry[] = []
    let cfg = JSON.parse(state.module.config)

    if (cfg instanceof Object) {
      let keys = Object.keys(cfg)
      keys.forEach(k => {
        ncfg.push({
          key: k,
          value: cfg[k]
        } as ConfigEntry)

      })

    }
    state.config = ncfg as ConfigEntry[]
  }
  state.loaded = true
}


function convertNanosecondsToString(nanoseconds: number): string {
  if (nanoseconds < 0) {
    nanoseconds = Math.abs(nanoseconds)
    // throw new Error("Input must be a non-negative number.");
  }

  if (nanoseconds === 0) {
    return "0ns";
  }

  const units = ["ns", "µs", "ms", "s", "ks"];
  let value = nanoseconds;
  let unitIndex = 0;

  while (value >= 1000 && unitIndex < units.length - 1) {
    if (unitIndex >= 3) {
      let seconds = value % 60
      let minutes = (value / 60) % 60
      let hours = (value / 60 / 60) % 60
      let days = Math.round((value / 60 / 60 / 24) % 24)
      let years = Math.round((value / 60 / 60 / 24 / 365))
      return `${years > 0 ? (years + 'y') : ''} ${(days > 0) ? (Math.round(days) + 'd') : ''} ${Math.round(hours)}h  ${Math.round(minutes)}m ${Math.round(seconds)}s`;
    }
    value /= 1000;
    unitIndex++;
  }

  return `${value.toFixed(2)}${units[unitIndex]}`;
}

function init() {
  Notification.requestPermission().then((result) => {
    if (result === "granted") {
      alert("Yippie!");
    }
  })
}
</script>

<template>

  <List v-if="state.loaded && state.moduleTimings.loaded">

    <List>
      <Element class="" foreground>
        <List class="justify-content-between align-items-center" row>
          <div class="d-flex align-items-start">
            <div class="notches px-1 mt-2 gap-1" style="height: 1.25rem">
              <div :class="`${true?'active':''}`" class="notch h-100"></div>
              <div :class="`${true?'active':''}`" class="notch h-100"></div>

            </div>
            <div class="d-flex flex-column gap-0 px-2">
              <div class="label-c1 label-w500 label-o5 lh-1">{{ state.module.name }}</div>
              <div class="label-c5 label-w400 label-o4 lh-1">v{{ state.module.version }}</div>
            </div>
          </div>
          <div></div>
          <div :class="`${state.module.state === 'running'?'text-success':''}`"
               class="label-c3 label-w400 label-o4 px-2 text-capitalize d-flex flex-column align-items-end">
            <div>
              {{ state.module.state }}
            </div>

            <Element class=" lh-1 d-flex justify-content-center px-0" foreground
                     v-if="false" style="width: 100%; border-radius: 1.5rem !important;">
              <div class="label-c4 label-w400 d-flex w-100">
                <div class="sf label-o3 label-c4 flex-fill">􀣔</div>
                <div v-if="state.module.state === 'running'" class="px-1">
                  <div v-if="!state.moduleTimings.run">
                    Not running.
                  </div>
                  <div v-else>
                    <Time :since="(new Date().valueOf() * 1000 *1000) - state.moduleTimings.run.stopNano" nano></Time>
                  </div>
                </div>
                <div v-else class="px-1">
                  <Time :since="(new Date().valueOf() * 1000 * 1000) - new Date(state.module.updated).valueOf() * 1000 * 1000"
                        nano></Time>

                </div>
              </div>
            </Element>

          </div>
        </List>
      </Element>

      <List row style="height: 3.25rem">
        <Element :cb="init" class="d-flex align-items-center justify-content-center gap-1" foreground mutable
                 style="width: 4rem">
          <div class="sf">􀊯</div>
          Reload
        </Element>

        <Element class="d-flex align-items-center justify-content-center gap-1" foreground mutable style="width: 4rem">
          <div class="sf">􀛶</div>
          Halt
        </Element>
      </List>
    </List>
    <ElementHeader title="Runtime"></ElementHeader>
    <ElementPair :title="!state.module.enabled?'Downtime':'Uptime'"
                 alt="Time since startup" icon="􀯛">
      <div v-if="state.moduleTimings.run">
        <Time :since="new Date(state.moduleTimings.run.stop).valueOf()" live></Time>
      </div>
      <div v-else>
        No Data
      </div>

    </ElementPair>
    <ElementPair v-if="state.module.enabled" alt="Module startup duration"
                 icon="􀣔" title="Startup">
      <div v-if="state.moduleTimings.run">

        <Time :since="state.moduleTimings.run.delta"
              nano></Time>
      </div>
      <div v-else>
        No Data
      </div>

    </ElementPair>
    <ElementHeader title="Update"></ElementHeader>
    <ElementPair v-if="state.module.enabled"
                 alt="Time since the last update"
                 icon="􀣔" title="Last Update">
      <div v-if="state.moduleTimings.update">
        <Time :since="new Date(state.moduleTimings.update.stop).valueOf()" live
              precise></Time>
      </div>
      <div v-else>
        No Data
      </div>

    </ElementPair>
    <ElementPair v-if="state.module.enabled"
                 alt="Module update duration"
                 icon="􀣔" title="Update">
      <div v-if="state.moduleTimings.update">
        <Time :since="new Date(state.moduleTimings.update.delta).valueOf()"
              nano></Time>
      </div>
      <div v-else>
        No Data
      </div>

    </ElementPair>
    <div class=""></div>
    <ElementPair v-if="state.module.enabled"
                 alt="Prescribed update interval"
                 icon="􀐱" title="Interval">
      <Time :since="state.module.interval"
            nano></Time>

    </ElementPair>
    <ElementPair v-if="state.module.enabled"
                 alt="Previous update interval"
                 icon="􀐱" title="Actual Interval">
      <div v-if="state.moduleTimings.update">
        <Time :since="new Date(state.moduleTimings.update.frequency).valueOf()"
              nano></Time>
      </div>
      <div v-else>
        No Data
      </div>

    </ElementPair>

    <ElementHeader title="Module"></ElementHeader>
    <ElementPair :value="state.module.name" icon="􀅳" title="Name"></ElementPair>
    <ElementPair :value="state.module.enabled?'Enabled':'Disabled'" icon="􀅳" title="Runtime"></ElementPair>

    <ElementHeader title="Instance"></ElementHeader>
    <ElementPair :value="state.module.state" icon="􀅳" title="State"></ElementPair>
    <ElementPair :value="state.module.uuid" copyable icon="􀅳" title="UUID"></ElementPair>


    <ElementHeader title="Implementation"></ElementHeader>
    <ElementPair v-if="state.module.enabled"
                 :value="convertNanosecondsToString((new Date().valueOf() - new Date(state.module.created).valueOf())*1000*1000).toString()"
                 icon="􀅳"
                 title="Created"></ElementPair>
    <ElementPair :value="state.module.type" icon="􀅳" title="Classification"></ElementPair>
    <ElementPair :value="state.module.path.replace('modules/', '')" icon="􀅳" title="Path"></ElementPair>
    <ElementPair :value="state.module.id" copyable icon="􀅳" title="ID"></ElementPair>
    <ElementHeader title="Metadata"></ElementHeader>
    <ElementPair :value="state.module.author" icon="􀅳" title="Author"></ElementPair>
    <ElementPair :value="state.module.version" icon="􀅳" title="Version"></ElementPair>
    <ElementHeader title="Configuration"></ElementHeader>
    <ElementPair v-for="key in state.config" v-if="state.config.length>0" :key="key.key" :title="key.key"
                 :value="key.value"
                 icon="􀅳"></ElementPair>

    <Element v-else class="px-4" foreground>No Configuration</Element>
  </List>
  <List v-else>
    Loading
  </List>

</template>

<style lang="scss" scoped>

</style>