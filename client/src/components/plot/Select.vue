<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {reactive} from "vue";

interface Props {
  selected: string
}

const props = defineProps<Props>()

const state = reactive({
  active: false
})


</script>

<template>
  <div class="element select">
    <div class="subplot" v-on:click="state.active = !state.active">{{ props.selected }}</div>
    <slot v-if="state.active" class="select-dropdown"></slot>
  </div>
</template>


<style>
.select {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;

}

.select > .subplot {
  height: 1.8rem;

}

.select-dropdown > .subplot {
  height: 1.8rem;

}

.select-dropdown {
  position: absolute !important;

  top: 0;
  left: 0;
  height: 100px;
}


.subplot:active {
  animation: click 100ms ease forwards !important;
}

.subplot:hover {
  animation: click 100ms ease forwards !important;
}

.subplot {
  justify-content: center;
}


</style>