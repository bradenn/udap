<script>
export default {
  data() {
    return {
      key: "",
      error: "",
      connection: ""
    }
  },
  created() {
    if (this.$root.session.token !== '' && this.$root.session.token !== 'unset') this.$router.push("/terminal/home")
  },
  methods: {
    register(key) {
      this.$http.get(`http://${this.$root.config.host}:${this.$root.config.port}/endpoints/register/${key}`)
          .then(this.success).catch(this.rejected)
    },
    test() {
      this.$http.get(`http://${this.$root.config.host}:${this.$root.config.port}/status`)
          .then(this.testResult).catch(this.rejected)
    },
    testResult(result) {
      if (result.status === 200) {
        this.connection = "<span class='text-success'>Connection succeeded.</span>"
      } else {
        this.connection = "<span class='text-danger'>Connection failed.</span>"
      }
    },
    success(result) {
      this.$root.session.token = result.data.token
      this.$router.push("/terminal/home")
    },
    rejected(result) {
      this.error = result
    }
  },
}
</script>

<template>
  <div class="container">
    <div class="row">
      <div class="col-12 d-flex justify-content-around mt-5">
        <div class="element px-3">
          <h4 class="my-1">Register Endpoint</h4>
          <div v-if="error" class="small text-danger">{{ error }}</div>
          <form class=" pt-2 pb-2">
            <label class="form-label">Host</label>
            <div class="input-group mb-1">
              <input v-model="this.$root.config.host" class="form-control" placeholder="localhost" type="text"/>
            </div>
            <label class="form-label">Port</label>
            <div class="input-group mb-1">
              <input v-model="this.$root.config.port" class="form-control" placeholder="3020" type="text"/>
            </div>
            <div class="d-flex justify-content-between align-items-center mb-1">
              <span class="small text-muted" v-html="this.connection"></span>
              <a class="text-muted btn btn-md btn-fog my-2" href="#" v-on:click="test()">
                Test Connection<i class="bi bi-arrow-repeat" style="margin-left: 0.75em;"></i>
              </a>
            </div>

            <label class="form-label">Access Key</label>
            <div class="input-group">
              <input v-model="key" class="form-control" placeholder="- - - - - - - -" type="text"/>
            </div>
            <div class="form-text">Endpoint keys can be generated from a diagnostic terminal.</div>
            <a class="text-muted btn btn-md btn-fog my-2 float-end" href="#" v-on:click="register(key)">
              Register <i class="bi bi-caret-right"></i>
            </a>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
