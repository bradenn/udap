<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import List from "udap-ui/components/List.vue";
import Element from "udap-ui/components/Element.vue";
import Title from "udap-ui/components/Title.vue";

import * as patterns from "udap-ui/vendor/hero-patterns"
import core from "@/core";
import {onBeforeMount, reactive} from "vue";


const preferences = core.preferences()

function selectPattern(pattern: string) {

  preferences.pattern.svg = getPattern(pattern)
  preferences.pattern.name = pattern

}

let state = reactive({
  loaded: false,
  patternList: Object.keys(patterns),
  scaleOptions: [1, 2, 3, 4] as number[],
  colors: [] as string[],
  backgrounds: [] as string[],
  opacities: [0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1],
  blurOptions: [2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 30]
})

function hslToHex(h: number, s: number, l: number): string {
  l /= 100;
  const a = s * Math.min(l, 1 - l) / 100;
  const f = (n: number): string => {
    const k = (n + h / 30) % 12;
    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
    return Math.round(255 * color).toString(16).padStart(2, '0');   // convert to Hex and prefix "0" if needed
  };
  return `#${f(0)}${f(8)}${f(4)}`;
}

function selectColor(accent: string) {
  preferences.accent = accent
  preferences.pattern.svg = getPattern(preferences.pattern.name)
}


onBeforeMount(() => {
  let numItems = 64
  for (let i = 0; i < numItems; i++) {
    state.colors.push(`${hslToHex((360 / numItems) * i, 50, 30)}`)
  }
  for (let i = 0; i < 48; i++) {
    state.backgrounds.push(`${hslToHex(0, 0, (50 / 48) * i)}`)
  }
  for (let i = 1; i <= 10; i++) {
    state.scaleOptions.push(i * 5)
  }
  state.loaded = true
})

function selectBackgroundColor(color: string) {
  preferences.background = color
}

function selectScale(scale: number) {
  preferences.pattern.scale = scale
}

function selectBlur(blur: number) {
  preferences.blur = blur
}

function selectOrientation(landscape: boolean) {
  preferences.landscape = landscape
}

function getPattern(pattern: string): string {
  //@ts-ignore
  return patterns[pattern](preferences.accent, preferences.pattern.opacity)
}

function selectOpacity(opacity: number) {
  preferences.pattern.opacity = opacity
}

</script>

<template>

  <List>

    <List>
      <Element v-if="false" :foreground="true">
        <div class="">
          <List :row="true" scroll-x>
            <Element v-for="w in state.scaleOptions" :accent="preferences.pattern.scale == w" :cb="() => selectScale(w)"
                     :foreground="true"

                     class="py-4 d-flex justify-content-center align-items-center"
                     style="min-width: 4rem; height: 2rem">
              {{ w }}
            </Element>

          </List>
        </div>
      </Element>

      <Element :foreground="true">
        <Title title="Element Blur"></Title>
        <div class="">
          <List :row="true" scroll-x>
            <Element v-for="b in state.blurOptions" v-if="state.loaded" :cb="() => selectBlur(b)"
                     :foreground="true"
                     :style="`${preferences.blur== b?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''} `"
                     class="py-4 d-flex justify-content-center align-items-center"
                     style="min-width: 4rem; height: 2rem"
                     surface>
              {{ b }}px
            </Element>


          </List>
        </div>
      </Element>
      <Element :foreground="true">
        <Title title="Orientation"></Title>
        <div class="">
          <List :row="true" scroll-x>
            <Element v-if="state.loaded" :cb="() => selectOrientation(true)"
                     :foreground="true"
                     :style="`${preferences.landscape?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''} `"
                     class="py-4 d-flex justify-content-center align-items-center"
                     style="min-width: 4rem; height: 2rem"
                     surface>
              Landscape
            </Element>
            <Element v-if="state.loaded" :cb="() => selectOrientation(false)"
                     :foreground="true"
                     :style="`${!preferences.landscape?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''} `"
                     class="py-4 d-flex justify-content-center align-items-center"
                     style="min-width: 4rem; height: 2rem"
                     surface>
              Portrait
            </Element>

          </List>
        </div>
      </Element>
    </List>


    <List :row="false" style="">

      <Element :foreground="true">
        <Title title="Background"></Title>
        <List :row="true" class="scroll-horizontal" style="overflow-x: scroll">
          <Element v-for="bg in state.backgrounds" v-if="state.loaded" :accent="preferences.background == bg"
                   :cb="() => selectBackgroundColor(bg)" :foreground="true"

                   :style="`background-color:${bg} !important;`"
                   class="py-4" style="min-width: 4rem">

          </Element>
        </List>
      </Element>
      <Element :foreground="true">
        <Title title="Accent"></Title>
        <List :row="true" class="scroll-horizontal" style="overflow-x: scroll;">
          <Element v-for="clr in state.colors" v-if="state.loaded" :accent="preferences.accent == clr"
                   :cb="() => selectColor(clr)"
                   :foreground="true"

                   :style="`background-color:${clr} !important;`"
                   class="py-4" style="min-width: 4rem">

          </Element>

        </List>

      </Element>
      <Element :foreground="true">
        <Title title="Accent Opacity"></Title>
        <List :row="true" class="scroll-horizontal" style="overflow-x: scroll">
          <Element v-for="opacity in state.opacities" v-if="state.loaded"
                   :accent="preferences.pattern.opacity == opacity"
                   :cb="() => selectOpacity(opacity)"
                   :foreground="true"

                   :style="`background-color:${preferences.accent} !important; opacity: ${opacity};`"
                   class="py-4" style="min-width: 4rem">

          </Element>
        </List>
      </Element>

      <Element class="" foreground>
        <Title title="Patterns"></Title>
        <div style="height: 100%">
          <List scroll-y style="max-height: 50vh">
            <div class="sample-grid">
              <Element v-for="pattern in state.patternList" v-if="state.loaded"
                       :accent="preferences.pattern.name == pattern"
                       :cb="() => selectPattern(pattern)"
                       :foreground="true"
                       :style="`background-image:${getPattern(pattern)}; `" class="py-4">
                <div class="label-c5 label-o5 px-2 label-w600 lh-1 mb-1">{{ pattern }}</div>
              </Element>
            </div>
          </List>
        </div>
      </Element>
    </List>

  </List>

</template>

<style scoped>
.sample-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(1rem, 1fr));
  grid-gap: 0.25rem;
}
</style>