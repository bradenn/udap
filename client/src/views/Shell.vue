<script>

export default {
  components: {Element},
  data() {
    return {
      command: "",
      predictions: ""
    }
  },
  created() {

  },
  methods: {
    loadPredictions() {
      if (this.command.startsWith("meta")) {
        this.predictions = this.$root
      } else {
        this.predictions = ""
      }
    }
  },
}
</script>

<template>
  <div class="d-flex justify-content-start align-items-center flex-column h-100 mt-5">
    <div class="col-8">
      <input v-model="command" class="form-control terminal-input " placeholder="udap"
             style="border-bottom-right-radius: 0 !important; border-bottom-left-radius: 0 !important;"
             type="text" v-on:keyup="loadPredictions"/>
    </div>
    <div class="col-8">
      <div class="terminal-input scroll"
           style="border-top-right-radius: 0 !important; border-top-left-radius: 0 !important;">
        <h6 class="fw-bold mt-2">Raw JSON</h6>
        <hr class="mt-0 mb-1">
        <div class="code h-100">
          <pre>{{ this.$root }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

.scroll {
  scroll-behavior: auto;
}

.code {
  max-height: 50vh;
  padding: 0.5em 0.25em;
  overflow-x: hidden;
  overflow-y: scroll;
  font-size: 14px !important;
  font-weight: 500;
  font-family: 'JetBrains Mono', sans-serif;
}
</style>
