<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, onUnmounted, reactive} from "vue";

interface AtomProps {
  neutrons: number
  protons: number
  electrons: number
}

const props = defineProps<AtomProps>()

const state = reactive({
  ctx: {} as CanvasRenderingContext2D,
  width: 0,
  height: 0,
  runner: 0,
  tick: 0
})

onMounted(() => {
  setupCanvas()
})

onUnmounted(() => {

})


// The Bohr atom is composed of a nucleus surrounded by electron rings.
// The number of electrons per ring is defined by 2*n^2
// We can determine the number of rings from an electron count by sequentially subtracting each ring until we get to 0
function numRings(): number {
  let cursor = props.electrons
  // Handle first ring base case
  if (cursor <= 2) {
    return 1
  }
  cursor -= 2;
  let ring = 2;
  while (true) {
    if (cursor - electronsPerRing(ring) > 0) {
      cursor -= electronsPerRing(ring)
      ring++;
    } else {
      return ring
    }

  }
}

function electronsPerRing(ring: number): number {
  return Math.pow(ring, 2) * 2
}

function drawAtom() {
  state.tick++;
  let ctx = state.ctx
  ctx.strokeStyle = "rgba(128,128,128,0.75)"
  ctx.fillStyle = "rgba(255,255,255,1)"
  ctx.font = "500 24px SF Pro"
  ctx.lineWidth = 4
  const rings = numRings()
  const epr = Object.keys(Array(rings).fill(0)).map(i => electronsPerRing((i * 1) + 1))
  console.log(epr)
  let usedElectrons = 0;

  let ringRadius = Math.max(state.width / 2.5 / numRings(), 100)
  for (let i = 1; i <= rings; i++) {
    ctx.beginPath()
    ctx.ellipse(state.width / 2, state.height / 2, ringRadius * i, ringRadius * i, 0, 0, Math.PI * 2, false)
    ctx.stroke();
    ctx.closePath()

    let numElectrons = Math.min(props.electrons - epr.slice(0, i - 1).reduce((a, b) => a + b, 0), epr[i - 1])
    usedElectrons += numElectrons
    let metrics = ctx.measureText(`${i}`)
    ctx.fillStyle = "rgba(255,255,255,0.7)"
    ctx.fillText(`${i}`, state.width / 2 - metrics.width / 2 + ringRadius / 2 + ringRadius * i, state.height / 2)
    let div = 2 * Math.PI / numElectrons;
    let radius = 12 - Math.log(i * 8)
    let spin = (2 * Math.PI / 1000) * state.tick % 1000 / i * 2
    for (let j = 0; j <= numElectrons; j++) {
      ctx.fillStyle = "rgb(246,204,38)"
      let x = Math.cos(j * div + spin) * ringRadius * i
      let y = Math.sin(j * div + spin) * ringRadius * i
      ctx.beginPath()
      ctx.ellipse(state.width / 2 + x, state.height / 2 + y, radius, radius, 0, 0, Math.PI * 2, false)
      ctx.fill();
      ctx.closePath()

    }
  }

  ctx.beginPath()
  ctx.ellipse(state.width / 2, state.height / 2, ringRadius / 1.5, ringRadius / 1.5, 0, 0, Math.PI * 2, false)
  ctx.fill();
  ctx.closePath()
  ctx.font = "600 32px SF Pro Rounded"
  let metricsP = ctx.measureText(`P ${props.protons}`)
  let metricsN = ctx.measureText(`N ${props.neutrons}`)
  ctx.fillStyle = "rgba(0,0,0,0.7)"
  ctx.fillText(`P ${props.protons}`, state.width / 2 - metricsP.width / 2, state.height / 2 + 32 / 3 - 16)
  ctx.fillText(`N ${props.neutrons}`, state.width / 2 - metricsN.width / 2, state.height / 2 + 32 / 3 + 16)

}

function redraw() {
  state.ctx.clearRect(0, 0, state.width, state.height)
  drawAtom()
  requestAnimationFrame(redraw)
}

function setupCanvas() {
  const chart = document.getElementById(`atom-${props.neutrons}-${props.electrons}-${props.electrons}`) as HTMLCanvasElement
  if (!chart) return;

  const ctx = chart.getContext('2d')
  if (!ctx) return
  state.ctx = ctx
  let scale = 1.75
  ctx.scale(scale, scale)


  chart.width = chart.clientWidth * scale
  chart.height = chart.clientHeight * scale

  state.width = chart.width
  state.height = chart.height
  ctx.translate(0, 0)
  ctx.clearRect(0, 0, state.width, state.height)

  redraw()


}


</script>

<template>
  <canvas :id="`atom-${props.neutrons}-${props.electrons}-${props.electrons}`"
          style="aspect-ratio: 1/1;"></canvas>
</template>

<style scoped>

</style>