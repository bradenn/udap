<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import SimpleKeyboard from "@/components/Keyboard.vue";
import {reactive, ref} from "vue"

import axios from "axios"

import {Preference} from "@/preferences"
import {PreferenceTypes} from "@/types";

// The reactive component for the spaces and cursor position
let state = reactive({
  spaces: ['', '', '', '', '', '', '', ''],
  cursor: 0
})

// The selected controller url from the previous dialog
let controller = new Preference(PreferenceTypes.Controller).get()

// Errors incurred during verification
let errorMessage = ref("")

function authenticate() {
  // Generated endpoint registration url using the security code and controller address
  let url = `http://${controller}/endpoints/register/${state.spaces.reduce((a, b) => a + b)}`
  // Make the request to the controller nexus
  axios.get(url).then(res => {
    // Set the token in localStorage
    new Preference(PreferenceTypes.Token).set(res.data.token)
    // Redirect the user to the authenticated portal
    window.location.href = "/terminal"
  }).catch(err => {
    // Notify of failures
    errorMessage.value = "Invalid security code. Try again."
    // Reset the input dialog
    resetCode()
  })
}

// Reset the security code input block
function resetCode() {
  // Reset characters
  state.spaces = ['', '', '', '', '', '', '', '']
  // Reset cursor to the zero position
  state.cursor = 0
}

// Enter a character and advance the cursor
function enterChar(char: string) {
  // If the incoming char is a backspace, decrement the cursor and clear the value.
  if (char === "{bksp}") {
    // Only decrement the cursor if it is bigger than zero
    if (state.cursor > 0) state.cursor--
    // Set the char at the cursor to an empty char
    state.spaces[state.cursor] = ''
    // Exit the function
    return
  }
  // Return if the security code block is full
  if (state.cursor >= 8) return
  // Add the provided char to the cursor position
  state.spaces[state.cursor] = char
  // Advance the cursor to the next position
  state.cursor++
}

</script>

<template>

  <div class="container">
    <div class="row  justify-content-center">
    </div>
    <div class="row mt-5 justify-content-center">
      <div class="element p-3 window-medium">
        <router-link class="label-md label-o5 label-w600 py-2" to="/setup/controller"><span
            class="label-o3 label-md">􀯶&nbsp;</span> Change Controller
        </router-link>
        <h3 class="mb-1 mt-2">Authentication.</h3>
        <div class="label-o4">Please enter this terminal's eight-digit security code.</div>
        <div class="d-flex flex-row gap justify-content-between mt-3">
          <div v-for="(v, k) in state.spaces" :key=k
               :class="`${state.cursor === k?'border-fog':'border-transparent'}`"
               class="surface character border label-o4">{{ v }}
          </div>
        </div>
        <div class="d-flex flex-row mt-3 justify-content-between align-items-center text-danger">
          <div class="label-lg label-w500">{{ errorMessage }}</div>
          <div class="surface d-inline-block label-lg label-o5 label-w600 px-4 py-2" @click="authenticate">
            Verify&nbsp;&nbsp;<span class="label-o3 label-md">􀆊</span></div>
        </div>
      </div>
    </div>
    <SimpleKeyboard :input="enterChar" keySet="alpha-numeric" keyboardClass="simple-keyboard"></SimpleKeyboard>
  </div>

</template>


<style scoped>
.border-transparent {
  border-color: transparent !important;
}

.view-medium {
  width: 18rem;
}

.character {
  width: 3rem;
  height: 3rem;
  font-size: 1.6rem;
  display: flex;
  align-content: center;
  align-items: center;
  justify-content: center;
}
</style>