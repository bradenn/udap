<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import {onMounted} from "vue";


onMounted(() => {
  setup()
})

function SphericalToPlanar(phi: number, theta: number, height: number): { x: number, y: number } {

  let r = Math.tan(theta) * height;
  let x = Math.cos(phi) * r;
  let y = Math.sin(phi) * r;


  return {x, y};
}

function PlanarToSpherical(cx: number, cy: number, x: number, y: number, height: number): {
  phi: number,
  theta: number
} {

  let r: number = Math.sqrt(Math.pow(cx - x, 2) + Math.pow(cy - y, 2));
  let phi: number = Math.atan(y / x) + (x >= 0 ? 0 : Math.PI);
  let theta: number = (Math.atan(r / height));


  return {phi, theta};
}

function setup() {
  let element = document.getElementById("beamDemoCanvas") as HTMLElement
  if (!element) return
  let canvas = element as HTMLCanvasElement
  if (!canvas) return

  let ctx = canvas.getContext("2d") as CanvasRenderingContext2D
  if (!ctx) return

  ctx.fillStyle = "rgba(255,255,255,1)"
  ctx.strokeStyle = "rgba(255,255,255,1)"
  ctx.rect(10, 10, 100, 100)
  ctx.stroke()
  ctx.beginPath()
  let lx = 0;
  let ly = 0
  for (let i = 0; i < 100; i++) {
    // let {x1, y1} = SphericalToPlanar(0, (Math.PI / 180) * i, 1)
    // let {x2, y2} = SphericalToPlanar(0, (Math.PI / 180) * i, 1)
    //
    // ctx.moveTo(x1, y1)
    // ctx.lineTo(x2, y2)
  }
  ctx.closePath()
  ctx.stroke()
}

</script>

<template>
  <Element>
    <canvas id="beamDemoCanvas" style=""></canvas>
  </Element>
</template>

<style>

</style>