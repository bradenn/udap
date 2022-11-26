<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";
import {v4 as uuidv4} from "uuid";
import Ticks from "@/components/Ticks.vue";
import core from "@/core";

const props = defineProps<{
    min: number,
    max: number,
    step: number,
    unit?: string,
    name: string,
    value: number,
    tags?: string[],
    change?: (a: number) => void
    live?: boolean
    confirm?: boolean
}>()

let ticks = reactive({
    w: 0,
    h: 0,
    x: 0,
    y: 0,
})

let slider = reactive({
    w: 0,
    h: 0,
    x: 0,
    y: 0,
    stops: 0,
    uuid: uuidv4().toString(),
    last: 0
})

let thumb = reactive({
    engaged: false,
    w: 0,
    h: 0,
    x: 0,
    y: 0,
    previous: 0,
    position: 0,
})

const xStop = 51;

function initSlider() {

    const dom = document.getElementById(`track-${slider.uuid}`) as HTMLElement
    if (!dom) return;

    dom.addEventListener("mousemove", handleDrag)

    dom.addEventListener("mousedown", handleDrag)

    dom.addEventListener("mouseup", mouseUp)

    let bounds = dom.getBoundingClientRect()

    slider.x = bounds.x
    slider.y = bounds.y

    slider.w = bounds.width - xStop
    slider.h = bounds.height

    slider.stops = (props.max - props.min) / props.step + 1

    thumb.w = slider.w / slider.stops

    thumb.position = Math.floor(props.value / props.step)
    thumb.x = thumb.position * thumb.w

    ticks.x = 0
    ticks.w = slider.w


}

function handleDrag(ev: MouseEvent) {
    if (Date.now() - slider.last < 50) return
    slider.last = Date.now()
    let mouseX = ev.clientX - slider.x - 25

    let protoPos = Math.floor((mouseX / slider.w) * slider.stops)

    if (protoPos >= 0 && protoPos < slider.stops) {
        if (thumb.position == protoPos) return
        thumb.position = protoPos
        thumb.x = thumb.position * thumb.w
        haptics.tap(2, 1, 25)
    }

}

const haptics = core.haptics()

onMounted(() => {

    initSlider()

})


function mouseUp() {
    if (!props.confirm && props.change) {
        props.change(props.min + thumb.position * props.step)
    } else {

    }
}

</script>

<template>
    <div class="d-flex justify-content-center w-100">
        <div :id="`track-${slider.uuid}`" class="element p-0" style="width: 100%">

            <div class="d-flex justify-content-between lh-1 " style="padding-top:4px">
                <div class="label-c1 label-o3 label-w500 p-1 px-2 pb-0">{{ props.name }}</div>
                <div v-if="props.tags" class="label-c1 label-o3 label-w500 p-1 px-2 pb-0">
                    {{ props.tags[thumb.position] }}
                </div>
                <div v-else class="label-c1 label-o3 label-w500 p-1 px-2 pb-0">
                    {{ props.min + thumb.position * props.step }}{{ props.unit }}
                </div>
            </div>
            <Ticks v-if="ticks.w > 0" :active="thumb.position" :min="props.min" :series="5" :step="props.step"
                   :style="`width: ${ticks.w}px; left: ${ticks.x + (thumb.w+25)/2 - 10}px; position:relative; height:1.25rem;`"
                   :tags="tags" :ticks="slider.stops"></Ticks>
            <div class="shuttle-path subplot"></div>
            <div :style="`left: ${thumb.x}px;  position: relative; width: ${thumb.w + 50}px;`"
                 class="shuttle-cock subplot m-1 mt-1">
                <div :style="`left: ${(thumb.w+25)/2}px;`" class="shuttle"></div>
                <div class="shuttle-center">􀌇</div>


            </div>

        </div>
    </div>

</template>

<style lang="scss" scoped>
.shuttle-path {
  position: absolute;
  width: calc(100% - 16px);
  margin: 0.25rem;
  //background-color: #0a58ca;
  height: 2rem;
}

.shuttle-text {
  position: absolute;
  left: 8px;
}

.shuttle-cock {

  height: 2rem;
  width: 100px;
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

  top: -28px;
}

.shuttle-center {
  display: flex;
  justify-content: center;
  width: 100%;
  left: 0;
  transform: rotateZ(90deg);
  font-size: 1rem;
  //outline: 1px solid white;
  color: rgba(255, 255, 255, 0.1);
  font-weight: 300;
  font-family: "SF Pro", sans-serf, serif;
  //text-shadow: 0 0 4px rgba(0, 0, 0, 0.1);
  border-radius: 3px;
  position: absolute;


}

.slider {
}
</style>