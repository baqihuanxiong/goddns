<template>
  <div>
    <v-container fluid>
      <v-card elevation="2">
        <v-card-title>Scanned IP: {{ currentIP }}</v-card-title>
        <v-card-text>
          Current using: {{ lastStatus.ip }}
          <br>
          Updated at {{ lastStatus.updatedAt }}
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary lighten-1" text @click="refreshStatus">
            Refresh
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import api from '@/api'

export default {
  data() {
    return {
      currentIP: "",
      lastStatus: {}
    }
  },

  created() {
    this.refreshStatus()
  },

  methods: {
    refreshStatus() {
      api.ip.lookupIP().then(res => {
        this.currentIP = res.ip
      })
      api.ip.getLastIP().then(res => {
        this.lastStatus = res
      })
    }
  }
}
</script>

<style scoped>

</style>
