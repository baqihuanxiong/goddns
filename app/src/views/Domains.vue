<template>
  <div>
    <v-card>
      <v-data-table
          :headers="headers"
          :items="domains"
      >
        <template v-slot:top>
          <v-toolbar flat>
            <v-toolbar-title>Domains</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-dialog v-model="dialogNew" max-width="600px">
              <template v-slot:activator="{ on, attrs }">
                <v-btn color="primary" dark class="mb-2" v-bind="attrs" v-on="on">
                  New Domain
                </v-btn>
              </template>
              <v-card>
                <v-card-title>New Domain</v-card-title>
                <v-card-text>
                  <v-container>
                    <v-row>
                      <v-col cols="12" sm="6" md="4">
                        <v-text-field
                            v-model="editedDomain.dn"
                            label="Domain name"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                      <v-select
                          v-model="editedDomain.provider"
                          :items="domainProviders"
                          label="Domain provider"
                      ></v-select>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="closeDialogNew">Cancel</v-btn>
                  <v-btn color="blue darken-1" text @click="saveDomain" :loading="savingDomain">Save</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
            <v-dialog v-model="dialogDelete" max-width="500px">
              <v-card>
                <v-card-title>Are you sure you want to delete this domain?</v-card-title>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="closeDialogDelete">Cancel</v-btn>
                  <v-btn color="blue darken-1" text @click="deleteDomainConfirm" :loading="deletingDomain">OK</v-btn>
                  <v-spacer></v-spacer>
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-toolbar>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-icon small @click="deleteDomain(item)">
            mdi-delete
          </v-icon>
        </template>
      </v-data-table>
    </v-card>
    <v-snackbar v-model="showMessage" timeout="3000" top>{{ message }}</v-snackbar>
  </div>
</template>

<script>
import api from '@/api'

export default {
  name: "Domains",

  data() {
    return {
      message: null,
      showMessage: false,
      savingDomain: false,
      deletingDomain: false,
      dialogNew: false,
      dialogDelete: false,
      headers: [
        {
          text: 'Domain Name',
          align: 'start',
          sortable: false,
          value: 'dn',
        },
        { text: 'Provider', value: 'provider' },
        { text: 'Actions', value: 'actions', sortable: false },
      ],
      domainProviders: [],
      domains: [],
      editedIndex: -1,
      editedDomain: {
        dn: '',
        provider: ''
      },
      defaultDomain: {
        dn: '',
        provider: ''
      }
    }
  },

  created() {
    this.initialize()
  },

  watch: {
    dialogNew (val) {
      val || this.closeDialogNew()
    },
    dialogDelete (val) {
      val || this.closeDialogDelete()
    },
  },

  methods: {
    initialize() {
      api.dns.listDomains().then(res => {
        this.domains = res
        // Object.assign(this.domains, res)
      })

      api.dns.listDomainProviders().then(res => {
        this.domainProviders = res
        // Object.assign(this.domainProviders, res)
      })
    },

    closeDialogNew() {
      this.dialogNew = false
      this.$nextTick(() => {
        this.editedDomain = Object.assign({}, this.defaultDomain)
        this.editedIndex = -1
      })
    },

    closeDialogDelete() {
      this.dialogDelete = false
      this.$nextTick(() => {
        this.editedDomain = Object.assign({}, this.defaultDomain)
        this.editedIndex = -1
      })
    },

    deleteDomainConfirm() {
      this.deletingDomain = true
      api.dns.deleteDomain(this.editedDomain.dn).then(() => {
        this.domains.splice(this.editedIndex, 1)
      }).catch(error => {
        console.log(this.message)
        this.message = error.response.data.error
        this.showMessage = true
      }).finally(() => {
        this.deletingDomain = false
      })
      this.closeDialogDelete()
    },

    saveDomain() {
      this.savingDomain = true
      if (this.editedIndex > -1) {
        Object.assign(this.domains[this.editedIndex], this.editedDomain)
      } else {
        const editedDomain = this.editedDomain
        api.dns.addDomain(editedDomain).then(() => {
          this.domains.push(editedDomain)
        }).catch(error => {
          console.log(this.message)
          this.message = error.response.data.error
          this.showMessage = true
        }).finally(() => {
          this.savingDomain = false
        })
      }
      this.closeDialogNew()
    },

    deleteDomain(item) {
      this.editedIndex = this.domains.indexOf(item)
      this.editedDomain = Object.assign({}, item)
      this.dialogDelete = true
    }
  }
}
</script>

<style scoped>

</style>