<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import {reactive} from "vue";

const numbers = [
  'C', '(', ')', '􀆛',
  'π', 'xⁿ', '􀓪', ',',
  7, 8, 9, '􀅼',
  4, 5, 6, '􀅽',
  1, 2, 3, '􀅾',
  'E', 0, '.', '􀅿',
]


let state = reactive({
  line: "",
  value: 0
})


function calculatorType(value: any) {
  let button: HTMLAudioElement = new Audio('/sound/selection.mp3');
  button.play()
  let form = `${value}`.replace("π", "Math.PI").replace("􀓪", "Math.sqrt(").replace("xⁿ", "Math.pow(").replace("􀅿", "/").replace('􀅾', '*').replace('􀅼', '+').replace('􀅽', '-').replace('􀆀', '=')
  switch (form) {
    case 'C':
      state.line = ""
      state.value = 0
      return
    case '􀆛':
      if (state.line === "") {
        state.value = 0
        return
      }
      state.line = state.line.slice(0, state.line.length - 1)
      state.value = eval(state.line) || 0
      break
    default:
      state.line += form
      state.value = eval(state.line) || 0
  }


}

</script>
<template>
  <div class="element" style="border-radius: 0.65rem">
    <Plot :cols="1" :rows="1" class="mb-1" style="">
      <div class="d-flex flex-column align-items-end p-1">
        <div class="calculator-result mb-0 pb-0">{{ Math.round(state.value * 10000000000) / 10000000000 }}</div>
        <div class="label-o2 label-c1 mb-1">{{ state.line }}&nbsp;</div>
      </div>
    </Plot>
    <Plot :cols="4" :rows="6" class="" style="">
      <Radio v-for="i in numbers.keys()" :active="false" :fn="() => calculatorType(numbers[i])"
             :title="`${numbers[i]}`"/>
    </Plot>
  </div>
</template>

<style lang="scss" scoped>

.calculator-result {
  font-size: 1.4rem;
  line-height: 1.8rem;
  font-family: "IBM Plex Sans", sans-serif;
}
</style>