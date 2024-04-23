<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import core from "@/core";
import type {RouteRecord} from "vue-router"
import {AppMeta} from "@/views/apps/index";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import {onMounted, reactive} from "vue";
import * as patterns from "udap-ui/vendor/hero-patterns"
import List from "udap-ui/components/List.vue";

const router = core.router()

const apps = [
  {
    name: "Monitor",
    link: "/monitor",
    icon: '􁂥'
  },
  {
    name: "Thermostat",
    link: "/thermostat",
    icon: '􁁋'
  },
  {
    name: "Lasers",
    link: "/sentry",
    icon: '􀑃'
  },
  {
    name: "Lights",
    link: "/lights",
    icon: '􁓼'
  },
  {
    name: "Todo",
    link: "/apps/todo",
    icon: '􀷾'
  },
  {
    name: "Message",
    link: "/apps/message",
    icon: '􀈠'
  },
  {
    name: "Sensors",
    link: "/apps/sensors",
    icon: '􁔊'
  },
  {
    name: "Actions",
    link: "/apps/actions",
    icon: '􀒘'
  },
  {
    name: "Conversions",
    link: "/apps/conversions",
    icon: '􀅮'
  },
  {
    name: "Movies",
    link: "/apps/movies",
    icon: '􀜤'
  },
  {
    name: "Network",
    link: "/apps/network",
    icon: '􁆬'
  },
]

const preferences = core.preferences();

interface AppSection {
  title: string,
  apps: AppMeta[],

}

const state = reactive({
  groups: [] as AppSection[],
  pattern: "" as string,
  patternName: "heavyRain" as string,
  brightness: 50 as number,
  scale: 50 as number,
  color: 10 as number,
  patternList: Object.keys(patterns),
})

onMounted(() => {
  getApps()
  state.patternName = preferences.pattern.name

  selectPattern(state.patternName)
})

function getApps(): AppMeta[] {
  let appMeta: AppMeta[] = []
  let routes: RouteRecord[] = router.getRoutes();

  if (!routes) {
    return [];
  }

  let sections: AppSection[] = []

  for (let route of routes) {
    if (!route.path.startsWith("/apps/")) continue


    let meta = route.meta as unknown as AppMeta
    if (!meta || Object.keys(meta).length == 0) continue
    meta.link = route.path
    if (meta.child) continue
    if (meta.section) {
      let section = sections.find(s => s.title == meta.section);
      if (section) {
        section.apps.push(meta)
        sections = sections.map(s => s.title == meta.section ? section : s)
      } else {
        sections.push({title: meta.section, apps: [meta]})
      }
    } else {
      let section = sections.find(s => s.title == "default");
      if (section) {
        section.apps.push(meta)
        sections = sections.map(s => s.title == "default" ? section : s)
      } else {
        sections.push({title: "default", apps: [meta]})
      }
    }

    // if (media) {
    //   if (meta.section && meta.section == "media") {
    //     appMeta.push(meta)
    //   }
    //
    // } else {
    //   if (meta.child || meta.section) continue
    //   appMeta.push(meta)
    // }

  }
  state.groups = sections
  return appMeta
}

function getPattern(name: string): string {
  //@ts-ignore
  try {
    return patterns[name](preferences.accent, preferences.pattern.opacity)
  } catch {
    return ""
  }
}

function previewPattern(name: string): string {
  //@ts-ignore
  try {
    return patterns[name]("#666", 50)
  } catch {
    return ""
  }
}

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

function selectColor(scale: number) {
  state.color = scale
  selectPattern(state.patternName)
}

function selectScale(scale: number) {
  state.scale = scale
}


function selectBrightness(scale: number) {
  state.brightness = scale
  selectPattern(state.patternName)
}

function selectPattern(pattern: string) {
  state.patternName = pattern
  state.pattern = getPattern(pattern)
}

</script>

<template>
  <div class="scrollable-fixed">
    <List v-if="state.groups.length > 0" scroll-y>
      <div v-for="group in state.groups" :key="group.title">
        <ElementHeader :title="group.title" class="text-capitalize mb-1"></ElementHeader>
        <div class="app-row">
          <div v-for="app in group.apps" :key="app.link"
               class="d-flex flex-column align-items-center justify-content-center mb-4 ">
            <div style="aspect-ratio: 1/1; width: 100%;">
              <Element :style="`aspect-ratio: 1/1;`"
                       :to="app.link || ''"
                       class="subplot app-pattern d-flex align-items-center justify-content-center"
                       mutable>

                <div class="sf label-c0 icon-blur" style="font-size: 1.7rem">
                  <div>{{ app.icon }}</div>
                </div>
              </Element>

            </div>
            <div class="label-c6 label-o5 label-w600 pt-1">{{ app.name }}</div>
          </div>

        </div>
      </div>
      <!--      <Title title="Brightness (0 - 100%)"></Title>-->
      <!--      <List class="px-2">-->
      <!--        <Slider :min="0" :max="100" bg="dim" :change="(scale) => selectBrightness(scale)" :step="1"-->
      <!--                :value="state.brightness">-->
      <!--          dsds-->
      <!--        </Slider>-->
      <!--      </List>-->
      <!--      <Title title="Color (0deg - 360deg)"></Title>-->
      <!--      <List class="px-2">-->
      <!--        <Slider :min="0" :max="360" bg="hue" :change="(scale) => selectColor(scale)" :step="1" :value="state.color">-->
      <!--          dsds-->
      <!--        </Slider>-->
      <!--      </List>-->

      <!--      <Title title="Patterns"></Title>-->

      <!--      <List scroll-x>-->
      <!--        <div class="app-row">-->
      <!--          <Element v-for="pattern in state.patternList" :key="pattern"-->
      <!--                   :accent="preferences.pattern.name == pattern"-->
      <!--                   :cb="() => selectPattern(pattern)"-->
      <!--                   :foreground="true"-->
      <!--                   :style="`background-image:${previewPattern(pattern)}; display: flex; align-items: center; justify-content: center;`"-->
      <!--                   class="py-4">-->
      <!--            <div class="label-c5 label-o5 px-2 label-w600 lh-1 mb-1">{{ pattern }}</div>-->
      <!--          </Element>-->
      <!--        </div>-->
      <!--      </List>-->

    </List>
  </div>
</template>

<style lang="scss" scoped>

.app-pattern {

  background-position: center;
  width: 100%;

  aspect-ratio: 1/1;

}

//
//.icon-blur {
//  align-items: center;
//  justify-content: center;
//  width: 2rem;
//  height: 2rem;
//  background-color: rgba(255, 255, 255, 0.5);
//  backdrop-filter: blur(12px) !important;
//  -webkit-backdrop-filter: blur(12px) !important;
//}

.app-row {
  display: grid;
  margin-inline: 0.5rem;
  grid-template-columns: repeat(5, minmax(2rem, 1fr));
  grid-gap: 0.5rem;
}

.app-grid {
  display: grid;
  padding: 1rem;
  grid-template-columns: repeat(4, minmax(2rem, 1fr));
  grid-gap: 2rem;
}
</style>