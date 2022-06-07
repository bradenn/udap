<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Keypad from "@/components/Keypad.vue";
import {inject, reactive} from "vue";
import axios from "axios";
import type {Remote, Session, User} from "@/types";

let remote = inject("remote") as Remote
let session = inject("session") as Session

let state = reactive({
  trial: ""
})

function withPasscode(passkey: string) {
  let target = remote.users.find(u => u.username === "bradenn")
  if (!target) return
  target.password = passkey
  axios.post("http://localhost:3020/users/authenticate", JSON.stringify(target)).then(res => {
    session.user = target as User
  }).catch(err => {
    console.log(err)
    session.user = {} as User
    state.trial = err.data
  })
}

</script>

<template>
  <div class="h-100 d-flex justify-content-center align-content-center mt-4">
    {{ state.trial }}
    <Keypad :onEnter="withPasscode"></Keypad>
  </div>
</template>

<style lang="scss" scoped>
</style>