<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {reactive} from "vue";
import {v4 as uuidv4} from "uuid";
import core from "@/core";


const props = defineProps<{
  active: boolean,
  name?: string,
  unit?: string,
  change?: (a: boolean) => void
}>()

let state = reactive({
  active: props.active,
  uuid: uuidv4().toString(),
})

const haptics = core.haptics()


function apply() {
  if (props.change) {
    props.change(state.active);
  }
}

function toggle(pos: boolean) {
  console.log("Hello")
  state.active = pos
  haptics.tap(2, 1, 100)
  apply()
}


</script>

<template>
  <div class="element w-100 p-1">
    <div class="d-flex justify-content-between">
      <div class="label-c1 label-o3 label-w500 p-1 px-1 pb-1 pt-0">{{ props.name }}</div>
      <div class="label-c1 label-o3 label-w500 p-1 px-1 pb-1 pt-0">
        {{ state.active ? "ON" : "OFF" }}
      </div>
    </div>
    <div class="switch-path subplot" style="z-index: -1 !important;"></div>
    <div class="switch-spots">
      <div :class="`${!state.active?'subplot switch-space-active switch-right':''}`" class="switch-space"
           @click="(e) => toggle(false)">
        􀥥
      </div>
      <div :class="`${state.active?'subplot switch-space-active switch-left':''}`" class="switch-space"
           @click="(e) => toggle(true)">
        􀥤
      </div>
    </div>
  </div>

</template>

<style lang="scss" scoped>

.switch-space-active {
  color: rgba(255, 149, 0, 0.7) !important;

  box-shadow: 0 0 20px 1px rgba(0, 0, 0, 0.1), inset 0 0 2px 1px rgba(255, 255, 255, 0.05);
}

.switch-right {
  //animation: switch-move-right 150ms forwards ease-out;
}

@keyframes switch-move-right {
  0% {
    transform: translateX(-16px) scale3d(0.96, 0.96, 1);
  }
  100% {
    transform: translateX(0px) scale3d(1, 1, 1);
  }
}

.switch-left {
  //animation: switch-move-left 150ms forwards ease;
}

@keyframes switch-move-left {
  0% {
    transform: translateX(16px);
  }
  100% {
    transform: translateX(0px);
  }
}

.switch-space {
  display: flex;
  justify-content: center;
  font-weight: 700;
  align-items: center;
  font-size: 0.5rem;
  height: 1.8rem;
  color: rgba(255, 255, 255, 0.4);
}

.switch-spots {
  display: grid;
  grid-column-gap: 0.5rem;
  grid-row-gap: 0.5rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(2, 1fr);
}

.switch-path {
  position: absolute;
  width: calc(100% - 0.5rem);
  left: 8px;

  height: 1.8rem;
  //background-color: #0a58ca;
}

.shuttle-text {
  position: absolute;
  left: 8px;
}

.shuttle-cock {

  height: 2rem;
  //width: 100px;
  box-shadow: 0 0 20px 1px rgba(0, 0, 0, 0.1), inset 0 0 2px 1px rgba(255, 255, 255, 0.05);

}


.arrow-down {

}

.shuttle {
  width: 3px;
  height: 8px;
  background-color: rgba(255, 255, 255, 0.25);
  border-radius: 3px;
  position: relative;

  top: -28px
}

.shuttle-center {
  width: 3px;
  height: 8px;
  background-color: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
  position: absolute;


}

.slider {
}
</style>