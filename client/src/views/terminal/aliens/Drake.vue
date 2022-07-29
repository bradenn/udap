<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Plot from "@/components/plot/Plot.vue";
import {onMounted, reactive} from "vue";
import Radio from "@/components/plot/Radio.vue";

export interface Page {
  pageid: number;
  ns: number;
  title: string;
  extract: string;
}

export interface Normalized {
  from: string;
  to: string;
}

export interface Query {
  normalized: Normalized[];
  pages: Map<string, Page>;
}

export interface WikiQuery {
  batchcomplete: string;
  query: Query;
}

interface Template {
  latex: string
  desc: string
  default: number
  value: number
  style?: string
  step: number
  unit: string
  name: string
}

let defaults = [
  {
    latex: 'N',
    desc: "Number of civilizations with which humans could communicate",
    default: 0.2,
    unit: "civilizations"
  },
  {
    latex: 'R_{*}',
    desc: "Mean number of new stars are formed every year in the Milky Way galaxy",
    default: 2,
    step: 0.5,
    unit: "stars per year",
    name: "Star Formation"
  },
  {
    latex: 'f_P',
    desc: "Percent of newly formed stars that harbour planets",
    default: 0.95,
    step: 0.05,
    style: "percent",
    unit: "have planets",
    name: "Planet Formation"
  },
  {
    latex: 'n_e',
    desc: "Mean number of planets that could support life per star with planets",
    default: 3,
    step: 0.5,
    unit: "support life",
    name: "Planet Habitability"
  },
  {
    latex: 'f_l',
    desc: "Fraction of life-supporting planets that eventually develop life",
    default: 0.5,
    step: 0.1,
    style: "percent",
    unit: "develop life",
    name: "Abiogenesis"
  },
  {
    latex: 'f_i',
    desc: "Fraction of planets with life where life develops intelligence",
    default: 0.05,
    style: "percent",
    step: 0.01,
    unit: "gain intelligence",
    name: "Evolution of Intelligence"
  },
  {
    latex: 'f_c',
    desc: "Fraction of intelligent civilizations that develop and use communication",
    default: 0.2,
    style: "percent",
    step: 0.1,
    unit: "communicate",
    name: "Willingness to Communicate"
  }
  ,
  {
    latex: 'L',
    desc: "Mean length of time a civilization will communicate with other planets",
    default: 1000,
    step: 250,
    unit: "years",
    name: "Propagation Period"
  }
] as Template[]

let state = reactive({
  params: [] as Template[],
  civilizations: 0 as number,
  loaded: false,
  data: {} as WikiQuery
})

onMounted(() => {
  reset()
  fetchWiki()
})


function reset() {
  state.params = defaults.slice(1)
  state.params.forEach(p => {
    p.value = p.default
  })
  state.loaded = true
  compute()
}

function resetValue(tmpl: Template) {
  tmpl.value = tmpl.default
  compute()
}

function compute() {
  state.civilizations = Math.round(state.params.map(m => m.value).reduce((a, b) => a * b) * 100) / 100
}

function changeValue(tmpl: Template, sign: number) {
  if (sign == 0) {
    tmpl.value = tmpl.default
    return
  }
  if (tmpl.value + sign * tmpl.step > 0) {
    tmpl.value += sign * tmpl.step
    tmpl.value = Math.round(tmpl.value * 100) / 100
  }
  compute()
}

function fetchWiki() {
  let url = "https://en.wikipedia.org/w/api.php?action=query&prop=extracts&exlimit=2&titles=drake+equation&format=json&origin=*"
  // axios.get(url, {}).then(res => {
  //   state.data = res.data as WikiQuery
  // }).catch(err => {
  //
  // })
}

</script>

<template>
  <div v-if="state.loaded" class="d-flex gap-1">
    <!--    <div v-for="page in state.data.query.pages">-->
    <!--      <div v-html="page.extract"></div>-->
    <!--    </div>-->
    <div class="d-flex flex-column gap-2">
      <div class="d-flex justify-content-start">
        <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-rocket fa-fw`"></i></div>
        <div class="label-w500 opacity-100 label-xxl px-2 label-r">The Drake Equation</div>
        <div class="flex-fill"></div>

        <Plot :cols="1" :rows="1" small style="width: 6rem;">
          <Radio :active="false" :fn="() => reset()"
                 title="Reset"></Radio>
        </Plot>
      </div>

      <div class="equation-container">


        <Plot :cols="1" :rows="1">
          <div class="d-flex flex-column gap-0 px-1 justify-content-center align-items-center">
            <div class="label-md px-1 label-o5">
              <math-jax :latex="`${state.params.map(m => m.latex).join('\\cdot ')} = N`"/>
            </div>
            <div class="label-c1 px-1 label-o4">
              <math-jax :latex="`${state.params.map(m => m.value).join(' \\cdot ')} = ${state.civilizations}`"/>
            </div>
          </div>

        </Plot>
        <Plot :cols="1" :rows="1">
          <div class="d-flex flex-column justify-content-center align-items-start px-2">
            <div class=" label-sm pt-1 label-o5">
              <math-jax :latex="`N = ${state.civilizations}`"/>
            </div>
            <div class="label-c2 label-o4 label-r">Approximately
              {{ state.civilizations }} civilizations in the Milky Way
            </div>
          </div>
        </Plot>
      </div>

      <div class="option-container">
        <Plot v-for="tmpl in state.params" :alt-fn="() => resetValue(tmpl)" :cols="1" :rows="1" :title="tmpl.name"
              alt="Reset">
          <div class="d-flex justify-content-start px-1 align-items-start">
            <div class="d-flex pt-1 label-o5" style="width: 2rem;margin-inline: 0.5rem;">
              <math-jax :latex="tmpl.latex" style="width: 4rem;"/>
            </div>
            <div class="mt-1" style="width: 100%">
              <Plot :cols="3" :rows="1" class="mb-1" style="width: 100%; box-shadow: none;">
                <div class="subplot d-flex align-items-center justify-content-center label-o5 label-r label-w500"
                     v-on:click="(e) => changeValue(tmpl, -1)">
                  􀅽
                </div>
                <div class="d-flex align-items-center justify-content-center flex-column">
                  <div class="label-sm label-o4 label-r label-w500 lh-1">
                    {{ (((tmpl.style || "") === "percent") ? `${Math.round(tmpl.value * 100.0)} %` : `${tmpl.value}`) }}
                  </div>
                  <div class="label-c2 label-o3 label-r label-w500 lh-1">{{ tmpl.unit }}</div>
                </div>
                <div class="subplot d-flex align-items-center justify-content-center label-o5 label-r label-w500"
                     v-on:click="(e) => changeValue(tmpl, 1)">
                  􀅼
                </div>
              </Plot>
              <div class="label-c1 label-o4 label-r pb-1 mx-1" style="">{{ tmpl.desc }}</div>
            </div>
          </div>
        </Plot>
      </div>
    </div>


  </div>
</template>

<style lang="scss" scoped>
.option-container {
  display: grid;
  grid-column-gap: 0.5rem;
  grid-row-gap: 0.5rem;
  grid-auto-flow: row;
  grid-template-rows: repeat(4, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.equation-container {
  display: grid;
  grid-column-gap: 0.5rem;
  grid-row-gap: 0.5rem;
  grid-auto-flow: row;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}
</style>
