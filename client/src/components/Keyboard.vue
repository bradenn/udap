<script lang="ts" setup>
import Keyboard from "simple-keyboard";
import "simple-keyboard/build/css/index.css";
import {onMounted, onUnmounted, reactive} from "vue";
import type {KeyboardOptions} from "simple-keyboard/build/interfaces";

const props = defineProps<{
  keyboardClass: string
  input: (data: string) => void
  keySet: string,
}>()

let keyboardOptions: KeyboardOptions = {
  onKeyPress: onKeyPress,
  theme: "hg-theme-default hg-theme-udap",
  physicalKeyboardHighlight: false,
  syncInstanceInputs: true,
  mergeDisplay: true,
  preventMouseDownDefault: true,
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
      '{control} {alt} {meta} {space} {meta} {alt} {control}'
    ],
    'shift': [
      '~ ! @ # $ % ^ & * ( ) _ + {bksp}',
      '{tab} Q W E R T Y U I O P { } |',
      '{lock} A S D F G H J K L : " {enter}',
      '{shift} Z X C V B N M < > ? {shift}',
      '{control} {alt} {meta} {space} {meta} {alt} {control}'
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
    "{space}": "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" +
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" +
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" +
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" +
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" +
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;",
    "{meta}": "⌘",
    "{at}": " ",
  }
}

interface KeyboardState {
  keyboard: Keyboard
  shift: boolean
}

let state: KeyboardState = reactive({
  keyboard: {} as Keyboard,
  shift: false
})

onMounted(() => {
  state.keyboard = new Keyboard(`.${props.keyboardClass}`, keyboardOptions)

  window.addEventListener("keydown", typeManual)
})

onUnmounted(() => {
  window.removeEventListener("keydown", typeManual)
})

function onKeyPress(button: string) {
  if (button === "{shift}" || button === "{lock}") {
    state.shift = !state.shift;
    handleShift();
    return
  }

  props.input(button)
  if (state.shift) {
    handleShift();
    state.shift = false
  }
}


function handleShift() {

  let currentLayout = state.keyboard.options.layoutName;
  let shiftToggle = currentLayout === "default" ? "shift" : "default";

  state.keyboard.setOptions({
    layoutName: shiftToggle
  });
}

function typeManual(e: any) {
  let keyPress = e["key"] || ""
  if (!keyPress) return
  if (keyPress === "Backspace") {
    onKeyPress("{bksp}")
  } else if (keyPress === "Shift" || keyPress.length > 1) {
    return
  } else {
    onKeyPress(keyPress)
  }

}

</script>

<template>
  <div :class="props.keyboardClass" class="simple-keyboard element"></div>
</template>


<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>

.simple-keyboard {
  position: absolute !important;
  width: 40rem;


  z-index: 1000;
  bottom: 2%;
  left: calc(50% - 20rem);

}
</style>
