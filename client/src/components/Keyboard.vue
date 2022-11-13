<script lang="ts" setup>
import Keyboard from "simple-keyboard";
import "simple-keyboard/build/css/index.css";
import {inject, onMounted, onUnmounted, reactive} from "vue";
import type {KeyboardOptions} from "simple-keyboard/build/interfaces";
import type {Haptics} from "@/haptics";


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
const haptics = inject("haptics") as Haptics

onMounted(() => {
    state.keyboard = new Keyboard(`.${props.keyboardClass}`, keyboardOptions)

    window.addEventListener("keydown", typeManual)
})

onUnmounted(() => {
    window.removeEventListener("keydown", typeManual)
})

function onKeyPress(button: string) {
    if (haptics) {
        haptics.tap(1, 1, 100)
    }
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

function typeManual(e: KeyboardEvent) {
    let key: string = e.key
    if (!key) return
    if (key === "Enter") {
        onKeyPress("{enter}")
    } else if (key === "Backspace") {
        onKeyPress("{bksp}")
    } else if (key.length === 1) {
        onKeyPress(key)
    }

}

</script>

<template>
    <div class="keyboard-frame element">
        <div :class="props.keyboardClass" class="simple-keyboard"></div>
    </div>
</template>


<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
.keyboard-frame {
    display: flex;
    /*outline: 1px solid white;*/
    justify-content: center;
    align-items: center;
    position: absolute !important;
    bottom: 4.5rem;
    left: calc(50% - 20.35rem);
    height: 13.5rem;
    width: 40.68rem;

}

.simple-keyboard {
    /*outline: 1px solid white;*/
    position: absolute !important;
    width: 40rem !important;
    /*outline: 1px solid white;*/
    padding: 1rem;

    z-index: 99000 !important;
}
</style>
