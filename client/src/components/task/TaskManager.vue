<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import type {Task} from "@/types";
import {TaskType} from "@/types";
import {inject, reactive, watchEffect} from "vue";
import PaneMenuToggle from "@/components/pane/PaneMenuToggle.vue";
import Keyboard from "@/components/Keyboard.vue";
import Scroll from "@/components/scroll/Scroll.vue";
import type {Haptics} from "@/views/terminal/haptics";
import Menu from "@/components/menu/Menu.vue";
import MenuLink from "@/components/menu/MenuLink.vue";

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
  if (props.tasks.length > 0) {
    state.tasks = props.tasks
    if (state.tasks.length > 0) {
      state.current = state.tasks[0]
      state.loaded = true
    }
  }
  return props.tasks
})

let props = defineProps<Tasks>()

const haptics = inject("haptics") as Haptics

function selectRadio(value: any, preview: string) {
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
  } else {
    list.push(value)
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

}

function selectTask(task: Task) {
  state.current = task
  haptics.tap(1, 1, 50)
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

  state.current.value = def;

}


</script>

<template>
  <div class="d-flex gap-1" style="height:100%">
    <div>

      <Menu class="mt-1" style="width: 15rem">
        <div class="label-lg label-w700 label-o5 lh-1 p-1 px-2 ">Macro</div>
        <MenuLink v-for="task in state.tasks" :active="task.title === state.current.title"
                  :subtext="task.preview" :title="task.title" @click="selectTask(task)"></MenuLink>
      </Menu>
    </div>
    <div class="d-flex flex-column gap-1 w-100">
      <div class="element d-flex flex-row justify-content-around gap-1 p-0">
      </div>
      <div class="element">
        <div v-if="state.loaded" class="d-flex justify-content-between align-items-center">
          <div class="label-md label-w700 label-o5 px-1 lh-1 mb-2 mt-1">{{ state.current.title }}</div>
          <div class="label-c2 label-w400 label-o3 px-2 lh-1 mb-2 mt-1">{{ state.current.preview }}</div>
        </div>
        <Scroll class="" style="max-height: 12rem; overflow-y: scroll;">
          <div>
            <div v-if="state.current.type === TaskType.String">
              <div class="d-flex mb-1 align-items-center justify-content-center">
                <div :class="`${state.textbox?'border border-accent':''}`"
                     class="subplot p-1 px-2 mt-0 label-c2 flex-grow-1"
                     @click="() => state.textbox = !state.textbox">
                  <div class="text-input">
                    <div class="label-xxs label-o3" style="white-space: pre;" v-text="state.current.value"></div>
                    <div class="cursor-blink"></div>
                    <div class="flex-fill"></div>
                  </div>
                </div>
              </div>

            </div>
            <div v-if="state.current.type === TaskType.Number">
              Number
            </div>
            <div v-else-if="state.current.type === TaskType.Radio">

              <PaneMenuToggle v-for="option in state.current.options" :active="state.current.value === option.value"
                              :fn="() => selectRadio(option.value, option.title)"
                              :subtext="option.description"
                              :title="option.title">

              </PaneMenuToggle>

            </div>
            <div v-if="state.current.type === TaskType.List">
              <PaneMenuToggle v-for="option in state.current.options" :active="isListSelected(option.value)"
                              :fn="() => selectList(option.value)"
                              :subtext="option.description"
                              :title="option.title">

              </PaneMenuToggle>
            </div>
          </div>
        </Scroll>
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


<style lang="scss" scoped>

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