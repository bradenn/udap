<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Module, Remote, Timing} from "@/types";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import Confirm from "@/components/plot/Confirm.vue";
import axios from "axios";

let remote = inject('remote') as Remote
let preferences = inject('preferences')

let state = reactive({
  modules: {} as Module[],
  timings: {} as Timing[],
  histories: new Map<string, number[]>(),
  maxes: new Map<string, number>(),
  threads: [] as number[],
  loading: true
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.modules = remote.modules.sort((a, b) => a.created > b.created ? 1 : -1) as Module[]
  state.threads = remote.metadata.system.threads as []
  remote.timings.filter(t => state.modules.find(m => m.id === t.pointer)).forEach(t => {
    let local = state.histories.get(t.pointer) || []
    let cand = t.delta / 1000
    if (local.includes(cand)) return
    local.push(cand)
    if (local.length > 20) {
      local = local.slice(1)
    }
    state.histories.set(t.pointer, local)
  })
  state.loading = false
  return remote
}

function toggleEnabled(id: string, enabled: boolean) {
  if (enabled) {
    axios.post(`http://localhost:3020/modules/${id}/enable`)
  } else {
    axios.post(`http://localhost:3020/modules/${id}/disable`)
  }
}

// groupBy creates several arrays of elements based on the value of a key
function groupBy<T>(xs: T[], key: string): T[] {
  return xs.reduce(function (rv: any, x: any): T {
    (rv[x[key]] = rv[x[key]] || []).push(x);
    return rv;
  }, {});
}
</script>

<template>
  <div>
    <div class="d-flex justify-content-between py-2 px-1">
      <div class="d-flex justify-content-start">
        <div class="label-w500 label-o3 label-xxl"><i :class="`fa-solid fa-layer-group fa-fw`"></i></div>
        <div class="label-w500 opacity-100 label-xxl px-2">Modules</div>
      </div>
      <div v-if="!state.loading" class="w-50">
        <div class="activity-container">
          <div v-for="(thread, k) in remote.metadata.system.threads" class="element tick-track"
               style="border-radius: 10px; padding-inline: 7px;">

            <div class="label-c3 label-o2 label-r label-w600 justify-content-center d-flex"
                 style="width: 0.8rem; height: 16px;">
              {{
                k + 1
              }}&nbsp;
            </div>
            <div :style="`width:${thread}%;`" class="tick">
            </div>
            <div :style="`width:${100-thread}%;`" class="tick-grey">
            </div>
          </div>

        </div>
      </div>
    </div>
    <div v-if="!state.loading" class="d-flex flex-column gap-1">


      <div class="module-container">

        <div v-for="module in state.modules" class="">

          <Plot :key="module.id" :cols="2" :rows="1">
            <div class="subplot subplot-inline px-0">
              <div :style="`background-color: rgba(${module.enabled?'25, 135, 84':'135, 100, 2'}, 0.53);`"
                   class="status-marker"></div>
              <div class="w-100">
                <div>
                  <div class="label-c1 label-o5 label-r lh-1">
                    <div class="d-flex justify-content-between">
                      <div>{{ module.name }}</div>

                    </div>

                  </div>
                  <div class="label-c4  label-o3 label-r py-0 overflow-ellipse" style="line-height: 0.55rem">{{
                      module.description
                    }}
                  </div>
                </div>

              </div>
              <!--              {{-->
              <!--                new Date().valueOf() - -->
              <!--                new Date(remote.timings.filter(t => t.pointer === module.id).valueOf()-->
              <!--              }}-->


              <div class="label-c3 label-o4 d-flex flex-column justify-content-end align-items-end">
                <div :class="`${module.enabled?'text-success':''}`" class="label-o3 text-uppercase">
                  &nbsp;{{ module.enabled ? module.state : 'Disabled' }}
                </div>
                <div v-if="state.histories" class="label-c3 label-o3 d-flex flex-row align-items-end time-marker-line">
                  <div v-for="marker in state.histories.get(module.id)?.map(d => d / 1000)"
                       :style="`height:${Math.log(marker)}px;`"
                       class="time-marker"></div>
                </div>
              </div>
            </div>
            <div class="d-flex gap-1 text-success justify-content-center">


              <Confirm v-if="!module.enabled" :fn="() => toggleEnabled(module.id, !module.enabled)" icon="􀊃"
                       title="Enable"></Confirm>
              <Confirm v-if="module.enabled" :fn="() => toggleEnabled(module.id, !module.enabled)" icon="􀆧"
                       title="Disable"></Confirm>
              <Radio :active="false" :fn="() => {}" sf="􀍟" style="width: 2.5rem;" title=""></Radio>
            </div>
          </Plot>
        </div>
      </div>


    </div>
    <div v-else class="element p-2">
      <div class="label-c1 label-o4 d-flex align-content-center gap-1">
        <div>
          <Loader size="sm"></Loader>
        </div>
        <div class="">Loading...</div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.tick {
  height: 18px;
  width: 100%;
  background-color: rgba(25, 135, 84, 0.3);
  border-radius: 2px;
  transition: width 250ms ease-in-out;
}

.tick-track {
  display: flex;
  align-items: start;
  align-content: start;
  line-height: 18px;
  justify-content: start;
  padding: 6px;
  gap: 0px;
}

.tick-grey {
  height: 18px;
  width: 100%;
  background-color: rgba(64, 64, 64, 0.2);
  border-radius: 2px;
  transition: width 250ms ease-in-out;
}

.overflow-ellipse {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis !important;
  text-wrap: none !important;
  max-width: 5.75rem;
}

.activity-container {
  width: 100%;

  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(2, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

.time-marker-line {
  display: flex;
  flex-direction: row;
  justify-content: center;
  height: 20px;
  width: 75px;
  align-items: center;
  gap: 1px;
  border-radius: 6px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}

.time-marker {
  width: 2px;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  height: 1px;
}

.status-marker {
  width: 4px !important;
  border-radius: 4px;
  height: 28px;

  margin-right: 14px;
  margin-left: 8px;


  background-color: rgba(25, 135, 84, 0.53);
}

.module-container {

  display: grid;
  grid-column-gap: 0.5rem;
  grid-row-gap: 0.5rem;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

</style>