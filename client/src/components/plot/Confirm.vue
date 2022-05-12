<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {reactive} from "vue";

interface Toggle {
  title?: string
  icon?: string
  active?: boolean
  fn?: () => void
}


let props = defineProps<Toggle>()

let state = reactive({
  intent: false,
  remaining: 6500,
  countdown: 5000,
  interval: 0,
  aborted: false
})


function abort(e: MouseEvent) {
  props.fn ? props.fn() : ''
  reset()
}

function reset() {
  state.intent = false
  clearInterval(state.interval)
  state.interval = 0
  state.remaining = 6500
}

function initialIntention(e: MouseEvent) {
  if (state.intent) {
    abort(e)
    return
  }
  state.intent = true
  state.interval = setInterval(countdown, 100)

}

function countdown() {
  state.remaining -= 100;
  if (state.remaining <= 0) {
    state.aborted = false
    reset()
  } else if (state.remaining <= 1500) {
    state.aborted = true
  }
}

</script>

<template>
  <div class="subplot label-o3 justify-content-center align-items-center label-r w-100 " @mousedown="initialIntention">
    <div v-if="!state.intent" class="d-flex align-items-center align-content-center gap-1 load-in">
      <div v-if="props.icon" class="label-o3 label-c2" v-html="props.icon"></div>
      <div v-if="props.title" class="label-o4">{{ props.title }}</div>
    </div>
    <div v-else-if="state.aborted" class="load-in">
      <div class="label-o3 label-c2">􀣔 Operation Timed Out</div>
    </div>
    <div v-else class="pulse load-in">
      􀬂 Tap to Authorize
    </div>
  </div>
</template>

<style scoped>
.pulse {
  color: rgba(25, 135, 84, 1);
  transition: all 500ms ease;
}

.load-in {
  animation: send 50ms ease forwards;
}

@keyframes send {
  0% {
    transform: scale3d(0.75, 0.75, 0.75);
  }
  100% {
    transform: scale3d(1, 1, 1);
  }
}

.subplot:active {
  animation: click 100ms ease forwards;
}

@keyframes click {
  0%, 100% {
    transform: scale(1.005);
  }
  50% {
    transform: scale(0.98);
  }
}

.subplot.abort {
  animation: shadows 3s infinite;
}

.abort-countdown {
  bottom: 8px;
  position: absolute;
  width: 3rem;
  display: flex;
  justify-content: start;

}

@keyframes shadows {
  0% {
    box-shadow: 5px 5px 8px 2px white;
  }
  25% {
    box-shadow: 5px 0px 8px 2px white;
  }
  50% {
    box-shadow: 0px 0px 8px 2px white;
  }
  75% {
    box-shadow: 0px 5px 8px 2px white;
  }
  100% {
    box-shadow: 5px 5px 8px 2px white;
  }
}

.abort-remaining {
  position: relative;
  align-items: start;
  height: 4px;
  border-radius: 8px;
  transition: width 100ms linear;
  background-color: #ff1037;
}
</style>