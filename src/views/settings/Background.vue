<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import List from "udap-ui/components/List.vue";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";

import * as patterns from "udap-ui/vendor/hero-patterns"
import core from "@/core";
import {onBeforeMount} from "vue";

const patternList = Object.keys(patterns)

const preferences = core.preferences()

function selectPattern(pattern: string) {
  preferences.pattern.svg = patterns[pattern](preferences.accent, preferences.pattern.opacity)
  preferences.pattern.name = pattern

}

function hslToHex(h, s, l) {
  l /= 100;
  const a = s * Math.min(l, 1 - l) / 100;
  const f = n => {
    const k = (n + h / 30) % 12;
    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
    return Math.round(255 * color).toString(16).padStart(2, '0');   // convert to Hex and prefix "0" if needed
  };
  return `#${f(0)}${f(8)}${f(4)}`;
}

function selectColor(accent: string) {
  preferences.accent = accent
  preferences.pattern.svg = patterns[preferences.pattern.name](preferences.accent, preferences.pattern.opacity)
}

const scaleOptions = [] as number[]
const colors = []
const backgrounds = []
const blurOptions = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 15, 16, 17, 18, 19, 20, 30]

onBeforeMount(() => {
  let numItems = 64
  for (let i = 0; i < numItems; i++) {
    colors.push(`${hslToHex((360 / numItems) * i, 50, 30)}`)
  }
  for (let i = 0; i < numItems; i++) {
    backgrounds.push(`${hslToHex(0, 13, (100 / numItems) * i)}`)
  }
  for (let i = 0; i <= 20; i++) {
    scaleOptions.push(i * 5)
  }
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

</script>

<template>
  <div class="d-flex gap-1 flex-column">
    <ElementHeader title="Background Scale"></ElementHeader>
    <Element>
      <div class="">
        <List :row="true" style="overflow-x: scroll">
          <Element v-for="w in scaleOptions" :cb="() => selectScale(w)" :foreground="true"
                   :style="`${preferences.pattern.scale == w?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''} `"

                   class="py-4 d-flex justify-content-center align-items-center"
                   style="min-width: 4rem; height: 2rem">
            {{ w }}
          </Element>

        </List>
      </div>
    </Element>
    <ElementHeader title="Foreground Blur"></ElementHeader>
    <Element>
      <div class="">
        <List :row="true" style="overflow-x: scroll">
          <Element v-for="b in blurOptions" :cb="() => selectBlur(b)" :foreground="true"
                   :style="`${preferences.blur== b?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''} `"

                   class="py-4 d-flex justify-content-center align-items-center"
                   style="min-width: 4rem; height: 2rem">
            {{ b }}px
          </Element>

        </List>
      </div>
    </Element>
    <ElementHeader title="Background Color"></ElementHeader>
    <Element>
      <List :row="false">
        <Element :foreground="true">
          <div class="label-c4 label-o5 px-2 label-w600 lh-1 mb-1">Background</div>
          <List :row="true" style="overflow-x: scroll">
            <Element v-for="bg in backgrounds" :cb="() => selectBackgroundColor(bg)" :foreground="true" :style="`${preferences.background == bg?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''} background-color:${bg} !important;`"

                     class="py-4"
                     style="min-width: 4rem">

            </Element>
          </List>
        </Element>
        <Element :foreground="true">
          <div class="label-c4 label-o5 px-2 label-w600 lh-1 mb-1">Accent</div>
          <List :row="true" style="overflow-x: scroll">
            <Element v-for="clr in colors" :cb="() => selectColor(clr)" :foreground="true" :style="`${preferences.accent == clr?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''} background-color:${clr} !important;`"

                     class="py-4"
                     style="min-width: 4rem">

            </Element>
          </List>

        </Element>
      </List>
    </Element>
    <ElementHeader title="Accent Color"></ElementHeader>

    <ElementHeader title="Background Pattern"></ElementHeader>
    <Element>
      <List>
        <Element v-for="pattern in patternList" :cb="() => selectPattern(pattern)" :foreground="true" :style="`background-image:${patterns[pattern](preferences.accent, preferences.pattern.opacity)}; ${preferences.pattern.name== pattern?('box-shadow: inset 0 0 0px 2px rgba(255,255,255,0.1);'):''}`"
                 class="py-4">
          {{ pattern }}
        </Element>
      </List>
    </Element>
  </div>
</template>

<style scoped>

</style>