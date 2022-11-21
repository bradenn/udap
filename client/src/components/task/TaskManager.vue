<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import type {Task} from "@/types";
import {TaskType} from "@/types";
import {inject, onMounted, reactive, watchEffect} from "vue";
import Keyboard from "@/components/Keyboard.vue";
import FixedScroll from "@/components/scroll/FixedScroll.vue";
import core from "@/core";
import Item from "@/components/element/Item.vue";
import type {Haptics} from "@/haptics";

interface Tasks {
    title: string,
    tasks: Task[],
    onComplete: (tasks: Task[]) => void
}

const state = reactive({
    tasks: [] as Task[],
    current: {} as Task,
    textbox: false,
    loaded: false,
})

watchEffect(() => {

    return props.tasks
})

onMounted(() => {
    if (props.tasks.length > 0) {
        state.tasks = props.tasks
        if (state.current.value == null) {
            state.current = state.tasks[0]
            state.loaded = true
        }
    }
})

const remote = core.remote()

let props = defineProps<Tasks>()

const haptics = inject("haptics") as Haptics

function selectRadio(value: any, preview: string) {
    haptics.tap(2, 3, 100)
    state.current.value = value
    state.current.preview = preview
}

function isListSelected(value: any): boolean {
    let list = state.current.value as any[]
    return list.includes(value)
}

function selectList(value: any) {

    let list = state.current.value as any[]
    if (list.includes(value)) {
        list = list.filter(l => l != value)
        haptics.tap(2, 1, 100)
        haptics.tap(1, 2, 50)
    } else {
        list.push(value)

        haptics.tap(2, 3, 100)
    }
    state.current.value = list
    state.current.preview = `${list.length} item${list.length == 1 ? '' : 's'}`
}

function nextTask() {

    let current = state.tasks.indexOf(state.current)
    if (current + 1 >= state.tasks.length) {
        props.onComplete(state.tasks)
        return
    }
    state.current = state.tasks[current + 1]
    updateLocal()
}

function updateLocal() {
    state.tasks = state.tasks.map(t => t.title === state.current.title ? state.current : t)
}

function selectTask(task: Task) {
    updateLocal()
    state.current = task
    haptics.tap(1, 1, 50)
}

function calcDuration(min: number) {
    return `${Math.floor(min / 60)}h ${min % 60}m `
}

function selectDuration(min: number) {
    state.current.value = Math.max(state.current.value + min, 0)
    state.current.preview = calcDuration(state.current.value)
}

function enterKey(key: string) {
    let def = state.current.value as string
    if (key === '{bksp}') {
        def = def.substring(0, def.length - 1);
    } else if (key === '{enter}') {
        nextTask()
    } else if (key === '{space}') {
        def += " "
    } else {
        def += key
    }

    updateLocal()
    state.current.value = def;
    state.current.preview = def;

}

function selectIcon(icon: string) {
    haptics.tap(2, 3, 100)
    state.current.value = icon
    state.current.preview = icon
}

function isIconSelected(icon: string) {
    return state.current.value === icon
}

function mouseDown(m: MouseEvent) {
    let elem = m.target as HTMLElement
    elem.className = "subplot-down subplot"
}


</script>

<template>
    <div class="d-flex gap-1 w-100" style="height:100%">

        <div class="d-flex flex-column gap-1 w-100">
            <div class="element d-flex flex-row justify-content-around gap-1 p-0">
                <div v-for="task in state.tasks" :id="task.title" class=" w-100 ">
                    <div :class="`${task.title === state.current.title?'':'subplot-inline'}`"
                         class="subplot d-flex justify-content-center gap-1"
                         style="height: 2rem; border-radius: 0.5rem !important;"
                         @click="selectTask(task)">
                        <div
                                class="px-1 lh-1 d-flex flex-column align-items-center gap-0">
                            <div :class="`${task.title === state.current.title?'text-accent':''}`"
                                 class="label-c2 label-o4 label-w500 lh-1">{{ task.title }}
                            </div>
                            <div class="label-c3 label-o3 label-w500 lh-1">{{ task.preview }}</div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="element">
                <div v-if="state.loaded" class="d-flex justify-content-between align-items-center">
                    <div class="label-md label-w700 label-o5 px-1 lh-1 mb-2 mt-1">{{ state.current.title }}</div>
                    <div class="label-c2 label-w400 label-o3 px-2 lh-1 mb-2 mt-1">{{ state.current.preview }}</div>
                </div>
                <FixedScroll class="" style="max-height: 12rem; overflow-y: scroll;">
                    <div>
                        <div v-if="state.current.type === TaskType.String">
                            <div class="d-flex mb-1 align-items-center justify-content-center">
                                <div :class="`${state.textbox?'border border-accent':''}`"
                                     class="subplot p-1 px-2 mt-0 label-c2 flex-grow-1"
                                     @click="() => state.textbox = !state.textbox">
                                    <div class="text-input">
                                        <div class="label-xxs label-o3" style="white-space: pre;"
                                             v-text="state.current.value"></div>
                                        <div class="cursor-blink"></div>
                                        <div class="flex-fill"></div>
                                    </div>
                                </div>
                            </div>

                        </div>
                        <div v-else-if="state.current.type === TaskType.Number"
                             class=" d-flex gap-1 w-100 justify-content-center">
                            <div class="d-flex gap-1">
                                <div class="subplot p-1 px-2" @click="() => selectDuration(-30)">- 30m</div>
                                <div class="subplot p-1 px-2" @click="() => selectDuration(-15)">- 15m</div>
                                <div class="subplot p-1 px-2" @click="() => selectDuration(-1)">- 1m</div>
                            </div>
                            <div class="subplot p-1 px-3">{{ calcDuration(state.current.value) }}</div>
                            <div class="d-flex gap-1">
                                <div class="subplot p-1 px-2" @click="()=> selectDuration(1)">+ 1m</div>
                                <div class="subplot p-1 px-2" @click="()=> selectDuration(15)">+ 15m</div>
                                <div class="subplot p-1 px-2" @click="()=> selectDuration(30)">+ 30m</div>
                            </div>

                        </div>
                        <div v-else-if="state.current.type === TaskType.Radio">
                            <div class="selection-grid px-1 pt-1">
                                <div v-for="option in state.current.options"
                                     :key="option.title"
                                     :class="`${option.value === state.current.value?'accent-selected':''}`"
                                     class="subplot p-1 d-flex flex-row align-items-start gap-0"
                                     @click="() => selectRadio(option.value, option.title)">
                                    <div class="d-flex flex-column p-1">
                                        <div class="lh-1 label-c2 label-o4">{{ option.title }}</div>
                                        <div class="lh-1 label-c3 label-o3">{{ option.description }}</div>
                                    </div>
                                    <div v-if="option.value === state.current.value"
                                         class="label-c2 label-o4 label-w500 text-accent lh-1 p-1">􀷙</div>
                                    <div v-else
                                         class="label-c2 label-o2 label-w400 lh-1 p-1">􀓞</div>
                                </div>
                            </div>
                        </div>
                        <div v-else-if="state.current.type === TaskType.List">
                            <div class="selection-grid px-1 pt-1 pb-2">
                                <Item v-for="option in state.current.options"
                                      :class="`${isListSelected(option.value)?'accent-selected':''}`"
                                      :description="option.description" :selected="isListSelected(option.value)"
                                      :title="`${option.title}`"
                                      :toggle="true"
                                      @click="() => selectList(option.value)"></Item>

                            </div>
                        </div>
                        <div v-else-if="state.current.type === TaskType.Icon">
                            <div class="selection-grid-xs px-1 pt-1">
                                <div v-for="option in state.current.options"
                                     :key="option.title"
                                     :class="`${isIconSelected(option.value)?'accent-selected':''}`"
                                     class="subplot p-1 d-flex flex-row align-items-start gap-0"
                                     @click="() => selectIcon(option.value)">
                                    <div class="d-flex flex-column p-1">
                                        <div class="lh-1 label-c1 label-o4">{{ option.title }}</div>

                                    </div>
                                    <div v-if="isIconSelected(option.value)"
                                         class="label-c2 label-o4 label-w500 text-accent lh-1 p-1">􀷙</div>
                                    <div v-else
                                         class="label-c2 label-o2 label-w400 lh-1 p-1">􀓞</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </FixedScroll>
                <div class="d-flex justify-content-between align-items-center mt-1">
                    <div class="label-c2 label-w400 label-o3 px-2 lh-1 ">{{ state.current.description }}</div>
                    <div class="subplot  d-flex justify-content-center" style="height:1.5rem; min-width:3.5rem;"
                         @click="() => nextTask()">
                        <div v-if="state.tasks.indexOf(state.current) < state.tasks.length-1"
                             class="label-c2 label-w500 label-o3 px-1 lh-1 text-accent">Next 􀯻</div>
                        <div v-else class="label-c2 label-w500 label-o3 px-1 lh-1 text-accent">Finish 􀯻</div>
                    </div>
                </div>

            </div>

        </div>
    </div>
    <Keyboard v-if="state.current.type === TaskType.String && state.textbox" :input="enterKey" keySet="d"
              keyboardClass="simple-keyboard"></Keyboard>

</template>


<style lang="scss">

.accent-selected:before {
  border: 2px solid rgba(255, 149, 0, 0.4) !important;
  content: ' ';
  position: absolute;
  border-radius: inherit;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.tactile-button:active {
  animation: accent-animation 10ms forwards linear;
}

@keyframes accent-animation {
  0% {
    transform: scale3d(1, 1, 1);
  }
  100% {
    transform: scale3d(0.99, 0.99, 0.99);
  }
}

.label-shimmer {
  mix-blend-mode: lighten;
}

.nav-bar {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  width: 100%;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.selection-grid {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  width: 100%;
  grid-template-rows: repeat(5, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

.selection-grid-xs {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  width: 100%;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(6, 1fr);
}

//.nav-bar > div {
//  width: 100%;
//  outline: 1px solid white;
//}

.popup-nav {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  width: 100%;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.nav-area {
  display: flex;
  align-items: center;
}

.pane-pop-up {
  z-index: 23 !important;
  position: absolute;
  width: 100%;
  height: 100%;
}


</style>