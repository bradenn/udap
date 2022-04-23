<script lang="ts" setup>
import Keyboard from "simple-keyboard";
import "simple-keyboard/build/css/index.css";
import {onMounted} from "vue";
import type {KeyboardOptions} from "simple-keyboard/build/interfaces";

const props = defineProps<{
  keyboardClass: string
  input: (data: string) => void
  keySet: string,
}>()

let keyboardOptions: KeyboardOptions = {
  onKeyPress: onKeyPress,
  theme: "hg-theme-default hg-theme-udap",
  physicalKeyboardHighlight: true,
  syncInstanceInputs: true,
  mergeDisplay: true,
  layout: {
    'default': props.keySet === "alpha-numeric" ? [
      '1 2 3 4 5 6 7 8 9 0 {bksp}',
      'Q W E R T Y U I O P',
      'A S D F G H J K L',
      'Z X C V B N M',
    ] : [
      '` 1 2 3 4 5 6 7 8 9 0 - = {bksp}',
      '{tab} q w e r t y u i o p [ ] \\',
      '{lock} a s d f g h j k l ; \' {enter}',
      '{shift} z x c v b n m , . / {shift}',
    ],
    'shift': [
      '~ ! @ # $ % ^ & * ( ) _ + {bksp}',
      '{tab} Q W E R T Y U I O P { } |',
      '{lock} A S D F G H J K L : " {enter}',
      '{shift} Z X C V B N M < > ? {shift}',
      '.com @ {space}'
    ]
  },
  display: {
    "{tab}": "tab",
    "{bksp}": "delete",
    "{enter}": "return",
    "{lock}": "caps",
    "{shift}": "shift",
    "{control}": "⌃",
    "{alt}": "⌥",
    "{space}": "&nbsp;",
    "{meta}": "⌘",
  }
}

interface KeyboardState {
  keyboard: Keyboard
}

let state: KeyboardState = {
  keyboard: {} as Keyboard
}
onMounted(() => {
  state.keyboard = new Keyboard(`.${props.keyboardClass}`, keyboardOptions)

})

function onKeyPress(button: string) {
  props.input(button)
  if (button === "{shift}" || button === "{lock}") handleShift();
}

function handleShift() {
  let currentLayout = state.keyboard.options.layoutName;
  let shiftToggle = currentLayout === "default" ? "shift" : "default";

  state.keyboard.setOptions({
    layoutName: shiftToggle
  });
}

</script>

<template>
  <div :class="props.keyboardClass" class="simple-keyboard"></div>
</template>


<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>

.simple-keyboard {
  position: absolute !important;
  width: 100%;
  z-index: 1000;
  bottom: 0;
  left: 0;

}
</style>
