<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {onMounted, reactive} from "vue";


export interface Allocation {
  name: string
  allocated: number
  source?: string
  active: number
  priority?: number
}

export interface AllocationBar {
  name: string
  allocatable: number
  allocations: Allocation[]
}

let props = defineProps<AllocationBar>()

let state = reactive({
  load: 0,
})

onMounted(() => {
  calculateLoad()
})

function calculateLoad() {
  state.load = Math.round(props.allocations.reduce((a: number, b: Allocation) => a + b.active, 0))
}

</script>

<template>
  <div>
    <div v-if="props.allocations">
      <h5 class="mb-1 px-2">{{ name }} <span class="label-xs label-o4 float-end label-w500">{{
          state.load
        }}A / {{ props.allocatable }}A</span></h5>
      <div class="power-chart element w-100 ">
        <div v-for="(device, i) in props.allocations.sort((a, b) => b.active-a.active)"
             :style="`width: ${(device.allocated/props.allocatable)*100}% !important;`"
             class="power-chart-point element">

          <div :style="`width: ${(device.active/device.allocated)*100}% !important;`"
               class="power-chart-meta element"></div>
          <div>{{ device.name }}</div>
          <div class="label-c3 label-o4 label-mono">{{ device.allocated }}A</div>
        </div>
        <div class="power-chart-buttons">

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
  height: 4rem;
  width: 8px;
  display: flex;
  flex-direction: column;
  align-items: start;
  justify-content: start;

}
</style>