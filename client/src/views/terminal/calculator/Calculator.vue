<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import {reactive, watchEffect} from "vue";
import Draw from "@/components/Draw.vue";
import core from "@/core";
import type {Attribute} from "@/types";

const remote = core.remote()

const keys = [
  'C', '(', ')', '􀆛',
  '7', '8', '9', '/',
  '4', '5', '6', '*',
  '1', '2', '3', '-',
  '0', '.', '(-)', '+',
]

const inputModifiers = [
  '', '1',
  '􀰌', '􀄫'
]

const operators = [
  {
    name: "sin",
    key: "sin",
    type: "function",
    args: 1
  },
  {
    name: "+",
    key: "+",
    type: "operand",
    args: 0
  },
  {
    name: "*",
    key: "*",
    type: "operand",
    args: 0
  },
  {
    name: "/",
    key: "/",
    type: "operand",
    args: 0
  }
]

let state = reactive({
  cursor: 0,
  value: 0,
  response: {} as Attribute,
  registers: {
    values: "",
    operators: new Map<number, any>()
  }
})

function init() {

}


function isOperator(input: string) {
  return operators.find(operator => operator.key == input)
}

function calculate() {
  state.value = Math.round(eval(state.registers.values) * 100) / 100
}

watchEffect(() => {
  state.response = remote.attributes.find(a => a.key == "result") as Attribute
})

function addValue(input: string) {

  switch (input) {
    case 'C':
      state.registers.values = ""
      state.cursor = 0
      state.value = 0
      return
    case '􀆛':
      if (state.cursor < 0) {
        state.value = 0
        return;
      }
      state.registers.values = state.registers.values.substring(0, state.registers.values.length - 1)
      break
    default:
      state.registers.values += input

  }
  calculate()
}

</script>
<template>
  <div class="d-flex justify-content-center gap-1 mt-4">

    <div class="d-flex flex-column gap-1">
      <div>
        <Plot :cols="4" :rows="1" style="width: 12rem; margin-bottom: 0.25rem;">
          <div
              class="grid-line d-flex justify-content-start flex-column w-100 align-items-end"
              style="height: 2.8rem;">
            <div class="label-xl label-w500 label-o5">{{ state.value }}</div>
            <div class="label-c2 label-o3">{{ state.registers.values }}</div>
          </div>
        </Plot>
        <Plot :cols="4" :rows="5" style="width: 12rem;">

          <div v-if="false"
               class="d-flex flex-row align-items-center gap-1 px-2 label-sm label-r grid-line">
            <div v-for="(value, k) in state.registers.values"
                 v-if="state.registers.values.length > 0">
              <div class="d-flex gap-1">
                <div v-if="state.registers.operators.has(k)"
                     class="label-w400 label-o4">{{
                    state.registers.operators.get(k).name
                  }}<span
                      v-if="state.registers.operators.get(k).type === 'function'">(</span>

                </div>
                <div class="label-w400 label-o4" @click="state.cursor = k">
              <span
                  :style="`${state.cursor === k?'border-bottom: 2px solid #b05b20;':''}`">
                {{ value }}
              </span>

                </div>
                <div v-if="state.registers.operators.has(k)">
                  <span
                      v-if="state.registers.operators.get(k).type === 'function'"
                      class="label-w400 label-o4">)</span>
                </div>
              </div>
            </div>
            <div v-else>
              <div class="label-w400 label-o4">
                {{ state.registers.values[state.cursor] }}(
              </div>
            </div>

          </div>
          <Radio v-for="item in keys" :active="false"
                 :fn="() => addValue(item)" :title="`${item}`"/>
        </Plot>
      </div>

    </div>
    <div class="element d-flex align-items-center justify-content-center"
         style="width: 13rem">
      <div>
        {{ state.response.value }}
      </div>
      <Draw style="width: 12rem; height: 12rem;"></Draw>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.grid-line {
  grid-column: span 4;
  grid-row: span 1;
}

.number-line-group {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 0.375rem;
  //border-radius: 0.375rem !important;
  font-family: "SF Pro Rounded", sans-serif !important;
  color: rgba(255, 255, 255, 0.5);
  background-color: rgba(34, 38, 45, 0);
  border-radius: 0.3rem;
  border: 1px solid rgba(255, 255, 255, 0.05);

  font-weight: 500;
  font-size: 0.625rem;
}

.cursor {
  color: #b05b20 !important;
}
</style>