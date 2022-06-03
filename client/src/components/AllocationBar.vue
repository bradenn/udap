<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {onMounted, reactive, watch} from "vue";


export interface Allocation {
  name: string
  allocated: number
  source?: string
  active: number
  battery?: boolean
  priority?: number
}

export interface AllocationBar {
  name: string
  allocatable: number
  allocations: Allocation[]
  isFixed?: boolean
}

let props = defineProps<AllocationBar>()

let state = reactive({
  allocatedLoad: 0,
  activeLoad: 0,
})

watch(props, (newProps, oldProps) => {
  calculateLoad()
})

onMounted(() => {
  calculateLoad()
})

function calculateLoad() {
  state.allocatedLoad = props.allocations.reduce((a: number, b: Allocation) => a + b.allocated, 0)
  state.activeLoad = props.allocations.reduce((a: number, b: Allocation) => a + b.active, 0)
}

function roundSmooth(input: number): number {
  return Math.round(input * 10) / 10
}

</script>

<template>
  <div v-if="props.allocations">

    <div class="d-flex justify-content-between align-items-end">
      <div class="label-xs label-r label-w400 label-o5 px-2">{{ props.name }}</div>
      <div class="label-c1 label-r label-w400 label-o4 px-2">{{
          roundSmooth((props.isFixed ? state.allocatedLoad : state.activeLoad) * 120)
        }}<span class="label-o2 label-c3 label-w600 label-r">W</span> / {{ roundSmooth(props.allocatable * 120) }}<span
            class="label-o2 label-c3 label-w600 label-r">W</span>
      </div>

    </div>
    <div v-if="props.isFixed">
      <div class="power-chart element w-100 ">
        <div v-for="device in props.allocations.sort((a, b) => b.active-a.active)"
             :style="`width: ${(device.allocated/props.allocatable)*100}% !important;`"
             class="power-chart-point element-secondary justify-content-center">

          <div :style="`width: ${(device.active/device.allocated)*100}% !important;`"
               class="power-chart-meta element-secondary "></div>

          <div class="label-w500 label-o4 label-r label-c2 lh-1"
               style="overflow: hidden; text-overflow: ellipsis; max-width: 100%;">
            {{ device.name }}
          </div>
          <div class="label-o2  label-c3">{{ roundSmooth(device.active * 120) }}W
          </div>
          <div class="label-o2  label-c3"></div>
          <span v-if="device.battery" class="label-c3 text-success">􀛨</span>
        </div>

        <div :style="`width: ${((props.allocatable-state.allocatedLoad)/ props.allocatable)*100 }% !important;`"
             class="power-chart-point subplot d-flex justify-content-center align-items-start px-2">
          <div class="label-w500 label-o3 label-r label-c2 lh-1">Unallocated</div>
          <div class="label-o2  label-c3">{{ roundSmooth((props.allocatable - state.allocatedLoad) * 120) }}W</div>

        </div>
      </div>
    </div>
    <div v-else>
      <div class="power-chart element w-100">
        <div v-for="(device, i) in props.allocations.sort((a, b) => b.active-a.active)"
             :style="`width: ${(device.active/props.allocatable)*100}% !important;`"
             class="power-chart-point element-secondary justify-content-center">

          <div class="label-w500 label-o4 label-r label-c2 lh-1 "
               style="overflow: hidden; text-overflow: ellipsis; max-width: 100%;">{{ device.name }}
          </div>
          <div class="label-o2  label-c3">{{ roundSmooth(device.active * 120) }}W</div>

          <span v-if="device.battery" class="label-c3 text-success">􀛨</span>
        </div>

        <div :style="`width: ${((props.allocatable - state.activeLoad) / props.allocatable)*100}% !important;`"
             class="power-chart-point subplot d-flex justify-content-center align-items-start px-2">
          <div class="label-w500 label-o3 label-r label-c2 lh-1">Unallocated</div>

          <div class="label-o2  label-c3 ">{{ roundSmooth((props.allocatable - state.activeLoad) * 120) }}W
          </div>

        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.power-chart-meta {
  position: absolute;
  top: 0;
  left: 0;
  border-radius: 0.25rem;
  height: 100%;
}

.power-chart-buttons {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
}

.power-chart-id {
  position: absolute;
  top: 0;
  left: 0;

  padding: 0.125rem 0.25rem;
}

.power-chart {
  display: flex;
  justify-content: start;
  font-size: 0.6rem;
  font-weight: 500;
  align-content: start;
  align-items: start;
  gap: 0.18rem;
}

.power-chart-point {
  min-width: 2.75rem;
  border-radius: 0.25rem;
  height: 3rem;
  width: 8px;
  display: flex;
  flex-direction: column;
  align-items: start;
  justify-content: start;
  position: relative;
  padding: 0.25rem;
}
</style>