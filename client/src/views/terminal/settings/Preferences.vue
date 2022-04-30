<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Header from "@/components/Header.vue";
import {inject} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types";
import Plot from "@/components/plot/Plot.vue";

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
      name: "Blueberry",
      identifier: "blueberry",
    },
    {
      name: "Dark Pattern",
      identifier: "cblack",
    },
    {
      name: "Dust",
      identifier: "dust",
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

function changeBackground(name: string) {
  new Preference(PreferenceTypes.Background).set(name)
  preferences.ui.background = name
}

function changeTheme(name: string) {
  new Preference(PreferenceTypes.Theme).set(name)
  preferences.ui.theme = name
}

function changeTouchmode(mode: string) {
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
            <div class="subplot w-100 d-flex justify-content-start" style="padding: 0.125rem;">
              <div :class="`${preferences.ui.background === background.identifier?'background-preview-active':''}`"
                   :style="`background-image: url('/custom/${background.identifier}@4x.png');`"
                   class="background-preview">
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
        <div class="d-flex gap">
          <div v-for="theme in defaults.themes" @click="changeTheme(theme.identifier)">
            <div
                :class="`${preferences.ui.theme === theme.identifier?'theme-preview-active':''}`"
                class="bg-blur surface theme-preview label-xxs label-o5 label-w600 px-4 py-2">
              <div class="label-xxs label-w500 label-o4">
                {{ theme.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="mt-2">
        <h4>Touch Mode</h4>
        <div class="d-flex gap">
          <div v-for="mode in defaults.touchModes" @click="changeTouchmode(mode.identifier)">
            <div
                :class="`${preferences.ui.mode === mode.identifier?'theme-preview-active':''}`"
                class="bg-blur surface theme-preview label-xxs label-o5 label-w600 px-4 py-2">
              <div class="label-xxs label-w500 label-o4">
                {{ mode.name }}
              </div>
            </div>
          </div>
        </div>
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
  border-radius: 0.25rem;

  background-size: cover;

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