<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {reactive} from "vue";
import core from "@/core";

const props = defineProps<{
  horizontal?: boolean
}>()

let state = reactive({
  dragging: false,
  target: {} as HTMLElement,
  px: {
    velocity: 0,
  },

  timeout: 0,

  last: 0,

  se: {} as ScrollingElement,

  xDecelerate: 0,
  yDecelerate: 0,
  yTime: 0,

  xStart: 0,
  yStart: 0,

  xActual: 0,
  yActual: 0,

  xVelocity: 0,
  yVelocity: 0,

  xLastPoll: 0,
  yLastPoll: 0,

  xLastPos: 0,
  yLastPos: 0,

  xSilence: false,
  ySilence: false,

  lastTap: 0,
})

function dragStart(e: MouseEvent) {
  state.xVelocity = 0
  state.yVelocity = 0
  clearInterval(state.xDecelerate)
  clearInterval(state.yDecelerate)
  state.xLastPos = 0
  state.yLastPos = 0
  state.xLastPoll = 0
  state.yLastPoll = 0
  state.dragging = false
  state.yStart = e.clientY
  state.xStart = e.clientX
}

const haptics = core.haptics()

function drag(e: MouseEvent) {

  if (!state.dragging && Math.abs(state.yStart - e.clientY) > 5 || (props.horizontal && Math.abs(state.xStart - e.clientX) > 5)) {
    state.dragging = true
  }

  if (!state.dragging) return
  let element = e.currentTarget as HTMLElement
  if (!element) return;
  state.target = element

  // haptics(1, 1, 50)


  if (props.horizontal) {
    let poll = new Date().valueOf()
    let dt = state.xLastPoll - poll
    let dx = state.xLastPos - e.clientX

    element.scrollLeft += (state.xStart - e.clientX)
    state.xActual = element.scrollLeft
    state.xStart = e.clientX
    if (state.xActual <= 0 || state.xActual + element.clientWidth >= element.scrollWidth) {
      if (!state.xSilence) {
        state.xSilence = true
        // haptics.tap(2, 1, 25)
      }
    } else {
      state.xSilence = false
    }

    state.xVelocity = dx / dt
    state.xLastPos = e.clientX
    state.xLastPoll = poll
  } else {
    cancelAnimationFrame(afY)
    let poll = new Date().valueOf()
    let dt = state.yLastPoll - poll
    let dx = state.yLastPos - e.clientY

    element.scrollTop += state.yStart - e.clientY
    state.yActual = dx
    state.yStart = e.clientY


    state.yVelocity = dx / dt
    state.yLastPos = e.clientY
    state.yLastPoll = poll
    if (element.scrollTop <= 0 && !state.ySilence) {
      // haptics(1, 1, 100)
      state.ySilence = true
    } else if (element.scrollTop > 0) {
      state.ySilence = false
    }
  }

}

const FRICTION = 0.8;
const MASS = 10;

interface ScrollingElement {
  position: number;
  velocity: number;
}

function getDeltaPosition(s: { velocity: number; deceleration: number }, elapsedTime: number): number {
  // Calculate the new velocity after applying deceleration
  const newVelocity = s.velocity + (s.deceleration * elapsedTime);
  state.yVelocity = newVelocity
  // Calculate the delta position using the average velocity during the elapsed time
  return (s.velocity + newVelocity) / 2 * elapsedTime;
}

function decelerate(time: number, currentPosition: number, previousPosition: number): number {
  let acceleration: number;

  if (currentPosition > previousPosition) {
    // Moving upwards
    acceleration = (FRICTION * MASS) / time;
  } else {
    // Moving downwards
    acceleration = -(FRICTION * MASS) / time;
  }

  // Modify the yVelocity of the global state object
  state.yVelocity -= acceleration * time;

  // Calculate and return the change in position
  return (state.yVelocity * time) + (0.5 * acceleration * Math.pow(time, 2));
}

let afX: number = 0;
let afY: number = 0;
const decayRate = 0.99;

function animateX(): void {
  state.xVelocity = state.xVelocity * 0.80
  state.target.scrollLeft -= state.xVelocity * 20

  if (Math.abs(state.xVelocity) <= 0.01 && state.xStart <= 0.1) {
    state.xVelocity = 0
    cancelAnimationFrame(afX);
  }
  afX = requestAnimationFrame(animateX);

}

function animateY(): void {

  state.yTime += 1000 / 60

  state.yVelocity *= 0.95
  state.target.scrollTop -= state.yVelocity * 10

  if (Math.abs(state.yVelocity) <= 0.01) {
    state.yVelocity = 0;
    state.yTime = 0
    cancelAnimationFrame(afY);
  }

  afY = requestAnimationFrame(animateY);
}


// state.xVelocity = state.xVelocity * 0.80
// state.target.scrollLeft -= state.xVelocity * 20
// state.xStart = state.xStart * 0.8
//
// if (Math.abs(state.xVelocity) <= 0.01 && state.xStart <= 0.1) {
//   state.xVelocity = 0
//   cancelAnimationFrame(af)
// }

function mouseDown() {
  state.lastTap = Date.now().valueOf()
  haptics.tap(1, 1, 50)
}

function dragStop(e: MouseEvent) {
  state.yStart = 0
  state.xStart = 0
  let element = e.currentTarget as HTMLElement
  if (!element) return;
  state.target = element
  state.timeout = 0
  if (props.horizontal) {

    animateX()
  } else {
    animateY()

  }

  // state.xDecelerate = setInterval(() => {
  //   state.xVelocity = state.xVelocity * 0.80
  //   element.scrollLeft -= state.xVelocity * 20
  //   state.xStart = state.xStart * 0.8
  //
  //   if (Math.abs(state.xVelocity) <= 0.01 && state.xStart <= 0.1) {
  //     state.xVelocity = 0
  //     clearInterval(state.xDecelerate)
  //   }
  // }, 16)

  // state.yDecelerate = setInterval(() => {
  //
  //   state.yVelocity = state.yVelocity * 0.965
  //   element.scrollTop -= state.yVelocity * 12
  //
  //   if (Math.abs(state.yVelocity) <= 0.0001) {
  //     state.yVelocity = 0;
  //     clearInterval(state.yDecelerate)
  //   }
  // }, 16)


  // setTimeout(() => {
  //   state.dragging = false
  // }, 250)

}

function click(e: MouseEvent) {
  if (state.dragging) {
    e.preventDefault()
    e.stopPropagation()
  }
}

</script>
<template>

  <div v-on:mousedown="dragStart" v-on:mousemove="drag" v-on:mouseup="dragStop"
       class="scroll-box" style="padding-right: 0.25rem"
       v-on:click.capture="click">
    <slot></slot>

  </div>
</template>


<style scoped>
.scroll-box {
  backface-visibility: hidden;
  perspective: 1000;

  -webkit-backface-visibility: hidden;
  -webkit-perspective: 1000;

  transform: translate3d(0, 0, 0);
}
</style>