<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {v4 as uuidv4} from "uuid";

import {onMounted, reactive} from "vue";
import Wave from "@/views/terminal/diagnostic/pages/Wave.vue";
import Wrap from "@/views/terminal/diagnostic/pages/Wrap.vue";

const state = reactive({
  uuid: uuidv4(),
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  ready: false,
  cycles: 1
})

onMounted(() => {
  animate()
})

function animate() {

  state.cycles = (state.cycles + 0.01) % 24

  requestAnimationFrame(animate)
}

</script>

<template>
  <div class="transform-grid">
    <div class="d-flex flex-column gap-1">
      <div class="element px-2 py-2">
        <div class="label-c2 label-o3 label-w600">
          Cycle: {{ Math.round(state.cycles * 100) / 100 }}
        </div>
      </div>
    </div>

    <div class="d-flex flex-column gap-1">
      <Wave :frequency="240" :phase="0"></Wave>
      <Wave :frequency="120" :phase="0"></Wave>
      <Wave :frequencies="[240, 120, 360, 320, 60]" :phase="0"></Wave>
    </div>
    <div class="d-flex flex-column gap-1">
      <Wrap :cycles="state.cycles" :frequencies="[240, 120, 360, 320, 60]"
            :phase="0"></Wrap>
    </div>
    <div class="d-flex flex-column gap-1">
      <Wrap :cycles="state.cycles" :frequency="120" :phase="0"></Wrap>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.transform-grid {
  display: grid;
  grid-template-rows: repeat(3, 1fr);
  grid-template-columns: repeat(4, 1fr);
  grid-gap: 0.25rem;
}

.inner-canvas {
  width: 100%;
  height: 100%;

}

.canvas-container {
  display: flex;
  flex-direction: row;
  justify-content: center;
  height: 100%;
  width: 100%;
  align-items: center;
  border-radius: 6px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}
</style>