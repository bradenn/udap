<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import {onMounted, reactive} from "vue";

let state = reactive({
  keyset: ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'] as string[],
  entry: ['', '', '', '', '', ''],
  cursor: 0,
  timeout: 0
})

function sortRandom(a: string, b: string): number {
  return 1 - Math.random() * 2
}

onMounted(() => {
  randomize()
})

function randomize() {
  state.keyset = state.keyset.sort(sortRandom)
}

function enterKey(a: string) {
  if (state.cursor < 6) {
    state.entry[state.cursor] = a
    state.cursor++
  }


}

function deleteKey() {
  if (state.cursor == 0) {
    return
  }

  state.cursor--
  state.entry[state.cursor] = ""
}

</script>

<template>
  <div class="d-flex flex-column gap-1">
    <Plot :cols="6" :rows="1" style="width: 13rem">
      <div v-for="number in Array(6).keys()" class="subplot d-flex justify-content-center">
        <div class="label-md label-r">
          <div v-if="state.cursor >= number">
            {{ state.entry[number] }}
          </div>
          <div v-else>
            &nbsp;
          </div>
        </div>
      </div>
    </Plot>

    <Plot :cols="4" :rows="5" style="width: 13rem">
      <div></div>
      <div></div>
      <Radio :active="false" :fn="randomize" icon="shuffle"></Radio>
      <Radio :active="false" :fn="deleteKey" icon="delete-left"></Radio>
      <div v-for="number in state.keyset" :key="state.keyset.indexOf(number)"
           class="subplot d-flex justify-content-center"
           @mousedown="() => enterKey(number)">
        <div class="label-md label-r">{{ number }}</div>
      </div>
    </Plot>

  </div>
</template>

<style scoped>
.keypad-char {

}
</style>