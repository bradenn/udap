<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>


import {onMounted, reactive, watchEffect} from "vue";
import moment from "moment";
import Countdown from "@/components/Countdown.vue";


interface ToastProps {
    title: string
    message: string
    severity: number
    time: number
    index?: number
}

const props = defineProps<ToastProps>()


const state = reactive({
    start: 0,
    time: "",
})

let iconMap = new Map<number, any>(
    [
        [
            0, {icon: "􀅴", color: "10,132,255,0.6"}
        ],
        [
            1, {icon: "􀁢", color: "48,209,88,0.6"}
        ],
        [
            2, {icon: "􀇾", color: "255,149,0,0.6"}
        ],
        [
            3, {icon: "􀘯", color: "255, 69, 58,0.6"}
        ]
    ]
)

onMounted(() => {
    state.start = props.time
})

watchEffect(() => {
    state.time = moment(props.time).fromNow()
})


</script>

<template>
    <div :class="`${props.time <= 300 && props.time !== -1 ? 'toast-dissolve' : ''}`" :style="`z-index: ${index || 2};`"
         class="element toast"
         style="width: 18rem; height: 2.25rem; ">
        <div class="d-flex align-items-start px-1">
            <div class="w-100 d-flex flex-column justify-content-start align-items-start">
                <div class="d-flex justify-content-between align-content-center align-items-center w-100">
                    <div class="d-flex align-items-center align-content-center justify-content-start ">
                        <div :style="`color: rgba(${iconMap.get(props.severity).color})`"
                             class="label-o2 label-c1 label-w500"
                             style="padding-right: 0.225rem">
                            {{ iconMap.get(props.severity).icon }}
                        </div>
                        <div class="label-c2 label-c1 label-o5 label-w600 lh-1">

                            {{ props.title }}
                        </div>
                    </div>
                    <div v-if="props.time > 0" class="lh-1">
                        <Countdown :percent="props.time/state.start"></Countdown>
                    </div>
                </div>
                <div class="label-c3 label-o3 label-w500 lh-1 px-3" style="letter-spacing: 0.3px">{{
                        props.message
                    }}
                </div>
            </div>
        </div>


    </div>
</template>

<style lang="scss" scoped>


.toast-dissolve.toast.element {

  z-index: 10 !important;
  animation: toast-dissolves 125ms linear forwards !important;
}

@keyframes toast-dissolves {
  0% {
    transform: scale(1) translateY(0);
    filter: blur(0);
  }
  75% {
    transform: scale(0.92) translateY(0);
    opacity: 0.9;
    filter: blur(12px);
  }
  85% {
    transform: scale(0.60) translateY(-10px);
    opacity: 0.1;
    filter: blur(18px);
  }
  100% {
    opacity: 0.0;
    transform: translateY(-15px);
    filter: blur(32px);

  }
}

.element.toast {
  animation: toast-resolves 125ms linear forwards !important;
}

@keyframes toast-resolves {
  0% {
    opacity: 0.0;
    transform: translateY(-20px);
    filter: blur(36px);

  }
  15% {
    transform: scale(0.60) translateY(-20px);
    opacity: 0.1;
    filter: blur(24px);
  }
  25% {
    transform: scale(0.92) translateY(0);
    opacity: 0.9;
    filter: blur(12px);
  }
  100% {
    transform: scale(1) translateY(0);
    filter: blur(0);

  }
}


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
}
</style>