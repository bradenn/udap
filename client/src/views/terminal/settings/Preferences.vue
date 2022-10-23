<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, reactive} from "vue";
import type {Preferences} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import axios from "axios";
import Subplot from "@/components/plot/Subplot.vue";
import Scroll from "@/components/Scroll.vue";

let state = reactive({
  loading: true,
})

let preferences = inject("preferences") as Preferences

const defaults = {
  backgrounds: [
    {
      name: "Viridian",
      identifier: "viridian",
    },
    {
      name: "Waves",
      identifier: "waves",
    },
    {
      name: "Swirl",
      identifier: "swirl",
    },
    {
      name: "Lies",
      identifier: "lies",
    },
    {
      name: "Ventura",
      identifier: "ventura",
    },
    {
      name: "Nebula",
      identifier: "nebula",
    },
    {
      name: "Blueberry",
      identifier: "blueberry",
    },
    {
      name: "Galaxy",
      identifier: "glax",
    },
    {
      name: "Geometry",
      identifier: "geometry",
    },
    {
      name: "Void",
      identifier: "void",
    },
    {
      name: "Milky Way",
      identifier: "milk",
    }
  ],
  touchModes: [
    {
      name: "Touch",
      identifier: "touch",
    },
    {
      name: "Cursor",
      identifier: "cursor",
    },
  ],
  screensavers: [
    {
      name: "Bubbles",
      identifier: "bubbles",
    },
    {
      name: "Warp",
      identifier: "warp",
    }
  ],
  themes: [
    {
      name: "Dark",
      identifier: "dark",
    },
    {
      name: "Light",
      identifier: "light",
    },
  ],
  blurModes: [
    {
      name: "Blurred",
      identifier: true,
    },
    {
      name: "Default",
      identifier: false,
    },
  ]
}


function loadImage(image: string) {

  axios.get(`/custom/${image}@2x.png`).then(res => {
    state.loading = false
  }).catch(err => {

  })
}

function changeBackground(name: string): any {
  state.loading = true
  loadImage(name);
  preferences.ui.background.image = name
}

function changeTheme(name: string): any {
  preferences.ui.theme = name
}

function changeTouchmode(mode: string): any {
  preferences.ui.mode = mode
}

function changeScreensaver(screensaver: string): any {
  preferences.ui.screensaver.selection = screensaver
}

function changeBlur(blurred: boolean) {
  preferences.ui.background.blur = blurred
}


</script>

<template>
  <div class="h-100">

    <div class="">
      <div class="element">
        <div class="d-flex align-items-center justify-content-between">
          <div class="label-c1  label-o4 label-w500 px-1 pb-1">Background</div>

        </div>
        <Scroll :horizontal="true" style="max-width: calc(100vw - 2.5rem); overflow-x: scroll;">
          <div class="d-flex gap-1" style="">
            <div v-for="background in defaults.backgrounds" class=" w-100" style="min-width: 8rem;"
                 @click="changeBackground(background.identifier)">
              <div class=" w-100 d-flex justify-content-start subplot " style="padding: 0.125rem;">
                <div :class="`${preferences.ui.background.image === background.identifier?'active':''}`"
                     :style="`background-image: url('/custom/${background.identifier}@2x.png');`"
                     class="background-preview ">
                  <div class="label-xxs label-w500 label-o5 pb-1">
                    {{ background.name }}
                  </div>
                </div>

              </div>

            </div>
          </div>
        </Scroll>


      </div>

      <div class="d-flex w-100 gap justify-content-between mt-2">

        <Plot :cols="2" :rows="1" class="flex-grow-1" title="Themes">
          <Subplot v-for="theme in defaults.themes" :active="preferences.ui.theme === theme.identifier"
                   :fn="() => changeTheme(theme.identifier)"
                   :name="theme.name" @click="">
          </Subplot>
        </Plot>

        <Plot :cols="2" :rows="1" class="flex-grow-1" title="Touch Mode">
          <Subplot v-for="mode in defaults.touchModes" :active="preferences.ui.mode === mode.identifier"
                   :fn="() => changeTouchmode(mode.identifier)"
                   :name="mode.name" @click="">
          </Subplot>
        </Plot>

        <Plot :cols="2" :rows="1" class="flex-grow-1" title="Screensavers">
          <Subplot v-for="screensaver in defaults.screensavers"
                   :active="preferences.ui.screensaver.selection === screensaver.identifier"
                   :fn="() => changeScreensaver(screensaver.identifier)"
                   :name="screensaver.name" @click="">
          </Subplot>
        </Plot>
        <Plot :cols="2" :rows="1" class="flex-grow-1" title="Themes">
          <Subplot v-for="theme in defaults.blurModes" :active="preferences.ui.background.blur === theme.identifier"
                   :fn="() => changeBlur(theme.identifier)"
                   :name="theme.name" @click="">
          </Subplot>
        </Plot>
      </div>
    </div>
  </div>
</template>

<style scoped>
.background-preview {
  display: flex;
  justify-content: end;
  flex-direction: column;
  align-items: center;
  height: 4rem;
  width: 100%;
  border-radius: 0.2rem;

  background-size: cover;

}

.background-preview .active {
  margin: 1rem;
}

.theme-preview {
  width: 8rem;
  display: flex;
  justify-content: center;
  border-radius: 0.5rem;
  box-shadow: 0 0 0.5rem 2px rgba(0, 0, 0, 0.2);
  align-content: center;
  align-items: center;
}


.theme-preview-active {
  background-color: rgba(255, 255, 255, 0.1);
}

.background-preview-active {
  box-shadow: inset 0 0 2px 2px rgba(255, 255, 255, 0.05);
}

</style>