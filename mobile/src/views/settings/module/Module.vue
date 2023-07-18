<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "udap-ui/core";
import type {Module} from "udap-ui/types";
import {onMounted, reactive, watchEffect} from "vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import ElementPair from "udap-ui/components/ElementPair.vue";
import Element from "udap-ui/components/Element.vue";

import type {ModuleTiming, RemoteTimings} from "udap-ui/timings";

interface ConfigEntry {
  key: string,
  value: string,
}

const router = core.router();
const remote = core.remote();
const timings = core.timings() as RemoteTimings
const state = reactive({
  moduleId: "" as string,
  module: {} as Module,
  config: [] as ConfigEntry[],
  moduleTimings: {} as ModuleTiming,
  loaded: false
})


onMounted(() => {
  poll()
})

watchEffect(() => {
  poll()
  return remote.modules
})


function poll() {
  if (!router.currentRoute.value) return
  state.moduleId = router.currentRoute.value.params["moduleId"] as string
  let module = remote.modules.find(m => m.id == state.moduleId)
  if (!module) return;
  state.module = module

  state.moduleTimings = timings.getModuleTimings(state.module.uuid)

  if (state.module.config.length > 1) {
    state.config = []
    let ncfg = []
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
    state.config = ncfg
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
</script>

<template>
  <List v-if="state.loaded">
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
          <div :class="`${state.module.state === 'running'?'text-success':''}`"
               class="label-c3 label-w400 label-o4 px-2 text-capitalize d-flex flex-column align-items-end">
            <div>
              {{ state.module.state }}
            </div>
            <Element class=" lh-1 d-flex justify-content-center px-0" foreground
                     style="width: 100%; border-radius: 1.5rem !important;">
              <div class="label-c4 label-w400 d-flex w-100">
                <div class="sf label-o3 label-c4 flex-fill">􀣔</div>
                <div v-if="state.module.state === 'running'" class="px-1">
                  {{
                    convertNanosecondsToString((new Date().valueOf() * 1000 * 1000) - state.moduleTimings.update.stopNano)
                  }}
                </div>
                <div v-else class="px-1">{{
                    convertNanosecondsToString((new Date().valueOf() * 1000 * 1000) - new Date(state.module.updated).valueOf() * 1000 * 1000)
                  }}
                </div>
              </div>
            </Element>

          </div>
        </List>
      </Element>

      <List row style="height: 3.25rem">
        <Element class="d-flex align-items-center justify-content-center gap-1" foreground mutable style="width: 4rem">
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
                 :value="convertNanosecondsToString((new Date().valueOf() - new Date(state.module.updated).valueOf())*1000*1000).toString()"
                 icon="􀅳"></ElementPair>
    <ElementPair v-if="state.module.enabled" :value="convertNanosecondsToString(state.moduleTimings.run.delta)" icon="􀅳"
                 title="Startup"></ElementPair>
    <ElementPair v-if="state.module.enabled" :value="convertNanosecondsToString(state.moduleTimings.update.delta)"
                 icon="􀅳"
                 title="Last Update"></ElementPair>


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

</template>

<style lang="scss" scoped>

</style>