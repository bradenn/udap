<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {reactive} from "vue";
import Subplot from "@/components/plot/Subplot.vue";

interface Props {
  selected: string
}

const props = defineProps<Props>()

const state = reactive({
  active: false
})

function dragScroll(e: DragEvent) {
  let element = e.target
  console.log(e)
}


</script>

<template>
  <div :class="state.active?'select-dropdown-show':''" class="element select-dropdown"
       tabindex="0" @focusout="(e) => {state.active = !state.active}">
    <div class="d-flex align-items-center justify-content-between lh-1" style="margin-top: -0.25rem; height: 1.8rem">
      <div class="label-c1  label-o4 label-w500 px-1">Zones</div>
      <div class="label-c2 label-w500 label-r label-o3 px-1 text-link"></div>
    </div>
    <div class="select-scroll" @drag="dragScroll">
      <slot></slot>
    </div>
    <div class="h-sep my-1"></div>
    <div class="d-flex gap-1 justify-content-between" style="height: 1.75rem">
      <div></div>
      <div class="d-flex">
        <Subplot :active="false" :fn="() => {}" name="Close"></Subplot>
        <Subplot :active="true" :fn="() => {}" name="Edit"></Subplot>
      </div>
    </div>
  </div>
  <div :class="state.active?'select-active':''" class="element select" v-on:click="state.active = !state.active">
    <div class="subplot d-flex">
      <div>
        <div>{{ props.selected }}</div>
      </div>
      <div class="p-1">􀆊</div>
    </div>

  </div>

</template>


<style lang="scss">
.select-scroll {
  width: 100%;
  overflow: scroll;
  height: 10rem;
  max-height: 10rem;
  padding-right: 0.25rem;
}

.label-micro {
  font-size: 0.4rem;
  height: 12px;
  top: 10px;
  color: rgba(255, 255, 255, 0.34)
}

.select-active {
  z-index: 200 !important;
}

.select {
  position: relative !important;

  display: flex;
  flex-direction: column;
  gap: 0.25rem;

  overflow: visible !important;

}

.select > .subplot {

  height: 1.8rem;

}

.select-scroll > * {
  height: 1.8rem;
}

.select-dropdown-show {
  display: block !important;
}

.select-dropdown {
  display: none;
  position: absolute !important;
  width: 12rem;

  margin-left: 11.25rem;

  z-index: 22;

}

.text-link {

}


.subplot {
  justify-content: center;
}


</style>