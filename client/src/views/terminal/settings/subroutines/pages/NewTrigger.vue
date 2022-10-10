<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, reactive} from "vue";
import type {Remote, Trigger} from "@/types";
import Keyboard from "@/components/Keyboard.vue";
import triggerService from "@/services/triggerService";

const router = useRouter()

const remote = inject("remote") as Remote


const fields = [{
  name: "Name",
  description: "The human-readable display name of this trigger",
  value: "",
  type: "text"
}, {
  name: "Description",
  description: "A brief description of the trigger activation source",
  value: "",
  type: "text"
}, {
  name: "Type",
  description: "Who or what will activate this trigger",
  value: "module",
  type: "radio",
  values: ["module", "system", "manual"]
}]

const state = reactive({
  trigger: {
    name: "",
    description: "",
    type: "",
  } as Trigger,
  fields: fields,
  field: 0,
  ready: false,
})

function goBack() {
  router.push("/terminal/settings/subroutines/create")
}

function selectTrigger(id: string) {
  if (state.trigger === id) {
    state.trigger = ""
  } else {
    state.trigger = id
  }
}

function setName(name: string) {
  state.trigger.name = name
  state.ready = !isDisabled()
}

function setDescription(description: string) {
  state.trigger.description = description
  state.ready = !isDisabled()
}

function isDisabled(): boolean {
  let t = state.trigger
  return t.name == "" || t.description == "" || t.type == ""
}

interface InputProps {
  name: string
  value: string
  description: string
  icon?: string
  type?: string
}

let input = inject("input") as (props: InputProps, cb: (a: string) => void) => void

function update(value: string) {
  if (value === "{bksp}") {
    state.fields[state.field].value = state.fields[state.field].value.substring(0, state.fields[state.field].value.length - 1)
  } else if (value === "{enter}") {
    state.field++;
  } else if (value === "{space}") {
    state.fields[state.field].value += " "
  } else {
    state.fields[state.field].value += value
  }
}

function setField(field: number) {
  state.field = field
}

function createTrigger() {
  state.trigger.name = state.fields[0].value
  state.trigger.description = state.fields[1].value
  state.trigger.type = state.fields[2].value

  triggerService.createTrigger(state.trigger).then(res => {
    goBack()
  }).catch(err => {

  })
}

</script>

<template>
  <div class="generic-grid">
    <div>

      <div class="d-flex gap-1">
        <div v-for="i in Array(state.field).keys()"
             class="element p-1 mt-1 flex-fill d-flex justify-content-between align-items-center"
             @click="() => setField(i)">
          <div class="label-c1 label-w700 label-o5 px-2 lh-2">{{ state.fields[i].name }}</div>
          <div class="label-c2 label-w500 label-o3 px-2 lh-1">{{ state.fields[i].value }}</div>
        </div>
      </div>
      <div class="element p-2 mt-1">
        <div class="label-c2 label-w500 label-o3 px-0 lh-1 mb-2 text-accent"
             @click="() => createTrigger()">􀆉 Back
        </div>
        <div class="label-md label-w700 label-o5 px-1 lh-1 mb-1">{{ state.fields[state.field].name }}</div>
        <div class="d-flex mb-1">
          <div v-if="state.fields[state.field].type === 'text'" class="subplot p-2 mt-0 label-c2 w-100">
            <div class="text-input w-100">
              <div class="label-xs" v-text="state.fields[state.field].value"></div>
              <div class="cursor-blink"></div>
              <div class="flex-fill"></div>
            </div>
          </div>
          <div v-else class="d-flex w-100 gap-1" style="height: 1.75rem">
            <div v-for="value in state.fields[state.field].values"
                 :class="`${value===state.fields[state.field].value?'text-accent':''}`"
                 class="subplot flex-grow-1 justify-content-center text-capitalize"
                 @click="state.fields[state.field].value = value">
              {{ value }}
            </div>
          </div>
        </div>
        <div class="label-c2 label-o4 px-1">{{ state.fields[state.field].description }}</div>
      </div>
      <div v-if="state.field >= 2" class="pt-1 ">
        <div class="element plot-1x4" style="" @click="() => createTrigger()">
          <div class="subplot justify-content-center">Create Trigger</div>
        </div>
      </div>
    </div>
    <div v-if="state.fields[state.field].type === 'text'" class="kb-input">
      <Keyboard :input="update" class="keyboard-location" keySet=""
                keyboardClass="simple-keyboard"
      ></Keyboard>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.kb-input {
  position: absolute;
  bottom: 6.5rem;
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
  /*grid-column-gap: 0.25rem;*/
  /*grid-row-gap: 0.25rem;*/
  /*grid-template-rows: repeat(1, 1fr);*/
  /*grid-template-columns: repeat(3, 1fr);*/
}
</style>