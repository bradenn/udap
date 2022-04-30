<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import {onMounted, reactive} from "vue";
import Plot from "@/components/plot/Plot.vue";
import Toggle from "@/components/plot/Toggle.vue";
import Radio from "@/components/plot/Radio.vue";


const colors = [
  {
    name: "White",
    style: "rgba(255,255,255, 1)"
  },
  {
    name: "Blue",
    style: "rgb(0,110,255)"
  },
  {
    name: "Red",
    style: "rgb(255,72,72)"
  },
  {
    name: "Green",
    style: "rgb(104,255,44)"
  }
]

let state = reactive({
  color: "rgba(255,255,255,1)",
  canvas: {} as HTMLCanvasElement,
  grid: true,
  size: 5,
})

onMounted(() => {
  initCanvas()
})

function setColor(color: string) {
  state.color = color
}

function toggleGrid(val: boolean) {
  state.grid = val
}


function clearCanvas() {
  let ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  if (!ctx) return;

  ctx.clearRect(0, 0, state.canvas.width, state.canvas.height)
}

function initCanvas() {
  let proto = document.getElementById("whiteboard-canvas")
  if (!proto) throw "Canvas 'whiteboard-canvas' not found"

  proto.style.width = '100%';
  proto.style.height = '100%';

  state.canvas = proto as HTMLCanvasElement
  // ...then set the internal size to match
  state.canvas.width = proto.offsetWidth;
  state.canvas.height = proto.offsetHeight;
  let ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  ctx.scale(1, 1)
  ctx.clearRect(0, 0, state.canvas.width, state.canvas.height)
}



function mouseDown(event: MouseEvent) {
  let ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  ctx.beginPath()
  let offsets = getOffsetSum(state.canvas)
  ctx.moveTo(event.pageX - offsets.x, event.pageY - offsets.y)
  ctx.lineTo(event.pageX - offsets.x, event.pageY - offsets.y)
}


function mousePlot(event: MouseEvent) {
  let ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  let offsets = getOffsetSum(state.canvas)
  ctx.lineTo(event.pageX - offsets.x, event.pageY - offsets.y)
  ctx.moveTo(event.pageX - offsets.x, event.pageY - offsets.y)
  ctx.strokeStyle = state.color
  ctx.globalAlpha = 1;
  ctx.lineWidth = state.size;
  ctx.lineJoin = "round";
  ctx.lineCap = "butt";
  ctx.globalCompositeOperation = "source-over";
  if (ctx.lineWidth >= 12) {
    ctx.lineWidth -= 0.1;
  }
  ctx.filter = "blur(1px)"
  ctx.imageSmoothingEnabled = true
  ctx.imageSmoothingQuality = "high"
  ctx.stroke()

}


function getOffsetSum(element: any) {
  var curleft = 0, curtop = 0;

  if (element.offsetParent) {
    do {
      curleft += element.offsetLeft;
      curtop += element.offsetTop;
      element = element.offsetParent;
    } while (element);
  }

  return {x: curleft, y: curtop};
}

function mouseStroke(event: MouseEvent) {
  let ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  ctx.closePath()


}

</script>
<template>
  <div class="blurs-10 d-flex flex-column gap-0 h-100 pb-2" @mousedown.stop>
    <div class="d-flex gap-1 w-100 mb-1">
      <Plot :cols="4" :rows="1">
        <Radio v-for="clr in colors" :active="state.color === clr.style" :fn="() => setColor(clr.style)"
               :title="clr.name"></Radio>
      </Plot>
      <Plot :cols="3" :rows="1" style="width: 8rem">
        <Radio :active="false" :fn="() => {}" :title="`${state.size}pt`" class="surface"></Radio>
        <Radio :active="false" :fn="() => state.size-=state.size<=1?0:1" title="-"></Radio>
        <Radio :active="false" :fn="() => state.size+=state.size>=12?0:1" title="+"></Radio>
      </Plot>
      <Plot :cols="2" :rows="1">
        <Radio :active="false" :fn="clearCanvas" title="Clear"></Radio>
        <Toggle :active="state.grid" :fn="() => toggleGrid(!state.grid)" title="Grid"></Toggle>
      </Plot>
    </div>


    <div class="element whiteboard-outer d-flex justify-content-center h-100">
      <canvas id="whiteboard-canvas" :class="`${state.grid?'whiteboard-canvas':''}`" @mousedown="mouseDown"
              @mousemove="mousePlot"
              @mouseup="mouseStroke"
              @mouseenter="mouseDown"
              @mouseleave="mouseStroke"
              @mouseout="mouseStroke"
              @mousedown.stop>
      </canvas>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.whiteboard-outer {
  height: 100%;
  width: 100%;
}

// Colors
$bg-color: rgba(255, 255, 255, 0.04);
$dot-color: rgba(255, 255, 255, 0.2);

// Dimensions
$dot-size: 28px;
$dot-space: 16px;

.whiteboard-canvas {

  background-color: $bg-color;
  opacity: 1;
  border-radius: 0.24rem;
  border: 1px solid $dot-color;
  background-image: radial-gradient($dot-color 0.980000000000000px, $bg-color 0.980000000000000px), radial-gradient($dot-color 0.98px, $bg-color 0.980000000000000px);
  background-size: $dot-size $dot-size;
  background-position: 0 0, calc($dot-size / 2) calc($dot-size / 2);
}

</style>