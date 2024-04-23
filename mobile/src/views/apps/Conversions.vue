<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import {reactive} from "vue";
import core from "udap-ui/core";
import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";


const remote = core.remote();

let state = reactive({
  show: false,
  context: [0, 0, 0, 0, 0, 0] as number[],
  results: 0 as number,
  conversion: {} as Conversion
})

interface Parameter {
  name: string
  symbol: string
  description: string
}

interface Conversion {
  name: string
  parameters: Parameter[]
  func: (num: number[]) => number

}

const conversions: Conversion[] = [
  {
    name: "&deg;C to &deg;F",
    parameters: [
      {
        name: "Degrees Celsius",
        symbol: "&deg; C",
        description: "Temperature in degree Celsius to convert",
      }
    ] as Parameter[],
    func: (num: number[]): number => {
      return (num[0] * (9.0 / 5.0)) + 32.0
    }
  },
  {
    name: "&deg;F to &deg;C",
    parameters: [
      {
        name: "Degrees Fahrenheit",
        symbol: "&deg; F",
        description: "Temperature in degree Fahrenheit to convert",
      }
    ],
    func: (num: number[]): number => {
      return (num[0] * (9.0 / 5.0)) + 32.0
    }
  },
  {
    name: "Warp to c",
    parameters: [
      {
        name: "Warp Factor",
        symbol: "warp",
        description: "Temperature in degree Fahrenheit to convert",
      }
    ],
    func: (num: number[]): number => {
      let wf = num[0]
      if (wf < 9) {
        return Math.pow(wf, (10 / 3))
      }
      let a = 0.03684678
      let n = 1.791275
      return Math.pow(wf, (((10 / 3) + a * Math.pow((-Math.log2(10 - wf)), n))))
    }
  }
]

function openModal(conversion: Conversion) {
  state.show = true

  state.conversion = conversion
  console.log(conversion)
}

function closeModal() {
  state.show = false
}

function valueChange() {

  state.results = state.conversion.func(state.context)
}

function update(et: Event, param: Parameter) {
  let e = et as InputEvent
  let val = e.target as HTMLInputElement
  if (!val) return
  console.log("hell")
  state.context[state.conversion.parameters.indexOf(param)] = parseFloat(val.value) || 0;
  valueChange()
}

</script>

<template>
  <div class="d-flex gap-1 flex-column">


    <Element>
      <List>
        <!--        <ElementHeader title="Sci-Fi"></ElementHeader>-->
        <ElementLink v-for="v in conversions" :cb="() => openModal(v)" :title="v.name" icon=""></ElementLink>

      </List>
    </Element>

    <div v-if="state.show" class="modal p-5 px-3" @touchstart="() => state.show?closeModal():{}">
      <div class="d-flex align-items-center justify-content-center flex-column">
        <div class="label-c0 label-o5 label-w600 mb-4" v-html="state.conversion.name"></div>
        <List class="flex-fill w-100">

          <Element v-for="param in state.conversion.parameters" foreground style="width: 100%">
            <div class="d-flex flex-row justify-content-between align-items-center w-100 px-0">
              <div class="d-flex align-items-center px-3">
                <!--                <div v-html="param.symbol" class="label-o5 label-c3 label-w700" style="width: 1.6rem"></div>-->
                <div class="label-o6 label-c3 label-w700 px-2">
                  {{ param.name }}
                </div>
              </div>
              <div class="d-flex align-items-center">
                <input :class="`input-number-active`"
                       class="input-number " inputmode="decimal"
                       type="number"
                       @input="(e) => update(e,param)" @touchstart.stop/>
                <div class="label-o5 label-c3 label-w700" style="width: 1.5rem; margin-left: 0.5rem;"
                     v-html="param.symbol"></div>
              </div>
              <div>

              </div>
            </div>
          </Element>
          <Element foreground>
            <div class="label-o6 label-c3 label-w700 px-2">
              {{ state.results }}
            </div>
          </Element>
        </List>

      </div>


    </div>

  </div>
</template>

<style lang="scss" scoped>
.input-number-active {
  background-color: rgba(255, 255, 255, 0.1);
}

.input-number:focus {
  outline: 1px solid rgba(255, 255, 255, 0.05);
}

.input-number {
  background-color: rgba(255, 255, 255, 0.05);
  //background-color: transparent;
  color: rgba(255, 255, 255, 0.6);
  border: none;
  border-radius: 0.25rem;
  //box-shadow: none !important;

  //padding-right: 1.5rem !important;
  padding: 0.8rem;
  width: 8rem;
  height: calc(3.25rem - 0.625rem);

  //text-align: right;
}

.sf-icon {
  color: hsl(0, 0%, 50%);
}

.modal {
  position: absolute;
  top: -0.5px;
  left: -0.5px;
  width: 100vw;
  height: 100vh;
  background-color: rgba(255, 255, 255, 0.01);
  box-shadow: inset 0 0 16px 8px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}
</style>