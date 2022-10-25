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
  <div :style="``" class="element p-2 px-1 pt-1 toast"
       style="width: 18rem; height: 2.5rem">
    <div class="d-flex align-items-start px-1">
      <div class="w-100 d-flex flex-column justify-content-center">
        <div class="d-flex justify-content-between lh-1 pt-1 ">
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
  height: 42px;

  margin: 8px 14px 8px 8px;

  background-color: rgba(25, 135, 84, 0.53);
}

@keyframes pulse {
  0% {
    outline: 10px solid white;
  }
  50% {
    outline: 80px solid white;
  }
  100% {
    outline: 10px solid white;
  }
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
  animation: pulse 400ms ease-out infinite;
}
</style>