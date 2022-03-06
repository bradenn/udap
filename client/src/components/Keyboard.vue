<script>
import Keyboard from "simple-keyboard";
import "simple-keyboard/build/css/index.css";

export default {
  name: "SimpleKeyboard",
  props: {
    keyboardClass: {
      default: "simple-keyboard",
      type: String
    },
    input: {
      type: Function
    },
    keyset: {
      type: String
    }
  },
  data: () => ({
    keyboard: null
  }),
  mounted() {
    this.keyboard = new Keyboard(this.keyboardClass, {
      onChange: this.onChange,
      onKeyPress: this.onKeyPress,
      theme: "hg-theme-default hg-theme-udap",
      physicalKeyboardHighlight: true,
      syncInstanceInputs: true,
      mergeDisplay: true,
      layout: {
        'default': this.$props.keyset === "alpha-numeric" ? [
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

    });
  },
  methods: {
    onChange(input) {
      this.$emit("onChange", input);
    },
    onKeyPress(button) {
      this.$emit("onKeyPress", button);
      this.$props.input(button)
      /**
       * If you want to handle the shift and caps lock buttons
       */
      if (button === "{shift}" || button === "{lock}") this.handleShift();
    },
    handleShift() {
      let currentLayout = this.keyboard.options.layoutName;
      let shiftToggle = currentLayout === "default" ? "shift" : "default";

      this.keyboard.setOptions({
        layoutName: shiftToggle
      });
    }
  },
  watch: {
    input(input) {
      this.keyboard.setInput(input);
    }
  }
};
</script>

<template>
  <div :class="keyboardClass"></div>
</template>


<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>

.simple-keyboard {
  position: absolute !important;
  width: 100%;
  bottom: 0;
  left: 0;

}
</style>
