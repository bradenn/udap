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

  <div class="side-app">
    <div class="element">
      <div class="label-o4 label-c1 label-w500 p-1 px-2 pb-0">Calculator</div>
      <div class="d-flex w-100 justify-content-end px-2">
        <div class="d-flex flex-column align-items-end">
          <div class="calculator-result mb-0 pb-0">{{ Math.round(state.value * 10000000000) / 10000000000 }}</div>
          <div class="label-o2 label-c1 mb-1">{{ state.line }}&nbsp;</div>
        </div>
      </div>
      <Plot :cols="4" :rows="6" class="" style="">
        <Radio v-for="i in numbers.keys()" :active="false" :fn="() => calculatorType(numbers[i])"
               :title="`${numbers[i]}`"/>
      </Plot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.side-app {
  position: absolute;
  right: 1rem;
  height: calc(100% - 2rem);
  width: 13rem;

  z-index: 22 !important;
}

.calculator-result {
  font-size: 1.4rem;
  line-height: 1.8rem;
  font-family: "IBM Plex Sans", sans-serif;
}
</style>