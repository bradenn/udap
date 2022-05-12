<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Header from "@/components/Header.vue";
import {inject, reactive} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import axios from "axios";
import Subplot from "@/components/plot/Subplot.vue";

interface Preferences {
  ui: {
    background: string
    theme: string
    mode: string
    blur: number
    grid: boolean
    watermark: boolean
    night: boolean
    outlines: boolean
  }
}

let state = reactive({
  loading: true,
})

const preferences = inject("preferences") as Preferences

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
      name: "Lies",
      identifier: "lies",
    },
    {
      name: "Yellow Lies",
      identifier: "yellowLies",
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
  themes: [
    {
      name: "Dark",
      identifier: "dark",
    },
    {
      name: "Light",
      identifier: "light",
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
  new Preference(PreferenceTypes.Background).set(name)
  state.loading = true
  loadImage(name);
  preferences.ui.background = name
}

function changeTheme(name: string): any {
  new Preference(PreferenceTypes.Theme).set(name)
  preferences.ui.theme = name
}

function changeTouchmode(mode: string): any {
  new Preference(PreferenceTypes.TouchMode).set(mode)
  preferences.ui.mode = mode
}

</script>

<template>
  <div>
    <Header class="px-2 pb-1" icon="bars-progress" name="Preferences"></Header>
    <div class="px-2">
      <div>
        <h4>Background</h4>
        <Plot :cols="5" :rows="2">
          <div v-for="background in defaults.backgrounds" @click="changeBackground(background.identifier)">
            <div class=" w-100 d-flex justify-content-start subplot " style="padding: 0.125rem;">
              <div :class="`${preferences.ui.background === background.identifier?'active':''}`"
                   :style="`background-image: url('/custom/${background.identifier}@2x.png');`"
                   class="background-preview ">
                <div class="label-xxs label-w500 label-o5 pb-1">
                  {{ background.name }}
                </div>
              </div>

            </div>

          </div>
        </Plot>

      </div>

      <div class="mt-2">
        <h4>Themes</h4>
        <Plot :cols="2" :rows="1">
          <Subplot v-for="theme in defaults.themes" :active="preferences.ui.theme === theme.identifier"
                   :fn="() => changeTheme(theme.identifier)"
                   :name="theme.name" @click="">
          </Subplot>
        </Plot>

      </div>
      <div class="mt-2">
        <h4>Touch Mode</h4>
        <Plot :cols="2" :rows="1">
          <Subplot v-for="mode in defaults.touchModes" :active="preferences.ui.mode === mode.identifier"
                   :fn="() => changeTouchmode(mode.identifier)"
                   :name="mode.name" @click="">
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