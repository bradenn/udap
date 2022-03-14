<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Header from "@/components/Header.vue";
import {inject} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types";

interface Preferences {
  background: string,
  theme: string
}

const preferences = inject("preferences") as Preferences

const defaults = {
  backgrounds: [
    {
      name: "Viridian",
      identifier: "viridian",
    },
    {
      name: "Nerf Dart",
      identifier: "nerfdart",
    },
    {
      name: "Dark Pattern",
      identifier: "cblack",
    },
    {
      name: "Orange Vacuum",
      identifier: "daemon",
    },
    {
      name: "Toaster",
      identifier: "toaster",
    },
    {
      name: "Milky Way",
      identifier: "milk",
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
  ]
}

function changeBackground(name: string) {
  new Preference(PreferenceTypes.Background).set(name)
  preferences.background = name
}

function changeTheme(name: string) {
  new Preference(PreferenceTypes.Theme).set(name)
  preferences.theme = name
}

</script>

<template>
  <div>
    <Header class="px-2 pb-1" icon="bars-progress" name="Preferences"></Header>
    <div class="px-2">
      <div>
        <h4>Background</h4>
        <div class="d-flex gap">
          <div v-for="background in defaults.backgrounds" @click="changeBackground(background.identifier)">
            <div
                :class="`${preferences.background === background.identifier?'background-preview-active':''}`"
                :style="`background-image: url('/custom/${background.identifier}@4x.png');`"
                class="background-preview">
              <div class="label-xxs label-w500 label-o4 pb-1">
                {{ background.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="mt-2">
        <h4>Themes</h4>
        <div class="d-flex gap">
          <div v-for="theme in defaults.themes" @click="changeTheme(theme.identifier)">
            <div
                :class="`${preferences.theme === theme.identifier?'theme-preview-active':''}`"
                class="bg-blur surface theme-preview label-xxs label-o5 label-w600 px-4 py-2">
              <div class="label-xxs label-w500 label-o4">
                {{ theme.name }}
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
  width: 8rem;
  border-radius: 0.5rem;
  aspect-ratio: 16/9;
  background-size: cover;
  box-shadow: 0 0 0.5rem 2px rgba(0, 0, 0, 0.2);
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