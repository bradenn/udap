<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import {reactive} from "vue";

const keys = [
  'C', '(', ')', '􀆛',
  'π', 'xⁿ', '􀓪', 'sin',
  '7', '8', '9', '􀅿',
  '4', '5', '6', '􀅾',
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
    name: "-",
    key: "-",
    type: "operand",
    args: 0
  }
]

let state = reactive({
  cursor: 0,
  registers: {
    values: [""] as string[],
    operators: new Map<number, any>()
  }
})

function isOperator(input: string) {
  return operators.find(operator => operator.key == input)
}

function handleOperator(operator: string) {
  if (state.registers.operators.has(state.cursor) || state.registers.values[state.cursor] !== "") {
    state.cursor++;
    state.registers.values[state.cursor] = ""
  }
  state.registers.operators.set(state.cursor, operators.find(o => o.key == operator))


}

function addValue(input: string) {

  if (isOperator(input)) {
    handleOperator(input)
    return
  }

  switch (input) {
    case 'C':
      state.registers.values = []
      state.registers.operators = new Map<number, any>()
      state.cursor = 0
      return
    case '􀆛':
      if (state.cursor < 0) return
      if (state.registers.values[state.cursor] === "") {
        if (state.cursor > 0) {
          state.registers.operators.delete(state.cursor)
        }
        return
      }
      state.registers.values[state.cursor] = state.registers.values[state.cursor].slice(0, state.registers.values[state.cursor].length - 1)
      break
    default:
      state.registers.values[state.cursor] += input
  }

}

</script>
<template>
  <div class="d-flex flex-column gap-1">
    <Plot :cols="2" :rows="2" style="width: 13rem">
      <Radio v-for="item in inputModifiers" :active="false"
             :fn="() => addValue(item)" :title="`${item}`"/>
    </Plot>
    <Plot :cols="2" :rows="1">
      <div class="number-line-group">

      </div>
      <div class="d-flex flex-row align-items-center gap-1 px-2 label-sm label-r">
        <div v-for="(value, k) in state.registers.values" v-if="state.registers.values.length > 0">
          <div class="d-flex gap-1">
            <div v-if="state.registers.operators.has(k)" class="label-w400 label-o4">{{
                state.registers.operators.get(k).name
              }}<span v-if="state.registers.operators.get(k).type === 'function'">(</span>

            </div>
            <div class="label-w400 label-o4" @click="state.cursor = k">
              <span :style="`${state.cursor === k?'border-bottom: 2px solid #b05b20;':''}`">
                {{ value }}
              </span>

            </div>
            <div v-if="state.registers.operators.has(k)">
              <span v-if="state.registers.operators.get(k).type === 'function'" class="label-w400 label-o4">)</span>
            </div>
          </div>
        </div>
        <div v-else>
          <div class="label-w400 label-o4">{{ state.registers.values[state.cursor] }}(</div>
        </div>

      </div>
    </Plot>
    <div>
      <Plot :cols="4" :rows="6" style="width: 12rem;">

        <Radio v-for="item in keys" :active="false"
               :fn="() => addValue(item)" :title="`${item}`"/>
      </Plot>
    </div>

  </div>
</template>

<style lang="scss" scoped>
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