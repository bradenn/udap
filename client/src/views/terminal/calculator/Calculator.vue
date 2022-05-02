<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";

const keys = [
  'C', '(', ')', '􀆛',
  'π', 'xⁿ', '􀓪', '^',
  7, 8, 9, '􀅿',
  4, 5, 6, '􀅾',
  1, 2, 3, '􀅽',
  0, '.', '(-)', '􀅼',
]

</script>
<template>
  <div class="d-flex flex-column gap-1">
    <div class="element p-2">
      <div>0.0</div>
    </div>

    <div class="d-flex flex-row justify-content-evenly gap-1">
      <div class="d-flex gap-1 flex-column" style="width: 12rem">
        <Plot :cols="4" :rows="6">

          <Radio v-for="item in keys" :active="false"
                 :title="`${item}`"/>
        </Plot>
        <Plot :cols="4" :rows="6">

          <Radio v-for="item in keys" :active="false"
                 :title="`${item}`"/>
        </Plot>
      </div>


      <Plot :cols="12" :rows="12" style="width: 100%">

        <Radio v-for="item in Array(12*12).keys()" :active="false"
               :title="`${item}`"/>
      </Plot>

    </div>
  </div>
</template>

<style lang="scss" scoped>

</style>