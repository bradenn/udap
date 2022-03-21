<script lang="ts" setup>


import {onMounted, reactive} from "vue";

let state = reactive<{
  canvas: HTMLCanvasElement
}>({
  canvas: {} as HTMLCanvasElement
})

onMounted(() => {
  drawCanvas()
})

function drawCanvas() {
  const canvas = document.getElementById("canvas-fractal") as HTMLCanvasElement
  const ctx = canvas.getContext('2d')
  if (!ctx) return

}

</script>

<template>
  <div class="canvas-container">
    <canvas id="canvas-fractal"></canvas>
  </div>
</template>

<style lang="scss" scoped>
.canvas-container {
  height: 200px;
  width: 200px;

  aspect-ratio: 1/1 !important;
}

</style>
