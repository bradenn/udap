<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, reactive} from "vue";
import PaneMenuToggle from "@/components/pane/PaneMenuToggle.vue";
import type {Macro, Remote, SubRoutine} from "@/types";
import Tasks from "@/components/task/Tasks.vue";
import Task from "@/components/task/Task.vue";
import Keyboard from "@/components/Keyboard.vue";
import subroutineService from "@/services/subroutineService";

const router = useRouter()

const remote = inject("remote") as Remote

interface InputField {
  name: string,
  description: string,
  value: any,
  type: string,
  values?: any[]
}

let fields: InputField[] = [
  {
    name: "Trigger",
    description: "What event should run this subroutine",
    value: "",
    type: "triggerId"
  },
  {
    name: "Macros",
    description: "What macros should run in this subroutine",
    value: "",
    type: "macros",
    values: []
  }
]

const state = reactive({
  subroutine: {} as SubRoutine,
  fields: fields as InputField[],
  name: "",
  trigger: {
    id: "",
    name: ""
  },
  macros: {
    macros: [] as Macro[],
    preview: ""
  },
  field: 0,
})

function goBack() {
  router.push("/terminal/settings/subroutines")
}

function createTrigger() {
  router.push("/terminal/settings/subroutines/trigger")
}

function selectTrigger(id: string) {
  if (state.trigger.id === id) {
    state.trigger.id = ""
    state.trigger.name = ""
  } else {
    state.trigger.id = id
    state.trigger.name = remote.triggers.find(t => t.id === id)?.name || ""
  }
}

function selectMacro(id: string) {
  if (state.macros.macros.find(m => m.id === id)) {
    state.macros.macros = state.macros.macros.filter(m => m.id !== id)
  } else {
    let mc = remote.macros.find(m => m.id === id)
    if (!mc) return
    state.macros.macros.push(mc)
    state.macros.preview = state.macros.macros.map(m => m.name).join(", ")
  }
}

function setField(field: number) {
  state.field = field
}

function update(value: string) {
  if (value === "{bksp}") {
    state.name = state.name.substring(0, state.name.length - 1)
  } else if (value === "{enter}") {
    state.field++;
  } else if (value === "{space}") {
    state.name += " "
  } else {
    state.name += value
  }
}

function createSubroutine() {
  subroutineService.createSubroutine({
    description: state.name,
    macros: state.macros.macros,
    triggerId: state.trigger.id
  } as SubRoutine).then(res => {
    console.log(res)
  }).catch(err => {
    console.log(err)
  })

}

</script>

<template>
  <div class="generic-grid">
    <Tasks title="New Subroutine">
      <Task :active="state.field === 0"
            :index="0"
            :next="setField"
            :preview="state.trigger.name"
            :value="state.trigger.id"
            description="What event should run this subroutine"
            title="Trigger">
        <PaneMenuToggle v-for="item in remote.triggers" :active="state.trigger.id === item.id"
                        :fn="() => selectTrigger(item.id)"
                        :subtext="item.description"
                        :title="item.name"></PaneMenuToggle>
      </Task>
      <Task :active="state.field === 1"
            :index="1"
            :next="setField"
            :preview="state.macros.preview"
            :value="state.macros.preview"
            description="What macros should this subroutine call"
            title="Macros">
        <PaneMenuToggle v-for="item in remote.macros"
                        :active="!!state.macros.macros.find(m => m.id === item.id)"
                        :fn="() => selectMacro(item.id)"
                        :subtext="item.description"
                        :title="item.name"></PaneMenuToggle>
      </Task>
      <Task :active="state.field === 2"
            :index="2"
            :next="setField"
            :value="state.name"
            description="The display name for this subroutine"
            title="Name">
        <div class="subplot p-2 mb-1 label-c2 " style="margin-left: 0.25rem;">
          <div class="text-input">
            <div class="label-xs" v-text="state.name"></div>
            <div class="cursor-blink"></div>
            <div class="flex-fill"></div>
          </div>
        </div>
      </Task>
      <div class="element">
        <div class="subplot justify-content-center label-w500 label-c1" style="height: 1.8rem"
             @click="createSubroutine">Create Subroutine</div>
      </div>
    </Tasks>
    <Tasks title="New Subroutine">
      <Task :active="state.field === 0"
            :index="0"
            :next="setField"
            :preview="state.trigger.name"
            :value="state.trigger.id"
            description="What event should run this subroutine"
            title="Trigger">
        <PaneMenuToggle v-for="item in remote.triggers" :active="state.trigger.id === item.id"
                        :fn="() => selectTrigger(item.id)"
                        :subtext="item.description"
                        :title="item.name"></PaneMenuToggle>
      </Task>
      <Task :active="state.field === 1"
            :index="1"
            :next="setField"
            :preview="state.macros.preview"
            :value="state.macros.preview"
            description="What macros should this subroutine call"
            title="Macros">
        <PaneMenuToggle v-for="item in remote.macros"
                        :active="!!state.macros.macros.find(m => m.id === item.id)"
                        :fn="() => selectMacro(item.id)"
                        :subtext="item.description"
                        :title="item.name"></PaneMenuToggle>
      </Task>
      <Task :active="state.field === 2"
            :index="2"
            :next="setField"
            :value="state.name"
            description="The display name for this subroutine"
            title="Name">
        <div class="subplot p-2 mb-1 label-c2 " style="margin-left: 0.25rem;">
          <div class="text-input">
            <div class="label-xs" v-text="state.name"></div>
            <div class="cursor-blink"></div>
            <div class="flex-fill"></div>
          </div>
        </div>
      </Task>
      <div class="element">
        <div class="subplot justify-content-center label-w500 label-c1" style="height: 1.8rem"
             @click="createSubroutine">Create Subroutine</div>
      </div>
    </Tasks>


    <div v-if="state.field === 2" class="kb-input">
      <Keyboard :input="update" class="keyboard-location" keySet="" keyboardClass="simple-keyboard"></Keyboard>
    </div>
  </div>
</template>

<style scoped>


.kb-input {
  position: absolute;
  bottom: 4rem;
}

.generic-grid > div {
  width: 18rem;
}

.generic-grid {
  display: flex;
  justify-content: center;
  /*width: 100%;*/
  /*height: 100%;*/
  /*display: grid;*/
  grid-column-gap: 0.25rem;
  /*grid-row-gap: 0.25rem;*/
  /*grid-template-rows: repeat(1, 1fr);*/
  /*grid-template-columns: repeat(3, 1fr);*/
}
</style>