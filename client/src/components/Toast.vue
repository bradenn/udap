<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>


import {reactive, watchEffect} from "vue";
import moment from "moment";

interface ToastProps {
  title: string
  message: string
  severity: number
  time: number
}

const props = defineProps<ToastProps>()


const state = reactive({
  time: ""
})

watchEffect(() => {
  state.time = moment(props.time).fromNow()
})


</script>

<template>
  <div class="element p-1 toast" style="width: 18rem; height: 2.5rem">
    <div class="d-flex">
      <div
          :style="`background-color: rgba(${props.severity === 0?'25, 135, 84':props.severity === 1?'255, 149, 0':'255,69,58'}, 0.53);`"
          class="status-marker"></div>
      <div class="w-100 d-flex flex-column justify-content-center">
        <div class="d-flex justify-content-between">
          <div class="label-c2 label-o5 label-w600  d-flex">
            {{ props.title }}
          </div>
          <div class="label-c2 label-o2 label-w500 px-1">{{ props.time / 1000 }}s</div>
        </div>
        <div class="label-c2 label-o3 label-w500 lh-1">{{ props.message }}</div>
      </div>
    </div>


  </div>
</template>

<style lang="scss" scoped>
.status-marker {
  width: 4px !important;
  border-radius: 4px;
  height: 28px;

  margin: 8px 14px 8px 8px;

  background-color: rgba(25, 135, 84, 0.53);
}

.module-icon {
  align-items: center;
  display: flex;
  justify-content: center;
  align-content: center;
  background-color: rgba(255, 255, 255, 0.08);
  border-radius: 0.2rem;
  margin-right: 0.25rem;
  padding-top: 0.25px;
  font-size: 0.4rem;
  width: 0.75rem;
  height: 0.75rem;
}


.toast {
  position: absolute;
  left: calc(50% - 8rem);
}
</style>