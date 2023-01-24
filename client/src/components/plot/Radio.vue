<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {reactive} from "vue";

interface Toggle {
  title?: string
  icon?: string
  sf?: string
  active: boolean
  disabled?: boolean
  fn: () => void
}


let props = defineProps<Toggle>()
let state = reactive({
  interval: 0,
  holding: false,
})

function handle(e: MouseEvent) {
  props.fn()
  state.holding = true
  setTimeout(() => {
    if (state.holding) {
      if (state.interval == 0) {
        state.interval = setInterval(props.fn, 125)
      } else {
        up()
      }
    } else {
      if (state.interval != 0) {
        up()
      }
    }
  }, 375)


}


function up() {

  state.holding = false
  if (state.interval != 0) clearInterval(state.interval)
  state.interval = 0;
}
</script>

<template>
  <div
      :class="`${props.active?'active':props.disabled===true?'disabled-radio':''}`"
      class="radio subplot d-flex justify-content-center"
      @click="(e) => props.fn()">
    <div><span v-if="props.icon"><i :class="`fa-${props.icon}`"
                                    class="fa-solid "></i>&nbsp;</span>{{
        props.title
      }}
      <div v-if="props.sf" class="label-o3 label-c2" v-html="props.sf"></div>
    </div>
    <slot></slot>
  </div>
</template>

<style scoped>
.radio.subplot:active {
  animation: click 100ms ease forwards;
}

.radio.subplot.disabled-radio {
  opacity: 0.5 !important;
}

@keyframes click {
  0% {
    transform: scale(1.0);
  }
  25% {
    transform: scale(0.98);
  }
  100% {
    transform: scale(1);
  }
}
</style>