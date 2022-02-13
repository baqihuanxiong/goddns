<template>
  <v-app>
    <v-main>
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4 lg4>
            <v-card class="elevation-3 pa-3" min-width="500">
              <v-card-title class="text-h4 font-weight-black primary--text">
                GoDDNS Dashboard Login
              </v-card-title>
              <v-card-text>
                <v-form>
                  <v-text-field
                      v-model="username"
                      label="Username"
                      append-icon="mdi-account"
                      :rules="[rules.required]"/>
                  <v-text-field
                      v-model="password"
                      label="Password"
                      id="password"
                      :type="hidePassword ? 'password' : 'text'"
                      :append-icon="hidePassword ? 'mdi-eye-off' : 'mdi-eye'"
                      :rules="[rules.required]"
                      @click:append="hidePassword = !hidePassword"
                      @keyup.enter="login"/>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn block color="primary" @click="login" :loading="loading">Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-flex>
        </v-layout>
      </v-container>
      <v-snackbar v-model="showResult" timeout="3000" top>{{ result }}</v-snackbar>
    </v-main>
  </v-app>
</template>

<script>
import auth from '@/api/authenticate'

export default {
  data() {
    return {
      loading: false,
      username: null,
      password: null,
      hidePassword: true,
      result: null,
      showResult: false,
      rules: {
        required: value => !!value || 'Required.'
      }
    }
  },

  methods: {
    login() {
      const username = this.username
      const password = this.password

      if (!username || !password) {
        this.result = "Username and Password can't be empty."
        this.showResult = true
        return
      }

      this.loading = true
      auth.login(username, password).then(res => {
        window.localStorage.setItem('jwt', res.jwt)
        this.$router.push({ path: 'Overview' })
      }).catch(error => {
        this.result = error.response.data.error
        this.showResult = true
      }).finally(() => {
        this.password = null
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>

</style>